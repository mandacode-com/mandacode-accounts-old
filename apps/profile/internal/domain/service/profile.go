package svcdomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/domain/dto"
)

type ProfileService interface {
	// GetProfileByID retrieves a profile by its ID.
	//
	// Parameter:
	//   - userID: The unique identifier of the user whose profile is to be retrieved.
	//
	// Returns:
	//   - *dto.Profile: The profile DTO if found.
	GetProfileByID(userID uuid.UUID) (*dto.Profile, error)

	// InitializeProfile creates a new profile for a user with the given ID.
	//
	// Parameters:
	//   - userID: The unique identifier of the user for whom the profile is to be created.
	//
	// Returns:
	//   - *dto.Profile: The newly created profile DTO.
	//	 - error: An error if the profile could not be created.
	InitializeProfile(userID uuid.UUID) (*dto.Profile, error)

	// UpdateProfile updates an existing profile with the given details.
	//
	// Parameters:
	//   - userID: The unique identifier of the user whose profile is to be updated.
	//   - email: The new email address of the profile owner (optional).
	//   - displayName: The new display name of the profile owner (optional).
	//   - bio: A new short biography of the profile owner (optional).
	//   - avatarURL: A new URL to the profile's avatar image (optional).
	UpdateProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string) (*dto.Profile, error)

	// DeleteProfile deletes a profile by its ID.
	//
	// Parameter:
	//   - userID: The unique identifier of the user whose profile is to be deleted.
	DeleteProfile(userID uuid.UUID) error
}
