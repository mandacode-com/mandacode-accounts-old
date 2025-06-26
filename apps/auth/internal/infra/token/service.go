package token

import (
	"context"

	tokendomain "mandacode.com/accounts/auth/internal/domain/service/token"
	tokenv1 "mandacode.com/accounts/proto/token/v1"
)

type TokenService struct {
	client tokenv1.TokenServiceClient
}

func New(client tokenv1.TokenServiceClient) tokendomain.TokenService {
	return &TokenService{client: client}
}

func (s *TokenService) GenerateAccessToken(ctx context.Context, userID string) (string, int64, error) {
	resp, err := s.client.GenerateAccessToken(ctx, &tokenv1.GenerateAccessTokenRequest{UserId: userID})
	if err != nil {
		return "", 0, err
	}
	return resp.Token, resp.ExpiresAt, nil
}

func (s *TokenService) VerifyAccessToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := s.client.VerifyAccessToken(ctx, &tokenv1.VerifyAccessTokenRequest{Token: token})
	if err != nil {
		return false, nil, err
	}
	return resp.Valid, resp.UserId, nil
}

func (s *TokenService) GenerateRefreshToken(ctx context.Context, userID string) (string, int64, error) {
	resp, err := s.client.GenerateRefreshToken(ctx, &tokenv1.GenerateRefreshTokenRequest{UserId: userID})
	if err != nil {
		return "", 0, err
	}
	return resp.Token, resp.ExpiresAt, nil
}

func (s *TokenService) VerifyRefreshToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := s.client.VerifyRefreshToken(ctx, &tokenv1.VerifyRefreshTokenRequest{Token: token})
	if err != nil {
		return false, nil, err
	}
	return resp.Valid, resp.UserId, nil
}

func (s *TokenService) GenerateEmailVerificationToken(ctx context.Context, userID, email, code string) (string, int64, error) {
	resp, err := s.client.GenerateEmailVerificationToken(ctx, &tokenv1.GenerateEmailVerificationTokenRequest{
		UserId: userID,
		Email:  email,
		Code:   code,
	})
	if err != nil {
		return "", 0, err
	}
	return resp.Token, resp.ExpiresAt, nil
}

func (s *TokenService) VerifyEmailVerificationToken(ctx context.Context, token string) (bool, *string, *string, *string, error) {
	resp, err := s.client.VerifyEmailVerificationToken(ctx, &tokenv1.VerifyEmailVerificationTokenRequest{Token: token})
	if err != nil {
		return false, nil, nil, nil, err
	}
	return resp.Valid, resp.UserId, resp.Email, resp.Code, nil
}
