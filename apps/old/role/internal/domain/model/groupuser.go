package model

import (
	"time"

	"github.com/google/uuid"
	"mandacode.com/accounts/role/ent"
)

type GroupUser struct {
	UserID    uuid.UUID `json:"user_id"`
	GroupID   uuid.UUID `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GroupUserFromEnt converts an ent GroupUser entity to a model GroupUser.
func GroupUserFromEnt(gu *ent.GroupUser) *GroupUser {
	if gu == nil {
		return nil
	}
	return &GroupUser{
		UserID:    gu.UserID,
		GroupID:   gu.GroupID,
		CreatedAt: gu.CreatedAt,
		UpdatedAt: gu.UpdatedAt,
	}
}

// ToEnt converts a model GroupUser to an ent GroupUser entity.
func (gu *GroupUser) ToEnt() *ent.GroupUser {
	if gu == nil {
		return nil
	}
	return &ent.GroupUser{
		UserID:    gu.UserID,
		GroupID:   gu.GroupID,
		CreatedAt: gu.CreatedAt,
		UpdatedAt: gu.UpdatedAt,
	}
}
