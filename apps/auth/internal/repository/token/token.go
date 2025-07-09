package tokenrepo

import (
	"context"

	"github.com/google/uuid"
	tokenv1 "github.com/mandacode-com/accounts-proto/go/token/v1"
	"github.com/mandacode-com/golib/errors"
	"github.com/mandacode-com/golib/errors/errcode"
	tokenmodels "mandacode.com/accounts/auth/internal/models/token"
)

type TokenRepository struct {
	client tokenv1.TokenServiceClient
}

// GenerateAccessToken creates a new access token for the user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: The ID of the user for whom the access token is generated.
//
// Returns:
//   - token: The generated access token.
//   - expiresAt: The expiration time of the token in Unix timestamp format.
//   - error: An error if the token generation fails, otherwise nil.
func (t *TokenRepository) GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, int64, error) {
	resp, err := t.client.GenerateAccessToken(ctx, &tokenv1.GenerateAccessTokenRequest{UserId: userID.String()})
	if err != nil {
		joinedErr := errors.Join(err, "Failed to generate access token")
		return "", 0, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}
	return resp.Token, resp.ExpiresAt, nil
}

// GenerateEmailVerificationToken creates a new email verification token for the user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: The ID of the user for whom the email verification token is generated.
//   - email: The email address to verify.
//   - code: The verification code associated with the email.
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

// GenerateRefreshToken creates a new refresh token for the user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: The ID of the user for whom the refresh token is generated.
//
// Returns:
//   - token: The generated refresh token.
func (t *TokenRepository) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, int64, error) {
	resp, err := t.client.GenerateRefreshToken(ctx, &tokenv1.GenerateRefreshTokenRequest{UserId: userID.String()})
	if err != nil {
		joinedErr := errors.Join(err, "Failed to generate refresh token")
		return "", 0, errors.Upgrade(joinedErr, errcode.ErrInternalFailure, "TokenGenerationFailed")
	}
	return resp.Token, resp.ExpiresAt, nil
}

// VerifyAccessToken checks if the provided access token is valid.
//
// Parameters:
//   - ctx: The context for the operation.
//   - token: The access token to verify.
//
// Returns:
//   - valid: A boolean indicating whether the token is valid.
//   - userID: The ID of the user associated with the token if valid, otherwise nil.
//   - error: An error if the verification fails, otherwise nil.
func (t *TokenRepository) VerifyAccessToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := t.client.VerifyAccessToken(ctx, &tokenv1.VerifyAccessTokenRequest{Token: token})
	if err != nil {
		joinErr := errors.Join(err, "Failed to verify access token")
		return false, nil, errors.Upgrade(joinErr, errcode.ErrInternalFailure, "TokenVerificationFailed")
	}
	return resp.Valid, resp.UserId, nil
}

// VerifyEmailVerificationToken checks if the provided email verification token is valid.
//
// Parameters:
//   - ctx: The context for the operation.
//   - token: The email verification token to verify.
//
// Returns:
//   - data: A pointer to an EmailVerificationResult containing the verification result.
//   - error: An error if the verification fails, otherwise nil.
func (t *TokenRepository) VerifyEmailVerificationToken(ctx context.Context, token string) (*tokenmodels.EmailVerificationResult, error) {
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
	data := &tokenmodels.EmailVerificationResult{
		Valid:  resp.Valid,
		UserID: userUUID,
		Email:  *resp.Email,
		Code:   *resp.Code,
	}
	return data, nil
}

// VerifyRefreshToken checks if the provided refresh token is valid.
//
// Parameters:
//   - ctx: The context for the operation.
//   - token: The refresh token to verify.
//
// Returns:
//   - valid: A boolean indicating whether the token is valid.
//   - userID: The ID of the user associated with the token if valid, otherwise nil.
//   - error: An error if the verification fails, otherwise nil.
func (t *TokenRepository) VerifyRefreshToken(ctx context.Context, token string) (bool, *string, error) {
	resp, err := t.client.VerifyRefreshToken(ctx, &tokenv1.VerifyRefreshTokenRequest{Token: token})
	if err != nil {
		joinErr := errors.Join(err, "Failed to verify refresh token")
		return false, nil, errors.Upgrade(joinErr, errcode.ErrInternalFailure, "TokenVerificationFailed")
	}
	return resp.Valid, resp.UserId, nil
}

func NewTokenRepository(client tokenv1.TokenServiceClient) *TokenRepository {
	return &TokenRepository{client: client}
}
