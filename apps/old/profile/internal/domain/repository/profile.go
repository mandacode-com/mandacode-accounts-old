package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/ent"
)

type ProfileRepository interface {
	// GetProfileByID retrieves a profile by its ID.
	//
	// Parameter:
	//   - userID: The unique identifier of the user whose profile is to be retrieved.
	//
	// Returns:
	//   - *ent.Profile: The profile entity if found.
	//   - error: An error if the profile could not be found or another issue occurred.
	GetProfileByID(userID uuid.UUID) (*ent.Profile, error)

	// InitializeProfile creates a new profile for a user with the given ID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user for whom the profile is to be created.
	//
	// Returns:
	//   - *ent.Profile: The newly created profile entity.
	//   - error: An error if the profile could not be created.
	InitializeProfile(userID uuid.UUID) (*ent.Profile, error)

	// UpdateProfile updates an existing profile with the given details.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose profile is to be updated.
	//   - email: The new email address of the profile owner (optional).
	//   - displayName: The new display name of the profile owner (optional).
	//   - bio: A new short biography of the profile owner (optional).
	//   - avatarURL: A new URL to the profile's avatar image (optional).
	//
	// Returns:
	//   - *ent.Profile: The updated profile entity.
	//   - error: An error if the profile could not be updated.
	UpdateProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string) (*ent.Profile, error)

	// DeleteProfile deletes a profile by its ID.
	//
	// Parameter:
	//   - userID: The unique identifier of the user whose profile is to be deleted.
	//
	// Returns:
	//   - error: An error if the profile could not be deleted or does not exist.
	DeleteProfile(userID uuid.UUID) error
}
