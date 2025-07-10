package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"

	localuserapp "mandacode.com/accounts/user/internal/app/user/local"
	oauthuserapp "mandacode.com/accounts/user/internal/app/user/oauth"
	userapp "mandacode.com/accounts/user/internal/app/user/user"
	localhandler "mandacode.com/accounts/user/internal/handler/user/local"
	oauthhandler "mandacode.com/accounts/user/internal/handler/user/oauth"
	userhandler "mandacode.com/accounts/user/internal/handler/user/user"
)

type HTTPRegisterer struct {
	UserApp      userapp.UserApp
	LocalUserApp localuserapp.LocalUserApp
	OAuthUserApp oauthuserapp.OAuthUserApp
	Logger       *zap.Logger
	Validator    *validator.Validate
	UIDHeader    string
}

func NewHTTPRegisterer(userApp userapp.UserApp, localUserApp localuserapp.LocalUserApp, oauthUserApp oauthuserapp.OAuthUserApp, logger *zap.Logger, validator *validator.Validate, uidHeader string) server.HTTPRegisterer {
	return &HTTPRegisterer{
		UserApp:      userApp,
		LocalUserApp: localUserApp,
		OAuthUserApp: oauthUserApp,
		Logger:       logger,
		Validator:    validator,
		UIDHeader:    uidHeader,
	}
}

func (r *HTTPRegisterer) Register(rg *gin.RouterGroup) error {
	// Register the user handlers
	userGroup := rg.Group("/user")
	userHandler, err := userhandler.NewUserHandler(r.UserApp, r.Validator, r.UIDHeader)
	if err != nil {
		r.Logger.Error("failed to create User HTTP handler", zap.Error(err))
		return err
	}
	userHandler.RegisterRoutes(userGroup)

	// Register the local user handlers
	localUserGroup := rg.Group("/local/user")
	localUserHandler, err := localhandler.NewLocalUserHandler(r.LocalUserApp, r.Validator, r.UIDHeader)
	if err != nil {
		r.Logger.Error("failed to create Local User HTTP handler", zap.Error(err))
		return err
	}
	localUserHandler.RegisterRoutes(localUserGroup)

	// Register the OAuth user handlers
	oauthUserGroup := rg.Group("/oauth/user")
	oauthUserHandler, err := oauthhandler.NewOAuthUserHandler(r.OAuthUserApp, r.Validator, r.UIDHeader)
	if err != nil {
		r.Logger.Error("failed to create OAuth User HTTP handler", zap.Error(err))
		return err
	}
	oauthUserHandler.RegisterRoutes(oauthUserGroup)

	return nil
}
