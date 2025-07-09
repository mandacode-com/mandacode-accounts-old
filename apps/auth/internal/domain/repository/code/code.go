package codedomain

import (
	"context"

	"github.com/google/uuid"
)

type CodeManager interface {
	// IssueCode issues a new login code for the given user ID.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The unique identifier of the user.
	//
	// Returns:
	//   - A string representing the issued login code.
	//   - An error if the code could not be issued.
	IssueCode(ctx context.Context, userID uuid.UUID) (string, error)

	// ValidateCode validates the provided login code for the given user ID.
	//
	// Parameters:
	//   - ctx: The context for the operation.
	//   - userID: The unique identifier of the user.
	//   - code: The login code to validate.
	//
	// Returns:
	//   - A boolean indicating whether the code is valid.
	//   - An error if the validation fails.
	ValidateCode(ctx context.Context, userID uuid.UUID, code string) (bool, error)
}
