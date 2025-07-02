package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"mandacode.com/accounts/profile/internal/app/profile"
	httphandler "mandacode.com/accounts/profile/internal/handler/http"
	"mandacode.com/lib/server/server"
)

type HTTPRegisterer struct {
	ProfileApp profile.ProfileApp
	Logger     *zap.Logger
	Validator  *validator.Validate
	UIDHeader  string
}

func NewHTTPRegisterer(profileApp profile.ProfileApp, logger *zap.Logger, validator *validator.Validate, uidHeader string) server.HTTPRegisterer {
	return &HTTPRegisterer{
		ProfileApp: profileApp,
		Logger:     logger,
		Validator:  validator,
		UIDHeader:  uidHeader,
	}
}

func (r *HTTPRegisterer) Register(rg *gin.RouterGroup) error {
	profileHandler, err := httphandler.NewProfileHTTPHandler(r.ProfileApp, r.Validator, r.UIDHeader)

	if err != nil {
		r.Logger.Error("failed to create HTTP handler", zap.Error(err))
		return err
	}

	profileHandler.RegisterRoutes(rg)

	return nil
}
