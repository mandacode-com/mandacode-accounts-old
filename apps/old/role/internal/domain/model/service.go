package model

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type Service struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"` // Optional field, can be empty
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServiceFromEnt converts an ent Service entity to a model Service.
func ServiceFromEnt(s *ent.Service) *Service {
	if s == nil {
		return nil
	}
	return &Service{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// ToEnt converts a model Service to an ent Service entity.
func (s *Service) ToEnt() *ent.Service {
	if s == nil {
		return nil
	}
	return &ent.Service{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
