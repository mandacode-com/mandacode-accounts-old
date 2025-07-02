package model

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ServiceID   uuid.UUID `json:"service_id"`
	IsActive    bool      `json:"is_active"`
	Description string    `json:"description,omitempty"` // Optional field, can be empty
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GroupFromEnt converts an ent Group entity to a model Group.
func GroupFromEnt(g *ent.Group) *Group {
	if g == nil {
		return nil
	}
	return &Group{
		ID:          g.ID,
		Name:        g.Name,
		ServiceID:   g.Edges.Service.ID,
		IsActive:    g.IsActive,
		Description: g.Description,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}

// ToEnt converts a model Group to an ent Group entity.
func (g *Group) ToEnt() *ent.Group {
	if g == nil {
		return nil
	}
	return &ent.Group{
		ID:          g.ID,
		Name:        g.Name,
		ServiceID:   g.ServiceID,
		IsActive:    g.IsActive,
		Description: g.Description,
		CreatedAt:   g.CreatedAt,
		UpdatedAt:   g.UpdatedAt,
	}
}
