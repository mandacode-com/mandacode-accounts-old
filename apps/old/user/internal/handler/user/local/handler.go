package localhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	localuserapp "mandacode.com/accounts/user/internal/app/user/local"
	localhandlerdto "mandacode.com/accounts/user/internal/handler/user/local/dto"
)

type LocalUserHandler struct {
	localUserApp localuserapp.LocalUserApp
	validator    *validator.Validate
	uidHeader    string
}

func NewLocalUserHandler(
	localUserApp localuserapp.LocalUserApp,
	validator *validator.Validate,
	uidHeader string,
) (*LocalUserHandler, error) {
	return &LocalUserHandler{
		localUserApp: localUserApp,
		validator:    validator,
		uidHeader:    uidHeader,
	}, nil
}

func (h *LocalUserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.CreateUser)
	rg.GET("/verify", h.VerifyUserEmail)
	rg.PATCH("/password", h.UpdatePassword)
}

func (h *LocalUserHandler) CreateUser(c *gin.Context) {
	var req localhandlerdto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	userID, err := h.localUserApp.CreateUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := localhandlerdto.CreateUserResponse{
		UserID: userID,
	}
	c.JSON(http.StatusCreated, response)
}

func (h *LocalUserHandler) VerifyUserEmail(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	verificationToken := c.Query("token")
	if verificationToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "verification token is required"})
		return
	}

	if err := h.localUserApp.VerifyUserEmail(c.Request.Context(), userUUID, verificationToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user email"})
		return
	}

	resposne := localhandlerdto.VerifyUserEmailResponse{
		Message: "Email verification successful",
	}
	c.JSON(http.StatusOK, resposne)
}

func (h *LocalUserHandler) UpdatePassword(c *gin.Context) {
	userID := c.GetHeader(h.uidHeader)
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	var req localhandlerdto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	if err := h.localUserApp.UpdatePassword(c.Request.Context(), userUUID, req.CurrentPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	response := localhandlerdto.UpdatePasswordResponse{
		Message: "Password updated successfully",
	}
	c.JSON(http.StatusOK, response)
}
