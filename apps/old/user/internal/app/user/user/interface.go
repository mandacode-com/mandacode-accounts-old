package userapp

import (
	"context"

	"github.com/google/uuid"
)

type UserApp interface {
	// DeleteUser deletes a user by ID
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to be deleted.
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// ArchiveUser archives a user by ID
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user to be archived.
	ArchiveUser(ctx context.Context, userID uuid.UUID) error

	// UpdateActiveStatus updates the active status of a user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose active status is to be updated.
	//	 - isActive: Indicates if the user is active.
	UpdateActiveStatus(ctx context.Context, userID uuid.UUID, isActive bool) error

	// UpdateVerifiedStatus updates the verified status of a user.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The ID of the user whose verified status is to be updated.
	//   - isVerified: Indicates if the user is verified.
	UpdateVerifiedStatus(ctx context.Context, userID uuid.UUID, isVerified bool) error
}
