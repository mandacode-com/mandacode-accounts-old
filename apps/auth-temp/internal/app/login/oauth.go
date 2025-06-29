package login

import (
	"context"
	"errors"

	"mandacode.com/accounts/auth/ent/oauthuser"
	logindomain "mandacode.com/accounts/auth/internal/domain/service/login"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	tokendomain "mandacode.com/accounts/auth/internal/domain/service/token"
	"mandacode.com/accounts/auth/internal/util"
	providerv1 "mandacode.com/accounts/proto/common/provider/v1"
)

type OAuthLoginApp struct {
	providers        *map[oauthuser.Provider]oauthdomain.OAuthService
	tokenService     tokendomain.TokenService
	oauthAuthService logindomain.OAuthLoginService
}

func NewOAuthLoginApp(
	providers *map[oauthuser.Provider]oauthdomain.OAuthService,
	tokenService tokendomain.TokenService,
	oauthAuthService logindomain.OAuthLoginService) *OAuthLoginApp {
	return &OAuthLoginApp{
		providers:        providers,
		tokenService:     tokenService,
		oauthAuthService: oauthAuthService,
	}
}

func (a *OAuthLoginApp) LoginOAuthUser(ctx context.Context, provider providerv1.OAuthProvider, oauthAccessToken string) (*string, *string, *string, error) {
	providerEnum, err := util.FromProtoToProvider(provider)
	if err != nil {
		return nil, nil, nil, err
	}
	oauthService, exists := (*a.providers)[providerEnum]
	if !exists {
		return nil, nil, nil, errors.New("unsupported provider")
	}
	// Get user info from the OAuth oauthAuthService
	userInfo, err := oauthService.GetUserInfo(oauthAccessToken)
	if err != nil {
		return nil, nil, nil, err
	}
	if userInfo.ProviderID == "" {
		return nil, nil, nil, errors.New("user info does not contain a valid provider ID")
	}

	user, err := a.oauthAuthService.LoginOAuthUser(ctx, providerEnum, userInfo.ProviderID)
	if err != nil {
		return nil, nil, nil, err
	}

	userID := user.ID.String()

	// Generate access and refresh tokens for the LoginOAuthUser
	accessToken, _, err := a.tokenService.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	refreshToken, _, err := a.tokenService.GenerateRefreshToken(ctx, user.ID.String())
	if err != nil {
		return nil, nil, nil, err
	}

	return &userID, &accessToken, &refreshToken, nil
}
