package token

import (
	"context"
	"errors"

	tokendto "mandacode.com/accounts/auth/internal/app/token/dto"
	tokendomain "mandacode.com/accounts/auth/internal/domain/token"
)

type tokenApp struct {
	tokenProvider tokendomain.TokenProvider
}

func NewTokenApp(tokenProvider tokendomain.TokenProvider) TokenApp {
	return &tokenApp{
		tokenProvider: tokenProvider,
	}
}

// RefreshToken implements TokenApp.
func (t *tokenApp) RefreshToken(ctx context.Context, refreshToken string) (*tokendto.NewToken, error) {
	valid, userID, err := t.tokenProvider.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token: " + err.Error())
	}
	if !valid {
		return nil, errors.New("refresh token is not valid")
	}

	// Generate new access and refresh tokens
	newAccessToken, _, err := t.tokenProvider.GenerateAccessToken(ctx, *userID)
	if err != nil {
		return nil, errors.New("failed to generate new access token: " + err.Error())
	}

	newRefreshToken, _, err := t.tokenProvider.GenerateRefreshToken(ctx, *userID)
	if err != nil {
		return nil, errors.New("failed to generate new refresh token: " + err.Error())
	}

	return &tokendto.NewToken{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// VerifyToken implements TokenApp.
func (t *tokenApp) VerifyToken(ctx context.Context, accessToken string) (*tokendto.VerifyTokenResult, error) {
	valid, userID, err := t.tokenProvider.VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return nil, errors.New("invalid access token: " + err.Error())
	}
	if !valid {
		return nil, errors.New("access token is not valid")
	}
	return &tokendto.VerifyTokenResult{
		UserID: *userID,
		Valid:  valid,
	}, nil
}
