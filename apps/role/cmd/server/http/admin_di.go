package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	groupapp "mandacode.com/accounts/role/internal/app/group"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
	permissionapp "mandacode.com/accounts/role/internal/app/permission"
	serviceapp "mandacode.com/accounts/role/internal/app/service"
	adminhandler "mandacode.com/accounts/role/internal/handler/http/admin"
	"mandacode.com/accounts/role/internal/middleware"
	"mandacode.com/lib/server/server"
)

type AdminRegisterer struct {
	PermissionApp   permissionapp.PermissionApp
	GroupApp        groupapp.GroupApp
	GroupUserApp    groupuserapp.GroupUserApp
	ServiceApp      serviceapp.ServiceApp
	Logger          *zap.Logger
	Validator       *validator.Validate
	UIDHeader       string
}

func NewAdminRegisterer(
	permissionApp permissionapp.PermissionApp,
	groupApp groupapp.GroupApp,
	groupUserApp groupuserapp.GroupUserApp,
	serviceApp serviceapp.ServiceApp,
	logger *zap.Logger,
	validator *validator.Validate,
	uidHeader string,
) server.HTTPRegisterer {
	return &AdminRegisterer{
		PermissionApp:   permissionApp,
		GroupApp:        groupApp,
		GroupUserApp:    groupUserApp,
		ServiceApp:      serviceApp,
		Logger:          logger,
		Validator:       validator,
		UIDHeader:       uidHeader,
	}
}

func (r *AdminRegisterer) Register(rg *gin.RouterGroup) error {
	adminMiddleware := middleware.NewAdminMiddleware(r.PermissionApp, r.UIDHeader)

	groupHandler, err := adminhandler.NewGroupHandler(r.GroupApp, r.Validator, adminMiddleware)
	if err != nil {
		r.Logger.Error("failed to create HTTP handler", zap.Error(err))
		return err
	}
	groupUserHandler, err := adminhandler.NewGroupUserHandler(r.GroupUserApp, r.Validator, adminMiddleware)
	if err != nil {
		r.Logger.Error("failed to create HTTP handler", zap.Error(err))
		return err
	}
	serviceHandler, err := adminhandler.NewServiceHandler(r.ServiceApp, r.Validator, adminMiddleware)
	if err != nil {
		r.Logger.Error("failed to create HTTP handler", zap.Error(err))
		return err
	}

	groupRoutes := rg.Group("/groups")
	groupHandler.RegisterRoutes(groupRoutes)

	groupUserRoutes := rg.Group("/gu")
	groupUserHandler.RegisterRoutes(groupUserRoutes)

	serviceRoutes := rg.Group("/services")
	serviceHandler.RegisterRoutes(serviceRoutes)

	r.Logger.Info("admin routes registered successfully")
	return nil
}
