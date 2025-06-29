package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	httpuc "mandacode.com/accounts/profile/internal/app/profile/http"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
	httphandler "mandacode.com/accounts/profile/internal/handler/http"
	"mandacode.com/lib/server/server"
)

type HTTPRegisterer struct {
	ProfileService svcdomain.ProfileService
	Logger         *zap.Logger
	Validator      *validator.Validate
}

func NewHTTPRegisterer(profileService svcdomain.ProfileService, logger *zap.Logger, validator *validator.Validate) server.HTTPRegisterer {
	return &HTTPRegisterer{
		ProfileService: profileService,
		Logger:         logger,
		Validator:      validator,
	}
}

func (r *HTTPRegisterer) Register(rg *gin.RouterGroup) error {
	updateProfileUC := httpuc.NewUpdateProfileUsecase(r.ProfileService)
	getProfileUC := httpuc.NewGetProfileUsecase(r.ProfileService)

	profileHandler, err := httphandler.NewProfileHTTPHandler(getProfileUC, updateProfileUC, r.Validator)
	if err != nil {
		r.Logger.Error("failed to create HTTP handler", zap.Error(err))
		return err
	}

	profileHandler.RegisterRoutes(rg)

	return nil
}
