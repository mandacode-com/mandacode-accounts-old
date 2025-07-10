package oauthlogin

import (
	"context"

	"mandacode.com/accounts/auth/ent/oauthuser"
	logindto "mandacode.com/accounts/auth/internal/app/login/dto"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/oauth"
	repodomain "mandacode.com/accounts/auth/internal/domain/repository"
	tokendomain "mandacode.com/accounts/auth/internal/domain/token"
	tokenprovider "mandacode.com/accounts/auth/internal/infra/token"
	protoutil "mandacode.com/accounts/auth/internal/util/proto"
)

type oauthLoginApp struct {
	tokenProvider  tokendomain.TokenProvider
	oauthProviders *map[oauthuser.Provider]oauthdomain.OAuthProvider
	oauthLoginRep  repodomain.OAuthUserRepository
}

func NewOAuthLoginApp(
	tokenProvider tokendomain.TokenProvider,
	oauthProviders *map[oauthuser.Provider]oauthdomain.OAuthProvider,
	oauthLoginRep repodomain.OAuthUserRepository,
) OAuthLoginApp {
	return &oauthLoginApp{
		tokenProvider:  tokenProvider,
		oauthProviders: oauthProviders,
		oauthLoginRep:  oauthLoginRep,
	}
}

func (a *oauthLoginApp) Login(ctx context.Context, provider oauthuser.Provider, oauthAccessToken string) (*logindto.LoginToken, error) {
	oauthProvider, ok := (*a.oauthProviders)[provider]
	if !ok {
		return nil, protoutil.ErrUnsupportedProvider
	}

	userInfo, err := oauthProvider.GetUserInfo(oauthAccessToken)
	if err != nil {
		return nil, err
	}

	user, err := a.oauthLoginRep.GetUserByProvider(provider, userInfo.ProviderID)
	if err != nil {
		return nil, err
	}

	userID := user.ID.String()

	// Generate access and refresh tokens for the OAuthLogin
	accessToken, _, err := a.tokenProvider.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, tokenprovider.ErrTokenGenerationFailed
	}
	refreshToken, _, err := a.tokenProvider.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return nil, tokenprovider.ErrTokenGenerationFailed
	}

	return &logindto.LoginToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
