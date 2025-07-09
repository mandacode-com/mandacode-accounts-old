package locallogin

import (
	"context"

	logindto "mandacode.com/accounts/mobile-auth/internal/app/login/dto"
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
	//   - LoginToken: A pointer to a LoginToken containing the access and refresh tokens if the login is successful.
	//   - Error: An error if the login operation fails, otherwise nil.
	Login(ctx context.Context, email, password string) (*logindto.LoginToken, error)
}
