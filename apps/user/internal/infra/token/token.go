package token

import (
	"context"

	"github.com/google/uuid"
	tokenv1 "mandacode.com/accounts/proto/token/v1"
	tokendomain "mandacode.com/accounts/user/internal/domain/port/token"
)

type TokenService struct {
	client tokenv1.TokenServiceClient
}

// GenerateEmailVerificationToken implements tokendomain.TokenService.
func (t *TokenService) GenerateEmailVerificationToken(ctx context.Context, userID uuid.UUID, email string, code string) (string, int64, error) {
	resp, err := t.client.GenerateEmailVerificationToken(ctx, &tokenv1.GenerateEmailVerificationTokenRequest{
		UserId: userID.String(),
		Email:  email,
		Code:   code,
	})
	if err != nil {
		return "", 0, err
	}
	return resp.Token, resp.ExpiresAt, nil
}

// VerifyEmailVerificationToken implements tokendomain.TokenService.
func (t *TokenService) VerifyEmailVerificationToken(ctx context.Context, token string) (bool, *string, *string, *string, error) {
	resp, err := t.client.VerifyEmailVerificationToken(ctx, &tokenv1.VerifyEmailVerificationTokenRequest{Token: token})
	if err != nil {
		return false, nil, nil, nil, err
	}
	if !resp.Valid {
		return false, nil, nil, nil, nil
	}
	return resp.Valid, resp.UserId, resp.Email, resp.Code, nil
}

func NewTokenService(client tokenv1.TokenServiceClient) tokendomain.TokenService {
	return &TokenService{client: client}
}
