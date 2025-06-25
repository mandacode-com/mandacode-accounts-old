package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
)

// LocalUserRepository defines the interface for local authentication repository operations.
type LocalUserRepository interface {
	GetUserByEmail(email string) (*ent.LocalUser, error)
	CreateLocalUser(userID uuid.UUID, email string, password string, isActive *bool, isVerified *bool) (*ent.LocalUser, error)
	DeleteLocalUser(userID uuid.UUID) error
	UpdateLocalUser(userID uuid.UUID, email *string, password *string, isActive *bool, isVerified *bool) (*ent.LocalUser, error)
}
