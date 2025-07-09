package tokendomain

import (
	"context"

	"github.com/google/uuid"
)

type EmailVerificationResult struct {
	Valid  bool      `json:"valid"`
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Code   string    `json:"code"`
}

type TokenRepository interface {
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
	GenerateAccessToken(ctx context.Context, userID uuid.UUID) (string, int64, error)

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
	VerifyAccessToken(ctx context.Context, token string) (bool, *string, error)

	// GenerateRefreshToken creates a new refresh token for the user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user for whom the refresh token is generated.
	//
	// Returns:
	//   - token: The generated refresh token.
	GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, int64, error)

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
	VerifyRefreshToken(ctx context.Context, token string) (bool, *string, error)

	// GenerateEmailVerificationToken creates a new email verification token for the user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user for whom the email verification token is generated.
	//   - email: The email address to verify.
	//   - code: The verification code associated with the email.
	GenerateEmailVerificationToken(ctx context.Context, userID uuid.UUID, email, code string) (string, int64, error)

	// VerifyEmailVerificationToken checks if the provided email verification token is valid.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - token: The email verification token to verify.
	//
	// Returns:
	//   - EmailVerificationTokenData: The data associated with the email verification token, including validity, user ID, email, and code.
	//   - error: An error if the verification fails, otherwise nil.
	VerifyEmailVerificationToken(ctx context.Context, token string) (*EmailVerificationResult, error)
}
