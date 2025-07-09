package oauthlogin

import (
	logindto "mandacode.com/accounts/web-auth/internal/app/login/dto"
	"mandacode.com/accounts/web-auth/internal/domain/model/provider"
)

type OAuthLoginApp interface {
	// GetLoginURL returns the URL to redirect the user for OAuth login.
	//
	// Parameters:
	//   - provider: The OAuth provider (e.g., "google", "naver").
	//
	// Returns:
	//   - A string representing the URL to redirect the user for OAuth login.
	//	 - An error if the URL generation fails, otherwise nil.
	GetLoginURL(provider provider.Provider) (string, error)

	// Login performs the OAuth login operation using the provided authorization code.
	//
	// Parameters:
	//   - provider: The OAuth provider (e.g., "google", "naver").
	//   - code: The authorization code received from the OAuth provider.
	//
	// Returns:
	//   - A pointer to a LoginToken containing the access and refresh tokens if the login is successful.
	//   - An error if the login operation fails, otherwise nil.
	Login(provider provider.Provider, code string) (*logindto.LoginToken, error)
}
