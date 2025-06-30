package userdto

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/auth/ent"
)

type LocalUser struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	IsActive   bool      `json:"is_active"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func NewLocalUserFromEnt(entUser *ent.LocalUser) *LocalUser {
	if entUser == nil {
		return nil
	}
	return &LocalUser{
		ID:         entUser.ID,
		Email:      entUser.Email,
		IsActive:   entUser.IsActive,
		IsVerified: entUser.IsVerified,
		CreatedAt:  entUser.CreatedAt,
		UpdatedAt:  entUser.UpdatedAt,
	}
}
