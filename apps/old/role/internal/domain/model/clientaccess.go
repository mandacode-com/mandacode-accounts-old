package model

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type ClientAccess struct {
	ID           uuid.UUID `json:"id"`
	ServiceID    uuid.UUID `json:"service_id"`
	ClientID     string    `json:"client_id"`
	ClientSecret string    `json:"client_secret"`
	Description  string    `json:"description,omitempty"` // Optional field, can be empty
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ServiceClientFromEnt converts an ent ServiceClient entity to a model ServiceClient.
func ServiceClientFromEnt(ca *ent.ClientAccess) *ClientAccess {
	if ca == nil {
		return nil
	}
	return &ClientAccess{
		ID:           ca.ID,
		ServiceID:    ca.ServiceID,
		ClientID:     ca.ClientID,
		ClientSecret: ca.ClientSecret,
		Description:  ca.Description,
		IsActive:     ca.IsActive,
		CreatedAt:    ca.CreatedAt,
		UpdatedAt:    ca.UpdatedAt,
	}
}
