package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"mandacode.com/accounts/profile/internal/app/profile"
	httphandlerdto "mandacode.com/accounts/profile/internal/handler/http/dto"
)

type ProfileHTTPHandler struct {
	app       profile.ProfileApp
	validator *validator.Validate
	uidHeader string
}

func NewProfileHTTPHandler(
	app profile.ProfileApp,
	validator *validator.Validate,
	uidHeader string,
) (*ProfileHTTPHandler, error) {
	return &ProfileHTTPHandler{
		app:       app,
		validator: validator,
		uidHeader: uidHeader,
	}, nil
}

func (h *ProfileHTTPHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/profile", h.GetProfile)
	rg.PUT("/profile", h.UpdateProfile)
}

func (h *ProfileHTTPHandler) GetProfile(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(userID)

	profile, err := h.app.GetProfile(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	response := httphandlerdto.GetProfileResponse{
		Profile: profile,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHTTPHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	var profileUpdate httphandlerdto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&profileUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(profileUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	profile, err := h.app.UpdateProfile(userUUID, profileUpdate.Email, profileUpdate.DisplayName, profileUpdate.Bio, profileUpdate.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	response := httphandlerdto.ProfileUpdateResponse{
		Profile: profile,
	}

	c.JSON(http.StatusOK, response)
}
