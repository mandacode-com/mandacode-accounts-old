package oauthlogin

import (
	"errors"

	logindto "mandacode.com/accounts/web-auth/internal/app/login/dto"
	authdomain "mandacode.com/accounts/web-auth/internal/domain/auth"
	"mandacode.com/accounts/web-auth/internal/domain/model/provider"
	oauthcodedomain "mandacode.com/accounts/web-auth/internal/domain/oauthcode"
)

type oauthLoginApp struct {
	oauthCodeMap  map[provider.Provider]oauthcodedomain.OAuthCode
	authenticator authdomain.Authenticator
}

func NewOAuthLoginApp(oauthCodeMap map[provider.Provider]oauthcodedomain.OAuthCode, authenticator authdomain.Authenticator) OAuthLoginApp {
	return &oauthLoginApp{
		oauthCodeMap:  oauthCodeMap,
		authenticator: authenticator,
	}
}

func (app *oauthLoginApp) GetLoginURL(provider provider.Provider) (string, error) {
	oauthCode, exists := app.oauthCodeMap[provider]
	if !exists {
		return "", errors.New("unsupported OAuth provider: " + provider.String())
	}

	loginURL, err := oauthCode.GetLoginURL()
	if err != nil {
		return "", err
	}

	return loginURL, nil
}

func (app *oauthLoginApp) Login(provider provider.Provider, code string) (*logindto.LoginToken, error) {
	oauthCode, exists := app.oauthCodeMap[provider]
	if !exists {
		return nil, errors.New("unsupported OAuth provider: " + provider.String())
	}

	accessToken, err := oauthCode.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	// Create a LoginToken with the access token
	loginTokenRes, err := app.authenticator.OAuthLogin(provider, accessToken)
	if err != nil {
		return nil, err
	}

	// Convert the LoginResponse to LoginToken
	loginToken := &logindto.LoginToken{
		AccessToken:  loginTokenRes.AccessToken,
		RefreshToken: loginTokenRes.RefreshToken,
	}

	return loginToken, nil
}
