package adminhandlerdto

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
)

type CreateGroupRequest struct {
	Name        string    `json:"name" validate:"required"`
	Description *string   `json:"description,omitempty"`
	ServiceID   uuid.UUID `json:"service_id" validate:"required"`
	IsActive    *bool     `json:"is_active"` // Indicates if the group is active or not
}
type CreateGroupResponse struct {
	*model.Group // Embedding model.Group to include group details in the response
}

type GetGroupByIDResponse struct {
	*model.Group // Embedding model.Group to include group details in the response
}

type GetAllGroupsResponse struct {
	Groups []*model.Group `json:"groups"` // Slice of model.Group to hold multiple group details
}

type UpdateGroupRequest struct {
	Name        *string    `json:"name,omitempty"`        // Optional field, can be nil
	Description *string    `json:"description,omitempty"` // Optional field, can be nil
	ServiceID   *uuid.UUID `json:"service_id,omitempty"`  // Optional field, can be nil
	IsActive    *bool      `json:"is_active,omitempty"`   // Optional field, can be nil
}
type UpdateGroupResponse struct {
	*model.Group // Embedding model.Group to include group details in the response
}

type DeleteGroupResponse struct {
	Success bool `json:"success"` // Indicates whether the deletion was successful
}
