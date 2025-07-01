package authhandler

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	locallogin "mandacode.com/accounts/web-auth/internal/app/login/local"
	authhandlerdto "mandacode.com/accounts/web-auth/internal/handler/auth/dto"
	"net/http"
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
	rg.GET("/verify/:code", h.VerifyCode)
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

	code, err := h.localloginApp.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login failed", "details": err.Error()})
		return
	}

	response := authhandlerdto.LocalLoginResponse{
		Code: code,
	}

	c.JSON(http.StatusOK, response)
}

func (h *LocalAuthHandler) VerifyCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	loginToken, err := h.localloginApp.VerifyCode(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "verification failed", "details": err.Error()})
		return
	}

	// save the refresh token in the session
	session := sessions.Default(c)
	session.Set("refresh_token", loginToken.RefreshToken)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save session", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, authhandlerdto.LocalVerifyCodeResponse{
		AccessToken: loginToken.AccessToken,
	})
}
