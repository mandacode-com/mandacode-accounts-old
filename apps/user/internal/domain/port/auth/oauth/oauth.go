package oauthuserdomain

import (
	"context"

	"github.com/google/uuid"
	oauthuserv1 "mandacode.com/accounts/proto/auth/user/oauth/v1"
	providerv1 "mandacode.com/accounts/proto/common/provider/v1"
)

type OAuthUserService interface {
	// GetUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to retrieve.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	GetUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider) (*oauthuserv1.GetUserResponse, error)

	// EnrollUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to enroll.
	//   - accessToken: The access token for the OAuth provider.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//   - isActive: Indicates if the user is active.
	//   - isVerified: Indicates if the user is verified.
	EnrollUser(ctx context.Context, userID uuid.UUID, accessToken string, provider providerv1.OAuthProvider, isActive bool, isVerified bool) (*oauthuserv1.EnrollUserResponse, error)

	// DeleteUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to be deleted.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	DeleteUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider) (*oauthuserv1.DeleteUserResponse, error)

	// DeleteAllProviders
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose OAuth providers are to be deleted.
	DeleteAllProviders(ctx context.Context, userID uuid.UUID) (*oauthuserv1.DeleteAllProvidersResponse, error)

	// SyncUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to sync.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//   - accessToken: The access token for the OAuth provider.
	SyncUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, accessToken string) (*oauthuserv1.SyncUserResponse, error)

	// UpdateActiveStatus
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose active status is to be updated.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//	 - isActive: Indicates if the user is active.
	UpdateActiveStatus(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, isActive bool) (*oauthuserv1.UpdateActiveStatusResponse, error)

	// UpdateVerifiedStatus
	//
	// Parameters:
	//	 - ctx: The context for the operation.
	//   - userID: The ID of the user whose verified status is to be updated.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//   - isVerified: Indicates if the user is verified.
	UpdateVerifiedStatus(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, isVerified bool) (*oauthuserv1.UpdateVerifiedStatusResponse, error)
}
