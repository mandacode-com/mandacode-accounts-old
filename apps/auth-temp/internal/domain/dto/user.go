package dto

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
)

// LocalUser represents a user authenticated via local credentials.
type LocalUser struct {
	ID         uuid.UUID
	Email      string
	IsActive   bool
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewLocalUser(id uuid.UUID, email string, isActive bool, isVerified bool, createdAt time.Time, updatedAt time.Time) *LocalUser {
	return &LocalUser{
		ID:         id,
		Email:      email,
		IsActive:   isActive,
		IsVerified: isVerified,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func NewLocalUserFromEntity(user *ent.LocalUser) *LocalUser {
	return &LocalUser{
		ID:         user.ID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

type LocalDeletedUser struct {
	ID        uuid.UUID
	DeletedAt time.Time
}

func NewLocalDeletedUser(id uuid.UUID, deletedAt time.Time) *LocalDeletedUser {
	return &LocalDeletedUser{
		ID:        id,
		DeletedAt: deletedAt,
	}
}

// OAuthUser represents a user authenticated via an OAuth provider.
type OAuthUser struct {
	ID         uuid.UUID
	Provider   oauthuser.Provider
	ProviderID string
	Email      string
	IsActive   bool
	IsVerified bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewOAuthUser(id uuid.UUID, provider oauthuser.Provider, providerID string, email string, isActive bool, isVerified bool, createdAt time.Time, updatedAt time.Time) *OAuthUser {
	return &OAuthUser{
		ID:         id,
		Provider:   provider,
		ProviderID: providerID,
		Email:      email,
		IsActive:   isActive,
		IsVerified: isVerified,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}

func NewOAuthUserFromEntity(user *ent.OAuthUser) *OAuthUser {
	return &OAuthUser{
		ID:         user.ID,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

type OAuthDeletedUser struct {
	ID        uuid.UUID
	Provider  *oauthuser.Provider
	DeletedAt time.Time
}

func NewOAuthDeletedUser(id uuid.UUID, provider *oauthuser.Provider, deletedAt time.Time) *OAuthDeletedUser {
	return &OAuthDeletedUser{
		ID:        id,
		Provider:  provider,
		DeletedAt: deletedAt,
	}
}
