package authhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	oauthlogin "mandacode.com/accounts/mobile-auth/internal/app/login/oauth"
	"mandacode.com/accounts/mobile-auth/internal/domain/model/provider"
	authhandlerdto "mandacode.com/accounts/mobile-auth/internal/handler/auth/dto"
)

type OAuthAuthHandler struct {
	oauthLoginApp oauthlogin.OAuthLoginApp
	validator     *validator.Validate
}

func NewOAuthAuthHandler(oauthLoginApp oauthlogin.OAuthLoginApp, validator *validator.Validate) (*OAuthAuthHandler, error) {
	if oauthLoginApp == nil {
		return nil, errors.New("oauthLoginApp cannot be nil")
	}
	if validator == nil {
		return nil, errors.New("validator cannot be nil")
	}

	return &OAuthAuthHandler{
		oauthLoginApp: oauthLoginApp,
		validator:     validator,
	}, nil
}

func (h *OAuthAuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/o/login/:provider", h.Login)
}

func (h *OAuthAuthHandler) Login(c *gin.Context) {
	providerStr := c.Param("provider")
	if providerStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	provider, err := provider.FromString(providerStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider", "details": err.Error()})
		return
	}

	// Get Access Token from header and parse it
	bearerAccessToken := c.GetHeader("Authorization")
	if bearerAccessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		return
	}
	accessToken := bearerAccessToken[len("Bearer "):]
	if accessToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access token is required"})
		return
	}

	loginToken, err := h.oauthLoginApp.Login(provider, accessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider", "details": err.Error()})
		return
	}

	response := authhandlerdto.OAuthLoginResponse{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}
