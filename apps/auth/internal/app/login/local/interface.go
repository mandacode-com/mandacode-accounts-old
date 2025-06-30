package locallogin

import (
	"context"

	logindto "mandacode.com/accounts/auth/internal/app/login/dto"
)

type LocalLoginApp interface {
	// Login authenticates a user with their email and password.
	//
	// Parameters:
	// - ctx: The context for the operation.
	// - email: The user's email address.
	// - password: The user's password.
	//
	// Returns:
	// - *logindto.LoginToken: The login token containing access and refresh tokens.
	// - error: An error if the login fails, otherwise nil.
	Login(ctx context.Context, email, password string) (*logindto.LoginToken, error)
}
