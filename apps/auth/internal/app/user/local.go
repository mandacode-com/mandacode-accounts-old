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

func (a *LocalUserApp) CreateLocalUser(ctx context.Context, userID string, email string, password string, isActive *bool, isVerified *bool) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := a.localUserService.CreateLocalUser(userUUID, email, password, isActive, isVerified)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *LocalUserApp) DeleteLocalUser(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return a.localUserService.DeleteLocalUser(userUUID)
}

func (a *LocalUserApp) UpdateLocalUser(ctx context.Context, userID string, email *string, password *string, isActive *bool, isVerified *bool) (*dto.LocalUser, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	user, err := a.localUserService.UpdateLocalUser(userUUID, email, password, isActive, isVerified)
	if err != nil {
		return nil, err
	}

	return user, nil
}
