package locallogin

import (
	"context"

	logindto "mandacode.com/accounts/web-auth/internal/app/login/dto"
)

type LocalLoginApp interface {
	// Login performs the local login operation.
	//
	// Parameters:
	//   - ctx: The context for the login operation.
	//   - email: The email of the user attempting to log in.
	//   - password: The password of the user attempting to log in.
	//
	// Returns:
	//   - Code: A string representing the result of the login operation.
	//   - Error: An error if the login operation fails, otherwise nil.
	Login(ctx context.Context, email, password string) (code string, err error)

	// VerifyCode verifies the provided verification code for the given email.
	//
	// Parameters:
	//   - ctx: The context for the verification operation.
	//   - code: The verification code to be verified.
	//
	// Returns:
	//   - LoginToken: A pointer to a LoginToken containing the access and refresh tokens if the verification is successful.
	//   - Error: An error if the verification fails, otherwise nil.
	VerifyCode(ctx context.Context, code string) (*logindto.LoginToken, error)
}
