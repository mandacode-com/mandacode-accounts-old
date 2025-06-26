package authdomain

import (
	"context"

	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/domain/dto"
)

type OAuthAuthService interface {
	// LoginOAuthUser logs in a user using OAuth provider and code.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - provider: The OAuth provider (e.g., Google, Facebook).
	//   - providerID: The unique identifier provided by the OAuth provider.
	//
	// Returns:
	//   - user: The OAuth user if login is successful.
	//   - error: An error if the login fails, otherwise nil.
	LoginOAuthUser(ctx context.Context, provider oauthuser.Provider, providerID string) (*dto.OAuthUser, error)
}
