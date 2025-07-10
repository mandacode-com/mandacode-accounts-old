package userapp

import (
	"context"

	"github.com/google/uuid"
	localuserdomain "mandacode.com/accounts/user/internal/domain/port/auth/local"
	oauthuserdomain "mandacode.com/accounts/user/internal/domain/port/auth/oauth"
	profiledomain "mandacode.com/accounts/user/internal/domain/port/profile"
	repodomain "mandacode.com/accounts/user/internal/domain/port/repository"
	usereventdomain "mandacode.com/accounts/user/internal/domain/port/user_event"
	"mandacode.com/accounts/user/internal/util"
)

type userApp struct {
	repository       repodomain.UserRepository
	localUserService localuserdomain.LocalUserService
	oauthUserService oauthuserdomain.OAuthUserService
	profileService   profiledomain.ProfileService
	userEventService usereventdomain.UserEventService
	codeGenerator    util.RandomCodeGenerator
}

// ArchiveUser implements UserApp.
func (u *userApp) ArchiveUser(ctx context.Context, userID uuid.UUID) error {
	// Generate a sync code for the user
	syncCode, err := u.codeGenerator.GenerateSyncCode()
	if err != nil {
		return err
	}

	// Archive the user in the repository
	if err := u.repository.ArchiveUser(userID, syncCode); err != nil {
		return err
	}

	if err := u.userEventService.ArchiveUser(userID); err != nil {
		return err
	}

	return nil
}

// DeleteUser implements UserApp.
func (u *userApp) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := u.repository.DeleteUser(userID); err != nil {
		return err
	}

	// Notify the user event service about the deletion
	if err := u.userEventService.DeleteUser(userID); err != nil {
		return err
	}

	return nil
}

// UpdateActiveStatus implements UserApp.
func (u *userApp) UpdateActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) error {
	if _, err := u.localUserService.UpdateActiveStatus(ctx, userID, isActive); err != nil {
		return err
	}

	return nil
}

// UpdateVerifiedStatus implements UserApp.
func (u *userApp) UpdateVerifiedStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	if _, err := u.localUserService.UpdateVerifiedStatus(ctx, userID, isVerified); err != nil {
		return err
	}

	return nil
}

// NewUserApp creates a new instance of userApp with the provided dependencies.
func NewUserApp(
	repository repodomain.UserRepository,
	localUserService localuserdomain.LocalUserService,
	oauthUserService oauthuserdomain.OAuthUserService,
	profileService profiledomain.ProfileService,
	userEventService usereventdomain.UserEventService,
	codeGenerator util.RandomCodeGenerator,
) UserApp {
	return &userApp{
		repository:       repository,
		localUserService: localUserService,
		oauthUserService: oauthUserService,
		profileService:   profileService,
		userEventService: userEventService,
		codeGenerator:    codeGenerator,
	}
}
