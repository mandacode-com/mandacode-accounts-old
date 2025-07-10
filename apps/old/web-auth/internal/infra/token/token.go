package tokenmgr

import (
	"context"

	"github.com/google/uuid"
	tokenv1 "github.com/mandacode-com/accounts-proto/auth/token/v1"
	tokenmodel "mandacode.com/accounts/web-auth/internal/domain/model/token"
	tokenmgrdomain "mandacode.com/accounts/web-auth/internal/domain/token"
)

type TokenManager struct {
	tokenClient tokenv1.TokenServiceClient
}

// RefreshToken implements tokenmgrdomain.TokenManager.
func (t *TokenManager) RefreshToken(ctx context.Context, refreshToken string) (result *tokenmodel.RefresedToken, err error) {
	resp, err := t.tokenClient.RefreshToken(ctx, &tokenv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}

	return &tokenmodel.RefresedToken{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// VerifyAccessToken implements tokenmgrdomain.TokenManager.
func (t *TokenManager) VerifyAccessToken(ctx context.Context, accessToken string) (result *tokenmodel.VerifyAccessTokenResult, err error) {
	resp, err := t.tokenClient.VerifyToken(ctx, &tokenv1.VerifyTokenRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}

	if resp.Valid {
		userUUID, err := uuid.Parse(resp.UserId)
		if err != nil {
			return nil, err
		}
		return &tokenmodel.VerifyAccessTokenResult{
			UserID: &userUUID,
			Valid:  true,
		}, nil
	}
	return &tokenmodel.VerifyAccessTokenResult{
		Valid: false,
		UserID: nil,
	}, nil
}

func NewTokenManager(tokenClient tokenv1.TokenServiceClient) tokenmgrdomain.TokenManager {
	return &TokenManager{
		tokenClient: tokenClient,
	}
}
