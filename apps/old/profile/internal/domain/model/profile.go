package model

import (
	"time"

	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/profile/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/profile/ent"
)

type Profile struct {
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid"`
	Email       *string   `json:"email,omitempty" validate:"omitempty,email"`
	DisplayName *string   `json:"display_name,omitempty" validate:"omitempty,min=1,max=100"`
	Bio         *string   `json:"bio,omitempty" validate:"omitempty,max=500"`
	AvatarURL   *string   `json:"avatar_url,omitempty" validate:"omitempty,url"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
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

// NewProfileFromEnt converts an ent.Profile instance to a Profile instance.
//
// Parameters:
//   - profile: The ent.Profile instance to convert.
//
// Returns:
//   - *Profile: A pointer to the newly created Profile instance.
func NewProfileFromEnt(profile *ent.Profile) *Profile {
	return &Profile{
		UserID:      profile.ID,
		Email:       &profile.Email,
		DisplayName: &profile.DisplayName,
		Bio:         &profile.Bio,
		AvatarURL:   &profile.AvatarURL,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}

// NewProfileFromProto converts a profilev1.Profile instance to a Profile instance.
//
// Parameters:
//   - profile: The profilev1.Profile instance to convert.
//
// Returns:
//   - *Profile: A pointer to the newly created Profile instance.
//   - error: An error if the conversion fails (e.g., invalid UUID).
func NewProfileFromProto(profile *profilev1.Profile) (*Profile, error) {
	userID, err := uuid.Parse(profile.UserId)
	if err != nil {
		return nil, err
	}
	return &Profile{
		UserID:      userID,
		Email:       profile.Email,
		DisplayName: profile.DisplayName,
		Bio:         profile.Bio,
		AvatarURL:   profile.AvatarUrl,
		CreatedAt:   profile.CreatedAt.AsTime(),
		UpdatedAt:   profile.UpdatedAt.AsTime(),
	}, nil
}

// ToProtoProfile converts a Profile instance to a profilev1.Profile instance.
//
// Returns:
//   - *profilev1.Profile: A pointer to the newly created profilev1.Profile instance.
func (p *Profile) ToProtoProfile() *profilev1.Profile {
	return &profilev1.Profile{
		UserId:      p.UserID.String(),
		Email:       p.Email,
		DisplayName: p.DisplayName,
		Bio:         p.Bio,
		AvatarUrl:   p.AvatarURL,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
