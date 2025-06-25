package dto

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthuser"
)

// LocalUser represents a user authenticated via local credentials.
type LocalUser struct {
	ID         uuid.UUID
	Email      string
	IsActive   bool
	IsVerified bool
}

// OAuthUser represents a user authenticated via an OAuth provider.
type OAuthUser struct {
	ID         uuid.UUID
	Provider   oauthuser.Provider
	ProviderID string
	Email      string
	IsActive   bool
	IsVerified bool
}
