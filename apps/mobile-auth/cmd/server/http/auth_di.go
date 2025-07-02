package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	locallogin "mandacode.com/accounts/mobile-auth/internal/app/login/local"
	oauthlogin "mandacode.com/accounts/mobile-auth/internal/app/login/oauth"
	authhandler "mandacode.com/accounts/mobile-auth/internal/handler/auth"
	"mandacode.com/lib/server/server"
)

type AuthHandlerRegisterer struct {
	LocalLoginApp locallogin.LocalLoginApp
	OAuthLoginApp oauthlogin.OAuthLoginApp
	Logger        *zap.Logger
	Validator     *validator.Validate
}

func NewAuthHandlerRegisterer(localLoginApp locallogin.LocalLoginApp, oauthLoginApp oauthlogin.OAuthLoginApp, logger *zap.Logger, validator *validator.Validate) server.HTTPRegisterer {
	return &AuthHandlerRegisterer{
		LocalLoginApp: localLoginApp,
		OAuthLoginApp: oauthLoginApp,
		Logger:        logger,
		Validator:     validator,
	}
}

func (r *AuthHandlerRegisterer) Register(rg *gin.RouterGroup) error {
	localAuthHandler, err := authhandler.NewLocalAuthHandler(r.LocalLoginApp, r.Validator)
	if err != nil {
		r.Logger.Error("failed to create local auth handler", zap.Error(err))
		return err
	}
	oauthAuthHandler, err := authhandler.NewOAuthAuthHandler(r.OAuthLoginApp, r.Validator)
	if err != nil {
		r.Logger.Error("failed to create OAuth auth handler", zap.Error(err))
		return err
	}

	localAuthHandler.RegisterRoutes(rg)
	oauthAuthHandler.RegisterRoutes(rg)

	return nil
}
