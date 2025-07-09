package tokenprovider

import (
	"context"

	tokendomain "mandacode.com/accounts/auth/internal/domain/token"
	tokenv1 "github.com/mandacode-com/accounts-proto/token/v1"
)

type TokenProvider struct {
	client tokenv1.TokenServiceClient
}

func NewTokenProvider(client tokenv1.TokenServiceClient) tokendomain.TokenProvider {
	return &TokenProvider{client: client}
}

func (s *TokenProvider) GenerateAccessToken(ctx context.Context, userID string) (string, int64, error) {
	resp, err := s.client.GenerateAccessToken(ctx, &tokenv1.GenerateAccessTokenRequest{UserId: userID})
	if err != nil {
		return "", 0, ErrTokenGenerationFailed
	}
	return resp.Token, resp.ExpiresAt, nil
}

func (s *TokenProvider) VerifyAccessToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := s.client.VerifyAccessToken(ctx, &tokenv1.VerifyAccessTokenRequest{Token: token})
	if err != nil {
		return false, nil, ErrTokenVerificationFailed
	}
	return resp.Valid, resp.UserId, nil
}

func (s *TokenProvider) GenerateRefreshToken(ctx context.Context, userID string) (string, int64, error) {
	resp, err := s.client.GenerateRefreshToken(ctx, &tokenv1.GenerateRefreshTokenRequest{UserId: userID})
	if err != nil {
		return "", 0, ErrTokenGenerationFailed
	}
	return resp.Token, resp.ExpiresAt, nil
}

func (s *TokenProvider) VerifyRefreshToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := s.client.VerifyRefreshToken(ctx, &tokenv1.VerifyRefreshTokenRequest{Token: token})
	if err != nil {
		return false, nil, ErrTokenVerificationFailed
	}
	return resp.Valid, resp.UserId, nil
}

func (s *TokenProvider) GenerateEmailVerificationToken(ctx context.Context, userID, email, code string) (string, int64, error) {
	resp, err := s.client.GenerateEmailVerificationToken(ctx, &tokenv1.GenerateEmailVerificationTokenRequest{
		UserId: userID,
		Email:  email,
		Code:   code,
	})
	if err != nil {
		return "", 0, ErrTokenGenerationFailed
	}
	return resp.Token, resp.ExpiresAt, nil
}

func (s *TokenProvider) VerifyEmailVerificationToken(ctx context.Context, token string) (bool, *string, *string, *string, error) {
	resp, err := s.client.VerifyEmailVerificationToken(ctx, &tokenv1.VerifyEmailVerificationTokenRequest{Token: token})
	if err != nil {
		return false, nil, nil, nil, ErrTokenVerificationFailed
	}
	return resp.Valid, resp.UserId, resp.Email, resp.Code, nil
}
