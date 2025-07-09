package adminhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	groupapp "mandacode.com/accounts/role/internal/app/group"
	adminhandlerdto "mandacode.com/accounts/role/internal/handler/http/admin/dto"
)

type GroupHandler struct {
	groupApp        groupapp.GroupApp
	validator       *validator.Validate
	adminMiddleware gin.HandlerFunc
}

// NewGroupHandler creates a new GroupHandler instance.
func NewGroupHandler(
	groupApp groupapp.GroupApp,
	validator *validator.Validate,
	adminMiddleware gin.HandlerFunc,
) (*GroupHandler, error) {
	return &GroupHandler{
		groupApp:        groupApp,
		validator:       validator,
		adminMiddleware: adminMiddleware,
	}, nil
}

// RegisterRoutes registers the routes for the GroupHandler.
func (h *GroupHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.Use(h.adminMiddleware)
	rg.POST("/", h.CreateGroup)
	rg.GET("/", h.GetAllGroups)
	rg.GET("/:id", h.GetGroupByID)
	rg.PUT("/:id", h.UpdateGroup)
	rg.DELETE("/:id", h.DeleteGroup)
}

// CreateGroup handles the creation of a new group.
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var dto adminhandlerdto.CreateGroupRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	group, err := h.groupApp.CreateGroup(dto.Name, dto.ServiceID, dto.IsActive, dto.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create group"})
		return
	}

	response := adminhandlerdto.CreateGroupResponse{
		Group: group,
	}

	c.JSON(http.StatusCreated, response)
}

// GetGroupByID retrieves a group by its ID.
func (h *GroupHandler) GetGroupByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	group, err := h.groupApp.GetGroupByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get group"})
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	response := adminhandlerdto.GetGroupByIDResponse{
		Group: group,
	}

	c.JSON(http.StatusOK, response)
}

// GetAllGroups retrieves all groups.
func (h *GroupHandler) GetAllGroups(c *gin.Context) {
	serviceIDStr := c.Query("service_id")
	var serviceID uuid.UUID
	if serviceIDStr != "" {
		var err error
		serviceID, err = uuid.Parse(serviceIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service ID format"})
			return
		}
	}

	groups, err := h.groupApp.GetGroupsByServiceID(serviceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get groups"})
		return
	}

	if len(groups) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no groups found"})
		return
	}

	response := adminhandlerdto.GetAllGroupsResponse{
		Groups: groups,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateGroup updates an existing group.
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	var dto adminhandlerdto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	group, err := h.groupApp.UpdateGroup(id, dto.Name, dto.ServiceID, dto.IsActive, dto.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update group"})
		return
	}

	response := adminhandlerdto.UpdateGroupResponse{
		Group: group,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteGroup deletes a group by its ID.
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID format"})
		return
	}

	if err := h.groupApp.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete group"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
