package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mandacode.com/accounts/web-auth/internal/app/token"
	tokenhandler "mandacode.com/accounts/web-auth/internal/handler/token"
	"mandacode.com/lib/server/server"
)

type TokenHandlerRegisterer struct {
	TokenApp  token.TokenApp
	Logger    *zap.Logger
	Validator *validator.Validate
}

func NewTokenHandlerRegisterer(tokenApp token.TokenApp, logger *zap.Logger, validator *validator.Validate) server.HTTPRegisterer {
	return &TokenHandlerRegisterer{
		TokenApp:  tokenApp,
		Logger:    logger,
		Validator: validator,
	}
}

func (r *TokenHandlerRegisterer) Register(rg *gin.RouterGroup) error {
	tokenHandler, err := tokenhandler.NewTokenHandler(r.TokenApp, r.Validator)
	if err != nil {
		r.Logger.Error("failed to create token handler", zap.Error(err))
		return err
	}

	tokenHandler.RegisterRoutes(rg)

	return nil
}
