package localuser

import (
	"context"
	"github.com/google/uuid"
	localuserv1 "mandacode.com/accounts/proto/auth/user/local/v1"
	localuserdomain "mandacode.com/accounts/user/internal/domain/port/auth/local"
)

type LocalUserService struct {
	client localuserv1.LocalUserServiceClient
}

// DeleteUser implements localuserdomain.LocalUserService.
func (l *LocalUserService) DeleteUser(ctx context.Context, userID uuid.UUID) (*localuserv1.DeleteUserResponse, error) {
	resp, err := l.client.DeleteUser(ctx, &localuserv1.DeleteUserRequest{UserId: userID.String()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// EnrollUser implements localuserdomain.LocalUserService.
func (l *LocalUserService) EnrollUser(ctx context.Context, userID uuid.UUID, email string, password string, isActive bool, isVerified bool) (*localuserv1.EnrollUserResponse, error) {
	resp, err := l.client.EnrollUser(ctx, &localuserv1.EnrollUserRequest{
		UserId:     userID.String(),
		Email:      email,
		Password:   password,
		IsActive:   &isActive,
		IsVerified: &isVerified,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUser implements localuserdomain.LocalUserService.
func (l *LocalUserService) GetUser(ctx context.Context, userID uuid.UUID) (*localuserv1.GetUserResponse, error) {
	resp, err := l.client.GetUser(ctx, &localuserv1.GetUserRequest{UserId: userID.String()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateActiveStatus implements localuserdomain.LocalUserService.
func (l *LocalUserService) UpdateActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) (*localuserv1.UpdateActiveStatusResponse, error) {
	resp, err := l.client.UpdateActiveStatus(ctx, &localuserv1.UpdateActiveStatusRequest{
		UserId:   userID.String(),
		IsActive: isActive,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateEmail implements localuserdomain.LocalUserService.
func (l *LocalUserService) UpdateEmail(ctx context.Context, userID uuid.UUID, email string) (*localuserv1.UpdateEmailResponse, error) {
	resp, err := l.client.UpdateEmail(ctx, &localuserv1.UpdateEmailRequest{
		UserId: userID.String(),
		Email:  email,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdatePassword implements localuserdomain.LocalUserService.
func (l *LocalUserService) UpdatePassword(ctx context.Context, userID uuid.UUID, currentPassword string, newPassword string) (*localuserv1.UpdatePasswordResponse, error) {
	resp, err := l.client.UpdatePassword(ctx, &localuserv1.UpdatePasswordRequest{
		UserId:          userID.String(),
		CurrentPassword: currentPassword,
		NewPassword:     newPassword,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateVerifiedStatus implements localuserdomain.LocalUserService.
func (l *LocalUserService) UpdateVerifiedStatus(ctx context.Context, userID uuid.UUID, isVerified bool) (*localuserv1.UpdateVerifiedStatusResponse, error) {
	resp, err := l.client.UpdateVerifiedStatus(ctx, &localuserv1.UpdateVerifiedStatusRequest{
		UserId:     userID.String(),
		IsVerified: isVerified,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewLocalUserService creates a new instance of LocalUserService with the provided client.
func NewLocalUserService(client localuserv1.LocalUserServiceClient) localuserdomain.LocalUserService {
	return &LocalUserService{client: client}
}
