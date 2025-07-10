package authhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	locallogin "mandacode.com/accounts/mobile-auth/internal/app/login/local"
	authhandlerdto "mandacode.com/accounts/mobile-auth/internal/handler/auth/dto"
)

type LocalAuthHandler struct {
	localloginApp locallogin.LocalLoginApp
	validator     *validator.Validate
}

func NewLocalAuthHandler(localloginApp locallogin.LocalLoginApp, validator *validator.Validate) (*LocalAuthHandler, error) {
	if localloginApp == nil {
		return nil, errors.New("localloginApp cannot be nil")
	}
	if validator == nil {
		return nil, errors.New("validator cannot be nil")
	}

	return &LocalAuthHandler{
		localloginApp: localloginApp,
		validator:     validator,
	}, nil
}

func (h *LocalAuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.Login)
}

func (h *LocalAuthHandler) Login(c *gin.Context) {
	var req authhandlerdto.LocalLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed", "details": err.Error()})
		return
	}

	loginToken, err := h.localloginApp.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed", "details": err.Error()})
		return
	}

	response := authhandlerdto.LocalLoginResponse{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}

	c.JSON(http.StatusOK, response)
}
