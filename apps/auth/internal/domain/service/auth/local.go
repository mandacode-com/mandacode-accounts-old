package authdomain

import (
	"context"

	"mandacode.com/accounts/auth/internal/domain/dto"
)

type LocalAuthService interface {
	// LoginLocalUser logs in a user with email and password.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - email: The user's email address.
	//   - password: The user's password.
	//
	// Returns:
	//   - user: The local user if login is successful.
	//	 - error: An error if the login fails, otherwise nil.
	LoginLocalUser(ctx context.Context, email string, password string) (*dto.LocalUser, error)
}
