package auth

import (
	"context"

	localloginv1 "github.com/mandacode-com/accounts-proto/auth/login/local/v1"
	oauthloginv1 "github.com/mandacode-com/accounts-proto/auth/login/oauth/v1"
	authdomain "mandacode.com/accounts/web-auth/internal/domain/auth"
	"mandacode.com/accounts/web-auth/internal/domain/model/provider"
	authresdto "mandacode.com/accounts/web-auth/internal/infra/auth/dto/response"
)

type Authenticator struct {
	localLoginClient localloginv1.LocalLoginServiceClient
	oauthLoginClient oauthloginv1.OAuthLoginServiceClient
}

func NewAuthenticator(localClient localloginv1.LocalLoginServiceClient, oauthClient oauthloginv1.OAuthLoginServiceClient) authdomain.Authenticator {
	return &Authenticator{
		localLoginClient: localClient,
		oauthLoginClient: oauthClient,
	}
}

func (a *Authenticator) LocalLogin(email, password string) (*authresdto.LoginResponse, error) {
	resp, err := a.localLoginClient.Login(context.Background(), &localloginv1.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &authresdto.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (a *Authenticator) OAuthLogin(provider provider.Provider, accessToken string) (*authresdto.LoginResponse, error) {
	resp, err := a.oauthLoginClient.Login(context.Background(), &oauthloginv1.LoginRequest{
		Provider:    provider.ToProto(),
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	return &authresdto.LoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
