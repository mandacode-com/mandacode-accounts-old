package profiledomain

import (
	"context"

	"github.com/google/uuid"
	profilev1 "github.com/mandacode-com/accounts-proto/profile/v1"
)

type ProfileService interface {
	// CreateProfile creates a new user profile.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user for whom the profile is created.
	CreateProfile(ctx context.Context, userID uuid.UUID) (*profilev1.CreateProfileResponse, error)

	// DeleteProfile deletes an existing user profile.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose profile is to be deleted.
	DeleteProfile(ctx context.Context, userID uuid.UUID) (*profilev1.DeleteProfileResponse, error)
}
