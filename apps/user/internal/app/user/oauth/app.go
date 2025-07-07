package oauthuserapp

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	providerv1 "mandacode.com/accounts/proto/common/provider/v1"
	oauthuserdomain "mandacode.com/accounts/user/internal/domain/port/auth/oauth"
	profiledomain "mandacode.com/accounts/user/internal/domain/port/profile"
	repodomain "mandacode.com/accounts/user/internal/domain/port/repository"
	usereventdomain "mandacode.com/accounts/user/internal/domain/port/user_event"
	"mandacode.com/accounts/user/internal/util"
)

type oauthUserApp struct {
	repository       repodomain.UserRepository
	oauthUserService oauthuserdomain.OAuthUserService
	profileService   profiledomain.ProfileService
	userEventService usereventdomain.UserEventService
	codeGenerator    util.RandomCodeGenerator
}

// CreateUser implements OAuthUserApp.
func (o *oauthUserApp) CreateUser(ctx context.Context, provider providerv1.OAuthProvider, accessToken string) (string, error) {
	var g errgroup.Group

	userID := uuid.New()

	// Initialize rollback functions
	rollbackFuncs, fail := util.NewRollback()
	rollbackFuncs = append(rollbackFuncs, func() {
		o.userEventService.UserCreationFailed(userID)
	})

	// Generate a sync code for the user
	syncCode, err := o.codeGenerator.GenerateSyncCode()
	if err != nil {
		return "", err
	}

	// Create the user in the repository
	g.Go(func() error {
		if _, err := o.repository.CreateUser(userID, syncCode); err != nil {
			return err
		}
		return nil
	})
	// Enroll the user with OAuth user service and create a profile
	g.Go(func() error {
		if _, err := o.oauthUserService.EnrollUser(ctx, userID, accessToken, provider, true, true); err != nil {
			return err
		}
		return nil
	})
	// Create a profile for the user
	g.Go(func() error {
		if _, err := o.profileService.CreateProfile(ctx, userID); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fail()
		return "", err
	}

	return userID.String(), nil
}

// // SyncUser implements OAuthUserApp.
// func (o *oauthUserApp) SyncUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, accessToken string) error {
// 	panic("unimplemented")
// }

// NewOAuthUserApp creates a new instance of oauthUserApp.
func NewOAuthUserApp(
	repository repodomain.UserRepository,
	oauthUserService oauthuserdomain.OAuthUserService,
	profileService profiledomain.ProfileService,
	userEventService usereventdomain.UserEventService,
	codeGenerator util.RandomCodeGenerator,
) OAuthUserApp {
	return &oauthUserApp{
		repository:       repository,
		oauthUserService: oauthUserService,
		profileService:   profileService,
		userEventService: userEventService,
		codeGenerator:    codeGenerator,
	}
}
