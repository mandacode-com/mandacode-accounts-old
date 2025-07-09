package profile

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/domain/model"
)

type ProfileApp interface {
	// InitializeProfile initializes a new profile with the given details.
	//
	// Parameters:
	//   - userID: The unique identifier of the user for whom the profile is to be initialized.
	//
	// Returns:
	//   - *Profile: The newly initialized profile.
	//	 - error: An error if the profile could not be initialized.
	InitializeProfile(userID uuid.UUID) (*model.Profile, error)

	// GetProfile retrieves a profile by its user ID.
	//
	// Parameter:
	//   - userID: The unique identifier of the user whose profile is to be retrieved.
	//
	// Returns:
	//   - *Profile: The profile if found.
	//   - error: An error if the profile could not be found or another issue occurred.
	GetProfile(userID uuid.UUID) (*model.Profile, error)

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
	//   - *Profile: The updated profile.
	//   - error: An error if the profile could not be updated.
	UpdateProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string) (*model.Profile, error)

	// DeleteProfile deletes a profile by its ID.
	//
	// Parameter:
	//   - userID: The unique identifier of the user whose profile is to be deleted.
	//
	// Returns:
	//   - error: An error if the profile could not be deleted or does not exist.
	DeleteProfile(userID uuid.UUID) error
}
