package oauthhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	oauthuserapp "mandacode.com/accounts/user/internal/app/user/oauth"
	oauthhandlerdto "mandacode.com/accounts/user/internal/handler/user/oauth/dto"
	"mandacode.com/accounts/user/internal/util"
)

type OAuthUserHandler struct {
	oauthUserApp oauthuserapp.OAuthUserApp
	validator    *validator.Validate
	uidHeader    string
}

func NewOAuthUserHandler(
	oauthUserApp oauthuserapp.OAuthUserApp,
	validator *validator.Validate,
	uidHeader string,
) (*OAuthUserHandler, error) {
	return &OAuthUserHandler{
		oauthUserApp: oauthUserApp,
		validator:    validator,
		uidHeader:    uidHeader,
	}, nil
}

func (h *OAuthUserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/:provider", h.CreateUser)
}

func (h *OAuthUserHandler) CreateUser(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}
	providerEnum := util.ConvertToOAuthProvider(provider)

	accessToken := c.GetHeader(h.uidHeader)

	userID, err := h.oauthUserApp.CreateUser(c.Request.Context(), providerEnum, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	response := oauthhandlerdto.CreateOAuthUserResponse{
		UserID: userID,
	}
	c.JSON(http.StatusCreated, response)
}
