package localauth

import (
	"context"

	"github.com/google/uuid"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	mailerdomain "mandacode.com/accounts/auth/internal/domain/infra/mailer"
	dbmodels "mandacode.com/accounts/auth/internal/domain/models/database"
	codedomain "mandacode.com/accounts/auth/internal/domain/repository/code"
	dbdomain "mandacode.com/accounts/auth/internal/domain/repository/database"
	tokendomain "mandacode.com/accounts/auth/internal/domain/repository/token"
	localauthdomain "mandacode.com/accounts/auth/internal/domain/usecase/localauth"
)

type signupUsecase struct {
	authAccount      dbdomain.AuthAccountRepository
	localAuth        dbdomain.LocalAuthRepository
	token            tokendomain.TokenRepository
	mailer           mailerdomain.Mailer
	emailCodeManager codedomain.CodeManager
	verifyEmailURL   string
}

// ResendVerificationEmail implements localauthdomain.SignupUsecase.
func (s *signupUsecase) ResendVerificationEmail(ctx context.Context, email string) (success bool, err error) {
	panic("unimplemented")
}

// Signup implements localauthdomain.SignupUsecase.
func (s *signupUsecase) Signup(ctx context.Context, input localauthdomain.SignupInput) (userID uuid.UUID, err error) {
	auth, err := s.localAuth.GetLocalAuthByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, errcode.ErrNotFound) {
		joinedErr := errors.Join(err, "failed to get user by email")
		return uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	if auth != nil {
		return uuid.Nil, errors.New("email already registered", PubConflict, errcode.ErrConflict)
	}

	// Create auth account
	account, err := s.authAccount.CreateAuthAccount(ctx, &dbmodels.CreateAuthAccountInput{
		UserID:   uuid.New(),
		IsActive: true,
	})
	if err != nil {
		joinedErr := errors.Join(err, "failed to create auth account")
		return uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "AuthAccountCreationFailed")
	}

	auth, err = s.localAuth.CreateLocalAuth(
		ctx,
		&dbmodels.CreateLocalAuthInput{
			AccountID:  account.ID,
			Email:      input.Email,
			Password:   input.Password,
			IsActive:   true,
			IsVerified: false,
		},
	)
	if err != nil {
		joinedErr := errors.Join(err, "failed to create user")
		return uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	// Generate email verification token
	code, err := s.emailCodeManager.IssueCode(ctx, account.UserID)

	token, _, err := s.token.GenerateEmailVerificationToken(ctx, account.UserID, input.Email, code)
	if err != nil {
		joinedErr := errors.Join(err, "failed to generate email verification token")
		return uuid.Nil, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	// Send verification email
	url := s.verifyEmailURL + "?token=" + token
	err = s.mailer.SendEmailVerificationMail(input.Email, url)
	if err != nil {
		return uuid.Nil, errors.Join(err, "failed to send verification email")
	}

	return account.UserID, nil
}

// VerifyEmail implements localauthdomain.SignupUsecase.
func (s *signupUsecase) VerifyEmail(ctx context.Context, email string, token string) (success bool, err error) {
	result, err := s.token.VerifyEmailVerificationToken(ctx, token)
	if err != nil {
		joinedErr := errors.Join(err, "failed to verify email verification token")
		return false, errors.Upgrade(joinedErr, errcode.ErrUnauthorized, PubAuthenticationFailed)
	}
	if !result.Valid {
		return false, errors.New("invalid or expired token", PubAuthenticationFailed, errcode.ErrUnauthorized)
	}

	// Check if the user exists and is not already verified
	auth, err := s.localAuth.GetLocalAuthByEmail(ctx, result.Email)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get user by ID")
		return false, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	if auth.Email != email {
		return false, errors.New("email does not match", PubAuthenticationFailed, errcode.ErrUnauthorized)
	}
	if auth.IsVerified {
		return false, errors.New("user already verified", PubConflict, errcode.ErrConflict)
	}

	// Check if the verification code matches
	account, err := s.authAccount.GetAuthAccountByUserID(ctx, auth.AccountID)
	if err != nil {
		joinedErr := errors.Join(err, "failed to get auth account by user ID")
		return false, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	// Verify the code
	valid, err := s.emailCodeManager.ValidateCode(ctx, account.UserID, result.Code)
	if err != nil {
		joinedErr := errors.Join(err, "failed to validate verification code")
		return false, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}
	if !valid {
		return false, errors.New("verification code is invalid or expired", PubAuthenticationFailed, errcode.ErrUnauthorized)
	}

	// Update the user's verification status
	_, err = s.localAuth.SetIsVerified(ctx, auth.AccountID, true)
	if err != nil {
		joinedErr := errors.Join(err, "failed to update user verification status")
		return false, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, PubInternalFailure)
	}

	return true, nil
}

// NewSignupUsecase creates a new instance of SignupUsecase.
func NewSignupUsecase(
	authAccount dbdomain.AuthAccountRepository,
	localAuth dbdomain.LocalAuthRepository,
	token tokendomain.TokenRepository,
	mailer mailerdomain.Mailer,
	emailCodeManager codedomain.CodeManager,
	verifyEmailURL string,
) localauthdomain.SignupUsecase {
	return &signupUsecase{
		authAccount:      authAccount,
		localAuth:        localAuth,
		token:            token,
		mailer:           mailer,
		emailCodeManager: emailCodeManager,
		verifyEmailURL:   verifyEmailURL,
	}
}
