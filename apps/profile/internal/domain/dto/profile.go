package dto

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Profile struct {
	UserID      uuid.UUID `validate:"required,uuid"`
	Email       *string   `validate:"omitempty,email"`
	DisplayName *string   `validate:"omitempty,min=1,max=100"`
	Bio         *string   `validate:"omitempty,max=500"`
	AvatarURL   *string   `validate:"omitempty,url"`
	CreatedAt   time.Time `validate:"required"`
	UpdatedAt   time.Time `validate:"required"`
}

// NewProfile creates a new Profile instance with the provided parameters.
//
// Parameters:
//   - userID: The unique identifier for the user.
//   - email: The email address associated with the profile (optional).
//   - displayName: The display name of the profile owner (optional).
//   - bio: A short biography of the profile owner (optional).
//   - avatarURL: URL to the profile's avatar image (optional).
//   - createdAt: The time when the profile was created.
//   - updatedAt: The time when the profile was last updated.
//
// Returns:
//   - *Profile: A pointer to the newly created Profile instance.
func NewProfile(userID uuid.UUID, email, displayName, bio, avatarURL *string, createdAt, updatedAt time.Time) *Profile {
	return &Profile{
		UserID:      userID,
		Email:       email,
		DisplayName: displayName,
		Bio:         bio,
		AvatarURL:   avatarURL,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

var validate = validator.New()

// Validate checks if the Profile instance is valid according to the defined validation rules.
func (p *Profile) Validate() error {
	return validate.Struct(p)
}

// IsValid checks if the Profile instance is valid by calling the Validate method.
func (p *Profile) IsValid() bool {
	if err := p.Validate(); err != nil {
		return false
	}
	return true
}
