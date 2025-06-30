package userdto

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
	"mandacode.com/accounts/auth/ent/oauthuser"
)

type OAuthUser struct {
	ID       uuid.UUID          `json:"id"`
	Provider oauthuser.Provider `json:"provider"`
	ProviderID string    `json:"provider_id"`
	Email      string    `json:"email"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewOAuthUserFromEnt(entUser *ent.OAuthUser) *OAuthUser {
	if entUser == nil {
		return nil
	}
	return &OAuthUser{
		ID:         entUser.ID,
		Provider:   entUser.Provider,
		ProviderID: entUser.ProviderID,
		Email:      entUser.Email,
		IsActive:   entUser.IsActive,
		IsVerified: entUser.IsVerified,
		CreatedAt:  entUser.CreatedAt,
		UpdatedAt:  entUser.UpdatedAt,
	}
}
