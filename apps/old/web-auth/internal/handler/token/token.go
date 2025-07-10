package tokenhandler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"mandacode.com/accounts/web-auth/internal/app/token"
	tokenhandlerdto "mandacode.com/accounts/web-auth/internal/handler/token/dto"
)

type TokenHandler struct {
	tokenApp token.TokenApp
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
		tokenApp: tokenApp,
		validator: *validate,
	}, nil
}

func (h *TokenHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/refresh", h.RefreshToken)
	rg.GET("/verify", h.VerifyAccessToken)
}

func (h *TokenHandler) RefreshToken(c *gin.Context) {
	// Extract the refresh token from the session
	session := sessions.Default(c)
	refreshToken := session.Get("refresh_token")
	if refreshToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found in session"})
		return
	}

	result, err := h.tokenApp.RefreshToken(c.Request.Context(), refreshToken.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token", "details": err.Error()})
		return
	}

	session.Set("refresh_token", result.RefreshToken)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session", "details": err.Error()})
		return
	}

	response := tokenhandlerdto.RefreshTokenResponse{
		AccessToken: result.AccessToken,
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
