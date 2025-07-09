package tokenhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mandacode.com/accounts/mobile-auth/internal/app/token"
	tokenhandlerdto "mandacode.com/accounts/mobile-auth/internal/handler/token/dto"
)

type TokenHandler struct {
	tokenApp  token.TokenApp
	validator validator.Validate
}

func NewTokenHandler(tokenApp token.TokenApp, validate *validator.Validate) (*TokenHandler, error) {
	if tokenApp == nil {
		return nil, errors.New("tokenApp cannot be nil")
	}
	if validate == nil {
		return nil, errors.New("validator cannot be nil")
	}

	return &TokenHandler{
		tokenApp:  tokenApp,
		validator: *validate,
	}, nil
}

func (h *TokenHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/refresh", h.RefreshToken)
	rg.GET("/verify", h.VerifyAccessToken)
}

func (h *TokenHandler) RefreshToken(c *gin.Context) {
	// Extract the refresh token from the body
	var req tokenhandlerdto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	result, err := h.tokenApp.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token", "details": err.Error()})
		return
	}

	response := tokenhandlerdto.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}

func (h *TokenHandler) VerifyAccessToken(c *gin.Context) {
	// Extract the access token from header
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		return
	}
	// Remove "Bearer " prefix if present
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}
	result, err := h.tokenApp.VerifyAccessToken(c.Request.Context(), accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token", "details": err.Error()})
		return
	}

	response := tokenhandlerdto.VerifyAccessTokenResponse{
		Valid:  result.Valid,
		UserID: result.UserID,
	}

	c.JSON(http.StatusOK, response)
}
