package authdomain

import (
	"mandacode.com/accounts/mobile-auth/internal/domain/model/provider"
	authresdto "mandacode.com/accounts/mobile-auth/internal/infra/auth/dto/response"
)

type Authenticator interface {
	// LocalLogin authenticates a user with a username and password.
	//
	// Parameters:
	//   - username: The username of the user.
	//   - password: The password of the user.
	//
	// Returns:
	//   - A pointer to a LoginResponse containing the access and refresh tokens.
	//   - An error if the authentication fails, otherwise nil.
	LocalLogin(username, password string) (*authresdto.LoginResponse, error)

	// OAuthLogin authenticates a user using OAuth with the specified provider and code.
	//
	// Parameters:
	//   - provider: The OAuth provider (e.g., "google", "naver").
	//   - accessToken: The access token received from the OAuth provider.
	//
	// Returns:
	//   - A pointer to a LoginResponse containing the access and refresh tokens.
	//   - An error if the authentication fails, otherwise nil.
	OAuthLogin(provider provider.Provider, accessToken string) (*authresdto.LoginResponse, error)
}
