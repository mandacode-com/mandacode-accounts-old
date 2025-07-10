package adminhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
	adminhandlerdto "mandacode.com/accounts/role/internal/handler/http/admin/dto"
)

type GroupUserHandler struct {
	groupUserApp    groupuserapp.GroupUserApp
	validator       *validator.Validate
	adminMiddleware gin.HandlerFunc
}

// NewGroupUserHandler creates a new GroupUserHandler instance.
func NewGroupUserHandler(
	groupUserApp groupuserapp.GroupUserApp,
	validator *validator.Validate,
	adminMiddleware gin.HandlerFunc,
) (*GroupUserHandler, error) {
	return &GroupUserHandler{
		groupUserApp:    groupUserApp,
		validator:       validator,
		adminMiddleware: adminMiddleware,
	}, nil
}

// RegisterRoutes registers the routes for the GroupUserHandler.
func (h *GroupUserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(h.adminMiddleware)
	rg.GET("/:group_id/users/:user_id", h.GetGroupUser)
	rg.GET("/:group_id/users", h.GetAllGroupUsersByGroupID)
	rg.POST("/:group_id/users", h.CreateGroupUser)
	rg.DELETE("/:group_id/users/:user_id", h.DeleteGroupUser)
	rg.DELETE("/:group_id/users", h.DeleteGroupUserByGroupID)
}

// GetGroupUser handles the retrieval of a specific group user by group ID and user ID.
func (h *GroupUserHandler) GetGroupUser(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	groupUser, err := h.groupUserApp.GetGroupUser(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve group user"})
		return
	}
	if groupUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group user not found"})
		return
	}

	c.JSON(http.StatusOK, adminhandlerdto.GetGroupUserResponse{
		GroupUser: groupUser,
	})
}

// GetAllGroupUsersByGroupID handles the retrieval of all users in a specific group by group ID.
func (h *GroupUserHandler) GetAllGroupUsersByGroupID(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	groupUsers, err := h.groupUserApp.GetGroupUsersByGroupID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve group users"})
		return
	}

	c.JSON(http.StatusOK, adminhandlerdto.GetAllGroupUsersResponse{
		GroupUsers: groupUsers,
	})
}

// CreateGroupUser handles the creation of a new group user.
func (h *GroupUserHandler) CreateGroupUser(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var request adminhandlerdto.CreateGroupUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	groupUser, err := h.groupUserApp.CreateGroupUser(request.UserID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group user"})
		return
	}

	c.JSON(http.StatusCreated, adminhandlerdto.CreateGroupUserResponse{
		GroupUser: groupUser,
	})
}

// DeleteGroupUser handles the deletion of a specific group user by group ID and user ID.
func (h *GroupUserHandler) DeleteGroupUser(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	err = h.groupUserApp.DeleteGroupUser(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group user"})
		return
	}

	c.JSON(http.StatusOK, adminhandlerdto.DeleteGroupUserResponse{
		Success: true,
	})
}

// DeleteGroupUserByGroupID handles the deletion of all users in a specific group by group ID.
func (h *GroupUserHandler) DeleteGroupUserByGroupID(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	err = h.groupUserApp.DeleteGroupUserByGroupID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group users"})
		return
	}

	c.JSON(http.StatusOK, adminhandlerdto.DeleteGroupUserByGroupIDResponse{
		Success: true,
	})
}
