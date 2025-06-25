package userdomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/internal/domain/dto"
)

type LocalUserService interface {
	// CreateLocalUser creates a new local user with the given details.
	//
	// Parameters:
	//   - userID: The unique identifier for the user.
	//   - email: The user's email address.
	//   - password: The user's password.
	//   - isActive: A pointer to a boolean indicating if the user is active.
	//   - isVerified: A pointer to a boolean indicating if the user's email is verified.
	//
	// Returns:
	//   - user: The created local user if the creation is successful.
	//	 - error: An error if the creation fails, otherwise nil.
	CreateLocalUser(userID uuid.UUID, email, password string, isActive, isVerified *bool) (*dto.LocalUser, error)

	// DeleteLocalUser deletes a local user by their userID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be deleted.
	//
	// Returns:
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteLocalUser(userID uuid.UUID) error

	// UpdateLocalUser updates the details of a local user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be updated.
	//   - email: A pointer to the new email address (nil if not updating).
	//   - password: A pointer to the new password (nil if not updating).
	//   - isActive: A pointer to a boolean indicating if the user is active (nil if not updating).
	//   - isVerified: A pointer to a boolean indicating if the user's email is verified (nil if not updating).
	//
	// Returns:
	//   - user: The updated local user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateLocalUser(userID uuid.UUID, email, password *string, isActive, isVerified *bool) (*dto.LocalUser, error)
}
