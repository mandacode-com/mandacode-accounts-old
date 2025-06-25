package auth

import (
	"context"

	authdomain "mandacode.com/accounts/auth/internal/domain/service/auth"
	tokendomain "mandacode.com/accounts/auth/internal/domain/service/token"
)

type LocalAuthApp struct {
	tokenService     tokendomain.TokenService
	localAuthService authdomain.LocalAuthService
}

func NewLocalAuthApp(tokenService tokendomain.TokenService, localAuthService authdomain.LocalAuthService) *LocalAuthApp {
	return &LocalAuthApp{
		tokenService:     tokenService,
		localAuthService: localAuthService,
	}
}

func (a *LocalAuthApp) LoginLocalUser(ctx context.Context, email, password string) (*string, *string, *string, error) {
	user, err := a.localAuthService.LoginLocalUser(ctx, email, password)
	if err != nil {
		return nil, nil, nil, err
	}

	userID := user.ID.String()

	// Generate access and refresh tokens for the LoginLocalUser
	accessToken, _, err := a.tokenService.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	refreshToken, _, err := a.tokenService.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	return &userID, &accessToken, &refreshToken, nil
}
