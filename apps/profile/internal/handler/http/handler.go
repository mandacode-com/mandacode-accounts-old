package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/app/profile"
	"mandacode.com/accounts/profile/internal/handler/http/request"
)

type ProfileHTTPHandler struct {
	getUC     profile.GetProfileUsecase
	updateUC  profile.UpdateProfileUsecase
	validator *validator.Validate
}

func NewProfileHTTPHandler(
	getUC profile.GetProfileUsecase,
	updateUC profile.UpdateProfileUsecase,
	validator *validator.Validate,
) (*ProfileHTTPHandler, error) {
	return &ProfileHTTPHandler{
		getUC:     getUC,
		updateUC:  updateUC,
		validator: validator,
	}, nil
}

func (h *ProfileHTTPHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/profile/:user_id", h.GetProfile)
	rg.PUT("/profile/:user_id", h.UpdateProfile)
}

func (h *ProfileHTTPHandler) GetProfile(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(userID)

	profile, err := h.getUC.GetProfile(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"profile": profile,
	})
}

func (h *ProfileHTTPHandler) UpdateProfile(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	var profileUpdate request.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&profileUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(profileUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	profile, err := h.updateUC.UpdateProfile(userUUID, profileUpdate.Email, profileUpdate.DisplayName, profileUpdate.Bio, profileUpdate.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": profile,
	})
}
