package userhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	userapp "mandacode.com/accounts/user/internal/app/user/user"
	userhandlerdto "mandacode.com/accounts/user/internal/handler/user/user/dto"
)

type UserHandler struct {
	userApp   userapp.UserApp
	validator *validator.Validate
	uidHeader string
}

func NewUserHandler(
	userApp userapp.UserApp,
	validator *validator.Validate,
	uidHeader string,
) (*UserHandler, error) {
	return &UserHandler{
		userApp:   userApp,
		validator: validator,
		uidHeader: uidHeader,
	}, nil
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.PATCH("/archive", h.ArchiveUser)
	rg.DELETE("/", h.DeleteUser)
	rg.PATCH("/active-status", h.UpdateActiveStatus)
	rg.PATCH("/verified-status", h.UpdateVerifiedStatus)
}

func (h *UserHandler) ArchiveUser(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userApp.ArchiveUser(c.Request.Context(), userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.userApp.DeleteUser(c.Request.Context(), userUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UpdateActiveStatus(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req userhandlerdto.UpdateActiveStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.userApp.UpdateActiveStatus(c.Request.Context(), userUUID, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update active status"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) UpdateVerifiedStatus(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req userhandlerdto.UpdateVerifiedStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.userApp.UpdateVerifiedStatus(c.Request.Context(), userUUID, req.IsVerified); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update verified status"})
		return
	}

	c.Status(http.StatusNoContent)
}
