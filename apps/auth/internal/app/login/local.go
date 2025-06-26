package login

import (
	"context"

	logindomain "mandacode.com/accounts/auth/internal/domain/service/login"
	tokendomain "mandacode.com/accounts/auth/internal/domain/service/token"
)

type LocalLoginApp struct {
	tokenService     tokendomain.TokenService
	localAuthService logindomain.LocalLoginService
}

func NewLocalLoginApp(tokenService tokendomain.TokenService, localAuthService logindomain.LocalLoginService) *LocalLoginApp {
	return &LocalLoginApp{
		tokenService:     tokenService,
		localAuthService: localAuthService,
	}
}

func (a *LocalLoginApp) LoginLocalUser(ctx context.Context, email, password string) (*string, *string, *string, error) {
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
