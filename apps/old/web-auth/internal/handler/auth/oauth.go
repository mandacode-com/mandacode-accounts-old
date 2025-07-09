package authhandler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	oauthlogin "mandacode.com/accounts/web-auth/internal/app/login/oauth"
	"mandacode.com/accounts/web-auth/internal/domain/model/provider"
	authhandlerdto "mandacode.com/accounts/web-auth/internal/handler/auth/dto"
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
	rg.GET("/o/login/:provider", h.GetLoginURL)
	rg.GET("/o/callback/:provider", h.OAuthCallback)
}

func (h *OAuthAuthHandler) GetLoginURL(c *gin.Context) {
	providerStr := c.Param("provider")
	if providerStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}

	provider, err := provider.FromString(providerStr)

	loginURL, err := h.oauthLoginApp.GetLoginURL(provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider", "details": err.Error()})
		return
	}

	response := authhandlerdto.OAuthLoginURLResponse{
		URL: loginURL,
	}

	c.JSON(http.StatusOK, response)
}

func (h *OAuthAuthHandler) OAuthCallback(c *gin.Context) {
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

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	loginToken, err := h.oauthLoginApp.Login(provider, code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed", "details": err.Error()})
		return
	}

	session := sessions.Default(c)
	session.Set("refresh_token", loginToken.RefreshToken)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session", "details": err.Error()})
		return
	}

	response := authhandlerdto.OAuthCallbackResponse{
		AccessToken: loginToken.AccessToken,
	}

	c.JSON(http.StatusOK, response)
}
