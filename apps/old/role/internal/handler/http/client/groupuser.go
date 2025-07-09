package clienthandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	groupapp "mandacode.com/accounts/role/internal/app/group"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
	permissionapp "mandacode.com/accounts/role/internal/app/permission"
	clienthandlerdto "mandacode.com/accounts/role/internal/handler/http/client/dto"
)

type ClientGroupUserHandler struct {
	permissionApp permissionapp.PermissionApp
	groupApp      groupapp.GroupApp
	groupUserApp  groupuserapp.GroupUserApp
	validator     *validator.Validate
}

func NewClientGroupUserHandler(
	permissionApp permissionapp.PermissionApp,
	groupApp groupapp.GroupApp,
	groupUserApp groupuserapp.GroupUserApp,
	validator *validator.Validate,
) (*ClientGroupUserHandler, error) {
	return &ClientGroupUserHandler{
		permissionApp: permissionApp,
		groupApp:      groupApp,
		groupUserApp:  groupUserApp,
		validator:     validator,
	}, nil
}

func (h *ClientGroupUserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/:group_id/users/:user_id/enroll", h.EnrollGroupUser)
	rg.POST("/:group_id/users/:user_id/exists", h.CheckGroupUser)
	rg.POST("/:group_id/users", h.GetAllGroupUsers)
	rg.DELETE("/:group_id/users/:user_id", h.DeleteGroupUser)
	rg.DELETE("/:group_id/users", h.DeleteGroupUserByGroupID)
}

func (h *ClientGroupUserHandler) checkClientAccess(groupID uuid.UUID, clientID, clientSecret string) (bool, error) {
	group, err := h.groupApp.GetGroupByID(groupID)
	if err != nil {
		return false, err
	}
	if group == nil {
		return false, nil // Group not found
	}
	valid, err := h.permissionApp.CheckClientAccess(group.ServiceID, clientID, clientSecret)
	if err != nil {
		return false, err // Error checking client access
	}
	return valid, nil // Return whether the client has access
}

func (h *ClientGroupUserHandler) EnrollGroupUser(c *gin.Context) {
	// Parse group ID and user ID from URL parameters
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

	// Get body data
	var request clienthandlerdto.EnrollGroupUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Validate request
	if err := h.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	// // Check client permissions
	valid, err := h.checkClientAccess(groupID, request.ClientID, request.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check client access"})
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	// Enroll user in group
	groupUser, err := h.groupUserApp.CreateGroupUser(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enroll user in group"})
		return
	}

	response := clienthandlerdto.EnrollGroupUserResponse{
		GroupUser: groupUser,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientGroupUserHandler) CheckGroupUser(c *gin.Context) {
	// Parse group ID and user ID from URL parameters
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

	// Get body data
	var request clienthandlerdto.CheckGroupUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Validate request
	if err := h.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	// // Check client permissions
	valid, err := h.checkClientAccess(groupID, request.ClientID, request.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check client access"})
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	exists, err := h.groupUserApp.CheckGroupUserExists(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check group user existence"})
		return
	}

	response := clienthandlerdto.CheckGroupUserResponse{
		IsEnrolled: exists,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientGroupUserHandler) GetAllGroupUsers(c *gin.Context) {
	// Parse group ID from URL parameters
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// Get body data
	var request clienthandlerdto.GetAllGroupUsersRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Validate request
	if err := h.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	// // Check client permissions
	valid, err := h.checkClientAccess(groupID, request.ClientID, request.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check client access"})
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	groupUsers, err := h.groupUserApp.GetGroupUsersByGroupID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group users"})
		return
	}

	response := clienthandlerdto.GetAllGroupUsersResponse{
		GroupUsers: groupUsers,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientGroupUserHandler) DeleteGroupUser(c *gin.Context) {
	// Parse group ID and user ID from URL parameters
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

	// Get body data
	var request clienthandlerdto.DeleteGroupUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Validate request
	if err := h.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	// // Check client permissions
	valid, err := h.checkClientAccess(groupID, request.ClientID, request.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check client access"})
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	err = h.groupUserApp.DeleteGroupUser(userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group user"})
		return
	}

	response := clienthandlerdto.DeleteGroupUserResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ClientGroupUserHandler) DeleteGroupUserByGroupID(c *gin.Context) {
	// Parse group ID from URL parameters
	groupID, err := uuid.Parse(c.Param("group_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	// Get body data
	var request clienthandlerdto.DeleteGroupUserByGroupIDRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Validate request
	if err := h.validator.Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	// // Check client permissions
	valid, err := h.checkClientAccess(groupID, request.ClientID, request.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check client access"})
		return
	}
	if !valid {
		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		return
	}

	err = h.groupUserApp.DeleteGroupUserByGroupID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group users by group ID"})
		return
	}

	response := clienthandlerdto.DeleteGroupUserByGroupIDResponse{
		Success: true,
	}

	c.JSON(http.StatusOK, response)
}
