package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
)

// LocalUserRepository defines the interface for local authentication repository operations.
type LocalUserRepository interface {
	GetUserByEmail(email string) (*ent.LocalUser, error)
	GetUserByID(userID uuid.UUID) (*ent.LocalUser, error)
	CreateUser(userID uuid.UUID, email string, password string, isActive *bool, isVerified *bool) (*ent.LocalUser, error)
	DeleteUser(userID uuid.UUID) error
	UpdateUser(userID uuid.UUID, email *string, password *string, isActive *bool, isVerified *bool) (*ent.LocalUser, error)
}
