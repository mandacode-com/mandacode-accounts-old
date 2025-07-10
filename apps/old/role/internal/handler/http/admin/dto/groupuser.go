package adminhandlerdto

import (
	"github.com/google/uuid"
	"mandacode.com/accounts/role/internal/domain/model"
)

// rg.GET("/:group_id/users/:user_id", h.GetGroupUser)
// rg.GET("/:group_id/users", h.GetAllGroupUsersByGroupID)
// rg.POST("/:group_id/users", h.CreateGroupUser)
// rg.DELETE("/:group_id/users/:user_id", h.DeleteGroupUser)
// rg.DELETE("/:group_id/users", h.DeleteGroupUserByGroupID)

type GetGroupUserResponse struct {
	GroupUser *model.GroupUser `json:"group_user"`
}

type GetAllGroupUsersResponse struct {
	GroupUsers []*model.GroupUser `json:"group_users"`
}

type CreateGroupUserRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required"`
}
type CreateGroupUserResponse struct {
	GroupUser *model.GroupUser `json:"group_user"`
}

type DeleteGroupUserResponse struct {
	Success bool `json:"success"`
}

type DeleteGroupUserByGroupIDResponse struct {
	Success bool `json:"success"`
}
