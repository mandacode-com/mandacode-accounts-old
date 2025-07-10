package localuser

import (
	"github.com/google/uuid"
	userdto "mandacode.com/accounts/auth/internal/app/user/dto"
)

type LocalUserApp interface {
	// GetUserByEmail retrieves a user by their email address.
	//
	// Parameters:
	// - email: The user's email address.
	//
	// Returns:
	// - *LocalUser: The user information if found.
	// - error: An error if the user is not found or if there is an issue retrieving the user.
	GetUserByEmail(email string) (*userdto.LocalUser, error)

	// CreateUser creates a new user with the provided details.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - email: The user's email address.
	// - password: The user's password.
	// - isActive: Optional flag indicating if the user is active.
	// - isVerified: Optional flag indicating if the user is verified.
	//
	// Returns:
	// - *LocalUser: The created user information.
	// - error: An error if the creation fails, otherwise nil.
	CreateUser(userID uuid.UUID, email, password string, isActive, isVerified *bool) (*userdto.LocalUser, error)

	// GetUser retrieves a user by their unique identifier.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	//
	// Returns:
	// - *LocalUser: The user information if found.
	// - error: An error if the user is not found or if there is an issue retrieving the user.
	GetUser(userID uuid.UUID) (*userdto.LocalUser, error)

	// UpdateEmail updates the email address of a user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - email: The new email address for the user.
	//
	// Returns:
	// - *LocalUser: The updated user information.
	// - error: An error if the update fails, otherwise nil.
	UpdateEmail(userID uuid.UUID, email string) (*userdto.LocalUser, error)

	// UpdatePassword updates the password of a user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - currentPassword: The user's current password for verification.
	// - newPassword: The new password to set for the user.
	//
	// Returns:
	// - *LocalUser: The updated user information.
	// - error: An error if the update fails, otherwise nil.
	UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) (*userdto.LocalUser, error)

	// UpdateActiveStatus updates the active status of a user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - isActive: The new active status to set for the user.
	//
	// Returns:
	// - *LocalUser: The updated user information.
	// - error: An error if the update fails, otherwise nil.
	UpdateActiveStatus(userID uuid.UUID, isActive bool) (*userdto.LocalUser, error)

	// UpdateVerificationStatus updates the verification status of a user.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	// - isVerified: The new verification status to set for the user.
	//
	// Returns:
	// - *LocalUser: The updated user information.
	// - error: An error if the update fails, otherwise nil.
	UpdateVerificationStatus(userID uuid.UUID, isVerified bool) (*userdto.LocalUser, error)

	// DeleteUser deletes a user by their unique identifier.
	//
	// Parameters:
	// - userID: The unique identifier for the user.
	//
	// Returns:
	// - error: An error if the deletion fails, otherwise nil.
	DeleteUser(userID uuid.UUID) error
}
