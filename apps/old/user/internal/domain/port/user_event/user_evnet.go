package usereventdomain

import "github.com/google/uuid"

type UserEventService interface {
	// DeleteUser deletes a user from the external service.
	//
	// Parameters:
	//   - userID: The ID of the user to be deleted.
	DeleteUser(userID uuid.UUID) error

	// ArchiveUser archives a user in the external service.
	//
	// Parameters:
	//   - userID: The ID of the user to be archived.
	ArchiveUser(userID uuid.UUID) error

	// UserCreationFailed handles the event when user creation fails in the external service.
	//
	// Parameters:
	//   - userID: The ID of the user for whom creation failed.
	UserCreationFailed(userID uuid.UUID) error
}
