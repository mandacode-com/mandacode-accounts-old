package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/user/ent"
)

type UserRepository interface {
	// CreateUser creates a new user and returns the user ID
	CreateUser(userID uuid.UUID, syncCode string) (*ent.User, error)

	// GetUser retrieves a user by ID
	GetUser(userID uuid.UUID) (*ent.User, error)

	// DeleteUser deletes a user by ID
	DeleteUser(userID uuid.UUID) error

	// ArchiveUser archives a user by ID
	ArchiveUser(userID uuid.UUID, syncCode string) error

	// UpdateEmailVerificationCode updates the email verification code for a user
	UpdateEmailVerificationCode(userID uuid.UUID, emailVerificationCode string) (*ent.User, error)
}
