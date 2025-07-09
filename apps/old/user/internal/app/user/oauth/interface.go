package oauthuserapp

import (
	"context"

	providerv1 "github.com/mandacode-com/accounts-proto/common/provider/v1"
)

type OAuthUserApp interface {
	// CreateUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//	 - provider: The OAuth provider (e.g., "google", "github").
	//   - accessToken: The access token received from the OAuth provider.
	CreateUser(ctx context.Context, provider providerv1.OAuthProvider, accessToken string) (string, error)

	// SyncUser
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to sync.
	//   - provider: The OAuth provider (e.g., "google", "github").
	//   - accessToken: The access token received from the OAuth provider.
	// SyncUser(ctx context.Context, userID uuid.UUID, provider providerv1.OAuthProvider, accessToken string) error
}
