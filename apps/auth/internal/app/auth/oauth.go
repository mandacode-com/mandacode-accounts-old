package auth

import (
	"context"
	"errors"

	"mandacode.com/accounts/auth/ent/oauthuser"
	authdomain "mandacode.com/accounts/auth/internal/domain/service/auth"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	tokendomain "mandacode.com/accounts/auth/internal/domain/service/token"
	"mandacode.com/accounts/auth/internal/util"
	providerv1 "mandacode.com/accounts/auth/proto/common/provider/v1"
)

type OAuthAuthApp struct {
	providers        *map[oauthuser.Provider]oauthdomain.OAuthService
	tokenService     tokendomain.TokenService
	oauthAuthService authdomain.OAuthAuthService
}

func NewOAuthAuthApp(
	providers *map[oauthuser.Provider]oauthdomain.OAuthService,
	tokenService tokendomain.TokenService,
	oauthAuthService authdomain.OAuthAuthService) *OAuthAuthApp {
	return &OAuthAuthApp{
		providers:        providers,
		tokenService:     tokenService,
		oauthAuthService: oauthAuthService,
	}
}

func (a *OAuthAuthApp) LoginOAuthUser(ctx context.Context, provider providerv1.OAuthProvider, oauthAccessToken string) (*string, *string, *string, error) {
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
