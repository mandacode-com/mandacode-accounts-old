package oauthuser

import (
	"context"
	"github.com/google/uuid"
	oauthuserv1 "mandacode.com/accounts/proto/auth/user/oauth/v1"
	providerv1 "mandacode.com/accounts/proto/common/provider/v1"
	oauthuserdomain "mandacode.com/accounts/user/internal/domain/port/auth/oauth"
)

type OAuthUserService struct {
	client oauthuserv1.OAuthUserServiceClient
}

// UpdateActiveStatus implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) UpdateActiveStatus(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, isActive bool) (*oauthuserv1.UpdateActiveStatusResponse, error) {
	resp, err := o.client.UpdateActiveStatus(ctx, &oauthuserv1.UpdateActiveStatusRequest{
		UserId:   userID.String(),
		Provider: provider,
		IsActive: isActive,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateVerifiedStatus implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) UpdateVerifiedStatus(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, isVerified bool) (*oauthuserv1.UpdateVerifiedStatusResponse, error) {
	resp, err := o.client.UpdateVerifiedStatus(ctx, &oauthuserv1.UpdateVerifiedStatusRequest{
		UserId:    userID.String(),
		Provider:  provider,
		IsVerified: isVerified,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteAllProviders implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) DeleteAllProviders(ctx context.Context, userID uuid.UUID) (*oauthuserv1.DeleteAllProvidersResponse, error) {
	resp, err := o.client.DeleteAllProviders(ctx, &oauthuserv1.DeleteAllProvidersRequest{UserId: userID.String()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteUser implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) DeleteUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider) (*oauthuserv1.DeleteUserResponse, error) {
	resp, err := o.client.DeleteUser(ctx, &oauthuserv1.DeleteUserRequest{
		UserId:   userID.String(),
		Provider: provider,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// EnrollUser implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) EnrollUser(ctx context.Context, userID uuid.UUID, accessToken string, provider providerv1.OAuthProvider, isActive bool, isVerified bool) (*oauthuserv1.EnrollUserResponse, error) {
	resp, err := o.client.EnrollUser(ctx, &oauthuserv1.EnrollUserRequest{
		UserId:      userID.String(),
		AccessToken: accessToken,
		Provider:    provider,
		IsActive:    &isActive,
		IsVerified:  &isVerified,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUser implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) GetUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider) (*oauthuserv1.GetUserResponse, error) {
	resp, err := o.client.GetUser(ctx, &oauthuserv1.GetUserRequest{
		UserId:   userID.String(),
		Provider: provider,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SyncUser implements oauthuserdomain.OAuthUserService.
func (o *OAuthUserService) SyncUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, accessToken string) (*oauthuserv1.SyncUserResponse, error) {
	resp, err := o.client.SyncUser(ctx, &oauthuserv1.SyncUserRequest{
		UserId:      userID.String(),
		Provider:    provider,
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// NewOAuthUserService creates a new instance of OAuthUserService.
func NewOAuthUserService(client oauthuserv1.OAuthUserServiceClient) oauthuserdomain.OAuthUserService {
	return &OAuthUserService{client: client}
}
