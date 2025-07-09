package oauthlogin

import (
	"context"

	"mandacode.com/accounts/auth/ent/oauthuser"
	logindto "mandacode.com/accounts/auth/internal/app/login/dto"
)

type OAuthLoginApp interface {
	// Login authenticates a user using OAuth with the specified provider and access token.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - provider: The OAuth provider to use for authentication.
	// - oauthAccessToken: The OAuth access token for the user.
	//
	// Returns:
	// - *logindto.LoginToken: The login token containing access and refresh tokens.
	// - error: An error if the login fails, otherwise nil.
	Login(ctx context.Context, provider oauthuser.Provider, oauthAccessToken string) (*logindto.LoginToken, error)
}
