package userdomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/internal/domain/dto"
)

type LocalUserService interface {
	// CreateUser creates a new local user with the given details.
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
	CreateUser(userID uuid.UUID, email, password string, isActive, isVerified *bool) (*dto.LocalUser, error)

	// GetUserByEmail retrieves a local user by their email address.
	//
	// Parameters:
	//   - email: The email address of the user to be retrieved.
	//
	// Returns:
	//   - user: The local user if found, otherwise nil.
	//   - error: An error if the retrieval fails, otherwise nil.
	GetUserByEmail(email string) (*dto.LocalUser, error)

	// GetUserByID retrieves a local user by their userID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be retrieved.
	//
	// Returns:
	//   - user: The local user if found, otherwise nil.
	//   - error: An error if the retrieval fails, otherwise nil.
	GetUserByID(userID uuid.UUID) (*dto.LocalUser, error)

	// DeleteUser deletes a local user by their userID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user to be deleted.
	//
	// Returns:
	//	 - user: The deleted local user if the deletion is successful.
	//   - error: An error if the deletion fails, otherwise nil.
	DeleteUser(userID uuid.UUID) (*dto.LocalDeletedUser, error)

	// UpdateEmail updates the email address of a local user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose email is to be updated.
	//   - newEmail: The new email address to set for the user.
	//
	// Returns:
	//   - user: The updated local user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateEmail(userID uuid.UUID, newEmail string) (*dto.LocalUser, error)

	// UpdatePassword updates the password of a local user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose password is to be updated.
	//   - currentPassword: The current password of the user (for verification).
	//   - newPassword: The new password to set for the user.
	//
	// Returns:
	//   - user: The updated local user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdatePassword(userID uuid.UUID, currentPassword, newPassword string) (*dto.LocalUser, error)

	// UpdateActiveStatus updates the active status of a local user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose active status is to be updated.
	//   - isActive: A boolean indicating if the user is active.

	// Returns:
	//   - user: The updated local user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateActiveStatus(userID uuid.UUID, isActive bool) (*dto.LocalUser, error)

	// UpdateVerifiedStatus updates the verified status of a local user.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose verified status is to be updated.
	//   - isVerified: A boolean indicating if the user's email is verified.
	//
	// Returns:
	//   - user: The updated local user if the update is successful.
	//   - error: An error if the update fails, otherwise nil.
	UpdateVerifiedStatus(userID uuid.UUID, isVerified bool) (*dto.LocalUser, error)
}
