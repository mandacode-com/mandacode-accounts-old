package localuserapp

import (
	"context"

	"github.com/google/uuid"
)

type LocalUserApp interface {
	// CreateUser creates a new user with the provided email and password.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - email: The email address of the user.
	//   - password: The password for the user.
	CreateUser(ctx context.Context, email string, password string) (string, error)

	// VerifyUserEmail verifies a user's email address.
	//
	// Parameters:
	//	 - ctx: The context for the operation.
	//   - userID: The ID of the user whose email is to be verified.
	//   - verificationToken: The token used to verify the user's email. 
	VerifyUserEmail(ctx context.Context, userID uuid.UUID, verificationToken string) error

	// UpdateEmail updates the email address of a user.
	//
	// Parameters:
	//	 - ctx: The context for the operation.
	//   - userID: The ID of the user whose email is to be updated.
	//	 - email: The new email address for the user.
	// UpdateEmail(ctx context.Context, userID uuid.UUID, email string) error

	// UpdatePassword updates the password for a user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose password is to be updated.
	//   - currentPassword: The current password of the user.
	//	 - newPassword: The new password for the user.
	UpdatePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error
}
