package localauthdomain

import (
	"context"

	"github.com/google/uuid"
)

type SignupUsecase interface {
	// Signup registers a new user with the provided email and password.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - input: The signup input containing email, password, and other optional parameters.
	Signup(ctx context.Context, input SignupInput) (userID uuid.UUID, err error)

	// VerifyEmail verifies the user's email address using the provided verification code.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - email: The email address of the user to verify.
	//   - token: The verification token sent to the user's email.
	VerifyEmail(ctx context.Context, email, token string) (success bool, err error)

	// ResendVerificationEmail resends the email verification link to the user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - email: The email address of the user to resend the verification link to.
	ResendVerificationEmail(ctx context.Context, email string) (success bool, err error)
}
