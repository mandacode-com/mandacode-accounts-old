package repodomain

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
)

// OAuthUserRepository defines the interface for OAuth authentication repository operations.
type OAuthUserRepository interface {
	GetUserByProvider(
		provider oauthuser.Provider,
		providerID string,
	) (*ent.OAuthUser, error)
	CreateUser(userID uuid.UUID, provider oauthuser.Provider, providerID string, email string, isActive *bool, isVerified *bool) (*ent.OAuthUser, error)
	DeleteUser(userID uuid.UUID) error
	DeleteUserByProvider(
		userID uuid.UUID,
		provider oauthuser.Provider,
	) error
	UpdateUser(userID uuid.UUID, provider oauthuser.Provider, providerID *string, email *string, isActive *bool, isVerified *bool) (*ent.OAuthUser, error)
}
