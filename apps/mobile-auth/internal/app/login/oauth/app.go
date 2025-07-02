package oauthlogin

import (
	logindto "mandacode.com/accounts/mobile-auth/internal/app/login/dto"
	authdomain "mandacode.com/accounts/mobile-auth/internal/domain/auth"
	"mandacode.com/accounts/mobile-auth/internal/domain/model/provider"
)

type oauthLoginApp struct {
	authenticator authdomain.Authenticator
}

// Login implements OAuthLoginApp.
func (o *oauthLoginApp) Login(provider provider.Provider, accessToken string) (*logindto.LoginToken, error) {
	loginToken, err := o.authenticator.OAuthLogin(provider, accessToken)
	if err != nil {
		return nil, err
	}

	return &logindto.LoginToken{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}, nil
}

// NewOAuthLoginApp creates a new instance of OAuthLoginApp with the provided authenticator.
func NewOAuthLoginApp(authenticator authdomain.Authenticator) OAuthLoginApp {
	return &oauthLoginApp{
		authenticator: authenticator,
	}
}
