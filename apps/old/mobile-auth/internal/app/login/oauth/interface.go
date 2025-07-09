package oauthlogin

import (
	logindto "mandacode.com/accounts/mobile-auth/internal/app/login/dto"
	"mandacode.com/accounts/mobile-auth/internal/domain/model/provider"
)

type OAuthLoginApp interface {
	// Login performs the OAuth login operation using the provided authorization code.
	//
	// Parameters:
	//   - provider: The OAuth provider (e.g., "google", "naver").
	//   - accessToken: The access token received from the OAuth provider.
	//
	// Returns:
	//   - A pointer to a LoginToken containing the access and refresh tokens if the login is successful.
	//   - An error if the login operation fails, otherwise nil.
	Login(provider provider.Provider, accessToken string) (*logindto.LoginToken, error)
}
