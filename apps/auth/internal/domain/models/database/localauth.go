package dbmodels

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
)

type CreateLocalAuthInput struct {
	AccountID  uuid.UUID `json:"account_id" validate:"required"`
	Email      string    `json:"email" validate:"required,email"`
	Password   string    `json:"password" validate:"required,min=8,max=100"`
	IsActive   bool      `json:"is_active" validate:"required"`
	IsVerified bool      `json:"is_verified" validate:"required"`
}

// SecureLocalAuth represents the local authentication details for an account.
// This does not include sensitive information like passwords.
type SecureLocalAuth struct {
	AccountID           uuid.UUID `json:"account_id" validate:"required"`
	Email               string    `json:"email" validate:"required,email"`
	IsActive            bool      `json:"is_active" validate:"required"`
	IsVerified          bool      `json:"is_verified" validate:"required"`
	LastLoginAt         time.Time `json:"last_login_at"`
	LastFailedLoginAt   time.Time `json:"last_failed_login_at"`
	FailedLoginAttempts int      `json:"failed_login_attempts,omitempty"`
}

// NewSecureLocalAuth creates a new instance of SecureLocalAuth with the provided account ID and email.
func NewSecureLocalAuth(localAuth *ent.LocalAuth) *SecureLocalAuth {
	return &SecureLocalAuth{
		AccountID:           localAuth.AuthAccountID,
		Email:               localAuth.Email,
		IsActive:            localAuth.IsActive,
		IsVerified:          localAuth.IsVerified,
		LastLoginAt:         localAuth.LastLoginAt,
		LastFailedLoginAt:   localAuth.LastFailedLoginAt,
		FailedLoginAttempts: localAuth.FailedLoginAttempts,
	}
}
