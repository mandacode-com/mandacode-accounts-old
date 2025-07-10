package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	groupapp "mandacode.com/accounts/role/internal/app/group"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
	permissionapp "mandacode.com/accounts/role/internal/app/permission"
	clienthandler "mandacode.com/accounts/role/internal/handler/http/client"
)

type ClientRegisterer struct {
	PermissionApp permissionapp.PermissionApp
	GroupApp      groupapp.GroupApp
	GroupUserApp  groupuserapp.GroupUserApp
	Logger        *zap.Logger
	Validator     *validator.Validate
	UIDHeader     string
}

func NewClientRegisterer(groupUserApp groupuserapp.GroupUserApp, logger *zap.Logger, validator *validator.Validate, uidHeader string) server.HTTPRegisterer {
	return &ClientRegisterer{
		GroupUserApp: groupUserApp,
		Logger:       logger,
		Validator:    validator,
		UIDHeader:    uidHeader,
	}
}

func (r *ClientRegisterer) Register(rg *gin.RouterGroup) error {
	clientGroupUserHandler, err := clienthandler.NewClientGroupUserHandler(r.PermissionApp, r.GroupApp, r.GroupUserApp, r.Validator)

	if err != nil {
		r.Logger.Error("failed to create HTTP handler", zap.Error(err))
		return err
	}

	clientGroupUserHandler.RegisterRoutes(rg)

	return nil
}
