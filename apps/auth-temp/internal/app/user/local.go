package user

import (
	"context"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/internal/domain/dto"
	userdomain "mandacode.com/accounts/auth/internal/domain/service/user"
)

type LocalUserApp struct {
	localUserService userdomain.LocalUserService
}

func NewLocalUserApp(localUserService userdomain.LocalUserService) *LocalUserApp {
	return &LocalUserApp{
		localUserService: localUserService,
	}
}

// Create a new local user
func (a *LocalUserApp) CreateUser(ctx context.Context, userID string, email string, password string, isActive *bool, isVerified *bool) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.CreateUser(userUUID, email, password, isActive, isVerified)
}

// Get user by id
func (a *LocalUserApp) GetUser(ctx context.Context, userID string) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.GetUserByID(userUUID)
}

// Delete user
func (a *LocalUserApp) DeleteUser(ctx context.Context, userID string) (*dto.LocalDeletedUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.DeleteUser(userUUID)
}

// Update email
func (a *LocalUserApp) UpdateEmail(ctx context.Context, userID string, email string) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.UpdateEmail(userUUID, email)
}

// Update password
func (a *LocalUserApp) UpdatePassword(ctx context.Context, userID string, currentPassword string, newPassword string) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.UpdatePassword(userUUID, currentPassword, newPassword)
}

// Update active status
func (a *LocalUserApp) UpdateActiveStatus(ctx context.Context, userID string, isActive bool) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.UpdateActiveStatus(userUUID, isActive)
}

// Update verified status
func (a *LocalUserApp) UpdateVerifiedStatus(ctx context.Context, userID string, isVerified bool) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return a.localUserService.UpdateVerifiedStatus(userUUID, isVerified)
}
