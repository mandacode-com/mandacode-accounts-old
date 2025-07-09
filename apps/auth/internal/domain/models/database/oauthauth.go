package dbmodels

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent/oauthauth"
)

type CreateOAuthAuthInput struct {
	AccountID  uuid.UUID          `json:"account_id" validate:"required"`
	Provider   oauthauth.Provider `json:"provider" validate:"required,oneof=google github facebook"`
	ProviderID string             `json:"provider_id" validate:"required"`
	Email      string             `json:"email" validate:"omitempty,email"`
	IsVerified bool               `json:"is_verified" validate:"required"`
}
