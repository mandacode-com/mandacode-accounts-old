package dbmodels

import (
	"time"

	"github.com/google/uuid"
)

type CreateAuthAccountInput struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
}

type UpdateAuthAccountInput struct {
	UserID              *uuid.UUID `json:"user_id" validate:"omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	LastFailedLoginAt   *time.Time `json:"last_failed_login_at,omitempty"`
	FailedLoginAttempts *int       `json:"failed_login_attempts,omitempty"`
}
