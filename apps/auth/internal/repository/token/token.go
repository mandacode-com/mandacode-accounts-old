package tokenrepo

import (
	"context"

	"github.com/google/uuid"
	tokenv1 "github.com/mandacode-com/accounts-proto/go/token/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokendomain "mandacode.com/accounts/auth/internal/domain/repository/token"
)

type TokenRepository struct {
	client tokenv1.TokenServiceClient
}

// GenerateAccessToken implements tokendomain.TokenRepository.
func (t *TokenRepository) GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, int64, error) {
	resp, err := t.client.GenerateAccessToken(ctx, &tokenv1.GenerateAccessTokenRequest{UserId: userID.String()})
	if err != nil {
		joinedErr := errors.Join(err, "Failed to generate access token")
		return "", 0, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}
	return resp.Token, resp.ExpiresAt, nil
}

// GenerateEmailVerificationToken implements tokendomain.TokenRepository.
func (t *TokenRepository) GenerateEmailVerificationToken(ctx context.Context, userID uuid.UUID, email string, code string) (string, int64, error) {
	resp, err := t.client.GenerateEmailVerificationToken(ctx, &tokenv1.GenerateEmailVerificationTokenRequest{
		UserId: userID.String(),
		Email:  email,
		Code:   code,
	})
	if err != nil {
		joinedErr := errors.Join(err, "Failed to generate email verification token")
		return "", 0, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}
	return resp.Token, resp.ExpiresAt, nil
}

// GenerateRefreshToken implements tokendomain.TokenRepository.
func (t *TokenRepository) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, int64, error) {
	resp, err := t.client.GenerateRefreshToken(ctx, &tokenv1.GenerateRefreshTokenRequest{UserId: userID.String()})
	if err != nil {
		joinedErr := errors.Join(err, "Failed to generate refresh token")
		return "", 0, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}
	return resp.Token, resp.ExpiresAt, nil
}

// VerifyAccessToken implements tokendomain.TokenRepository.
func (t *TokenRepository) VerifyAccessToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := t.client.VerifyAccessToken(ctx, &tokenv1.VerifyAccessTokenRequest{Token: token})
	if err != nil {
		joinErr := errors.Join(err, "Failed to verify access token")
		return false, nil, errors.Upgrade(joinErr, errcode.ErrInternalFailure, "TokenVerificationFailed")
	}
	return resp.Valid, resp.UserId, nil
}

// VerifyEmailVerificationToken implements tokendomain.TokenRepository.
func (t *TokenRepository) VerifyEmailVerificationToken(ctx context.Context, token string) (*tokendomain.EmailVerificationResult, error) {
	resp, err := t.client.VerifyEmailVerificationToken(ctx, &tokenv1.VerifyEmailVerificationTokenRequest{Token: token})
	if err != nil {
		joinErr := errors.Join(err, "Failed to verify email verification token")
		return nil, errors.Upgrade(joinErr, errcode.ErrInternalFailure, "TokenVerificationFailed")
	}

	if !resp.Valid {
		return nil, errors.New("Invalid email verification token", "InvalidToken", errcode.ErrUnauthorized)
	}

	userUUID, err := uuid.Parse(*resp.UserId)
	if err != nil {
		joinErr := errors.Join(err, "Failed to parse user ID from email verification token response")
		return nil, errors.Upgrade(joinErr, errcode.ErrInternalFailure, "TokenVerificationFailed")
	}
	data := &tokendomain.EmailVerificationResult{
		Valid:  resp.Valid,
		UserID: userUUID,
		Email:  *resp.Email,
		Code:   *resp.Code,
	}
	return data, nil
}

// VerifyRefreshToken implements tokendomain.TokenRepository.
func (t *TokenRepository) VerifyRefreshToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := t.client.VerifyRefreshToken(ctx, &tokenv1.VerifyRefreshTokenRequest{Token: token})
	if err != nil {
		joinErr := errors.Join(err, "Failed to verify refresh token")
		return false, nil, errors.Upgrade(joinErr, errcode.ErrInternalFailure, "TokenVerificationFailed")
	}
	return resp.Valid, resp.UserId, nil
}

func NewTokenRepository(client tokenv1.TokenServiceClient) tokendomain.TokenRepository {
	return &TokenRepository{client: client}
}
