package tokendomain

import (
	"context"

	"github.com/google/uuid"
)

type TokenService interface {
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
	//   - valid: A boolean indicating whether the token is valid.
	//   - userID: The ID of the user associated with the token if valid, otherwise nil.
	//   - email: The email address associated with the token if valid, otherwise nil.
	//   - code: The verification code associated with the email if valid, otherwise nil.
	//   - error: An error if the verification fails, otherwise nil.
	VerifyEmailVerificationToken(ctx context.Context, token string) (bool, *string, *string, *string, error)
}
