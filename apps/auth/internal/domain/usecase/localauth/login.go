package localauthdomain

import (
	"context"

	"github.com/google/uuid"
)

type LoginUsecase interface {
	// Login authenticates a user with the provided email and password.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - input: The login input containing email, password, and other optional parameters.
	Login(ctx context.Context, input LoginInput) (accessToken, refreshToken string, err error)

	// IssueLoginCode issues a login code for the user identified by the email.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - input: The input containing the email of the user to issue a login code for.
	IssueLoginCode(ctx context.Context, input LoginInput) (code string, userID uuid.UUID, err error)

	// VerifyLoginCode verifies the provided login code for the user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to verify.
	//   - code: The login code to be verified.
	VerifyLoginCode(ctx context.Context, userID uuid.UUID, code string) (accessToken, refreshToken string, err error)
}
