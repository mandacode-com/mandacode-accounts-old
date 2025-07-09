package localuserapp

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	localuserdomain "mandacode.com/accounts/user/internal/domain/port/auth/local"
	mailerdomain "mandacode.com/accounts/user/internal/domain/port/mailer"
	profiledomain "mandacode.com/accounts/user/internal/domain/port/profile"
	repodomain "mandacode.com/accounts/user/internal/domain/port/repository"
	tokendomain "mandacode.com/accounts/user/internal/domain/port/token"
	usereventdomain "mandacode.com/accounts/user/internal/domain/port/user_event"
	"mandacode.com/accounts/user/internal/util"
)

type localUserApp struct {
	repository       repodomain.UserRepository
	localUserService localuserdomain.LocalUserService
	profileService   profiledomain.ProfileService
	tokenService     tokendomain.TokenService
	mailerService    mailerdomain.Mailer
	userEventService usereventdomain.UserEventService
	codeGenerator    util.RandomCodeGenerator
}

// CreateUser implements LocalUserApp.
func (l *localUserApp) CreateUser(ctx context.Context, email string, password string) (string, error) {
	var g errgroup.Group

	userID := uuid.New()

	// Initialize rollback functions
	rollbackFuncs, fail := util.NewRollback()
	rollbackFuncs = append(rollbackFuncs, func() {
		l.userEventService.UserCreationFailed(userID)
	})

	// Generate a sync code for the user
	syncCode, err := l.codeGenerator.GenerateSyncCode()
	if err != nil {
		return "", err
	}

	// Create the user in the repository
	g.Go(func() error {
		if _, err := l.repository.CreateUser(userID, syncCode); err != nil {
			return err
		}
		rollbackFuncs = append(rollbackFuncs, func() {
			l.repository.DeleteUser(userID)
		})
		return nil
	})

	// Enroll the user with local user service and create a profile
	g.Go(func() error {
		if _, err := l.localUserService.EnrollUser(ctx, userID, email, password, true, false); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if _, err := l.profileService.CreateProfile(ctx, userID); err != nil {
			return err
		}
		return nil
	})

	// If any of the goroutines fail, rollback the changes
	if err := g.Wait(); err != nil {
		fail()
		return "", err
	}

	// Send email verification
	emailVerificationCode, err := l.codeGenerator.GenerateEmailVerificationCode()
	if err != nil {
		fail()
		return "", err
	}
	if _, err := l.repository.UpdateEmailVerificationCode(userID, emailVerificationCode); err != nil {
		fail()
		return "", err
	}
	emailVerificationToken, _, err := l.tokenService.GenerateEmailVerificationToken(ctx, userID, email, emailVerificationCode)
	if err != nil {
		fail()
		return "", err
	}
	if err := l.mailerService.SendEmailVerificationMail(email, emailVerificationToken); err != nil {
		fail()
		return "", err
	}

	return userID.String(), nil
}

// UpdateEmail implements LocalUserApp.
// func (l *localUserApp) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error {
// 	panic("unimplemented")
// }

// UpdatePassword implements LocalUserApp.
func (l *localUserApp) UpdatePassword(ctx context.Context, userID uuid.UUID, currentPassword string, newPassword string) error {
	if _, err := l.localUserService.UpdatePassword(ctx, userID, currentPassword, newPassword); err != nil {
		return err
	}

	return nil
}

// VerifyUserEmail implements LocalUserApp.
func (l *localUserApp) VerifyUserEmail(ctx context.Context, userID uuid.UUID, verificationToken string) error {
	var g errgroup.Group

	valid, userIDStr, email, code, err := l.tokenService.VerifyEmailVerificationToken(ctx, verificationToken)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("invalid verification token")
	}
	if userIDStr == nil || email == nil || code == nil {
		return errors.New("verification token is missing required fields")
	}

	// Check if the code in the token matches the user's code in the repository
	g.Go(func() error {
		userUID, err := uuid.Parse(*userIDStr)
		if err != nil {
			return errors.New("invalid user ID in verification token")
		}

		// Check if the user ID in the token matches the user ID in the repository
		repoUser, err := l.repository.GetUser(userUID)
		if err != nil {
			return err
		}
		if repoUser.EmailVerificationCode != *code {
			return errors.New("verification code in token does not match user's code")
		}
		return nil
	})

	// Check if the authenticated user Email matches the email in the token
	g.Go(func() error {
		getUserResponse, err := l.localUserService.GetUser(ctx, userID)
		if err != nil {
			return err
		}
		if getUserResponse == nil {
			return errors.New("user not found")
		}
		authUser := getUserResponse.User
		if authUser.Email != *email {
			return errors.New("email in token does not match user's email")
		}
		return nil
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		return err
	}

	if _, err := l.localUserService.UpdateVerifiedStatus(ctx, userID, true); err != nil {
		return err
	}

	if _, err := l.repository.UpdateEmailVerificationCode(userID, *code); err != nil {
		return err
	}

	return nil
}

// NewLocalUserApp creates a new instance of localUserApp with the provided dependencies.
func NewLocalUserApp(
	repository repodomain.UserRepository,
	localUserService localuserdomain.LocalUserService,
	profileService profiledomain.ProfileService,
	tokenService tokendomain.TokenService,
	mailerService mailerdomain.Mailer,
	userEventService usereventdomain.UserEventService,
) LocalUserApp {
	return &localUserApp{
		repository:       repository,
		localUserService: localUserService,
		profileService:   profileService,
		tokenService:     tokenService,
		mailerService:    mailerService,
		userEventService: userEventService,
	}
}
