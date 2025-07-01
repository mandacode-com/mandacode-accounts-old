package token

import (
	"context"

	tokenmodel "mandacode.com/accounts/web-auth/internal/domain/model/token"
	tokenmgrdomain "mandacode.com/accounts/web-auth/internal/domain/token"
)

type tokenApp struct {
	tokenManager tokenmgrdomain.TokenManager
}

// RefreshToken implements TokenApp.
func (t *tokenApp) RefreshToken(ctx context.Context, refreshToken string) (result *tokenmodel.RefresedToken, err error) {
	return t.tokenManager.RefreshToken(ctx, refreshToken)
}

// VerifyAccessToken implements TokenApp.
func (t *tokenApp) VerifyAccessToken(ctx context.Context, accessToken string) (result *tokenmodel.VerifyAccessTokenResult, err error) {
	return t.tokenManager.VerifyAccessToken(ctx, accessToken)
}

// RefreshToken implements TokenApp.
func NewTokenApp(tokenManager tokenmgrdomain.TokenManager) TokenApp {
	return &tokenApp{
		tokenManager: tokenManager,
	}
}
