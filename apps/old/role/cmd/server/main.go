package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/role/cmd/server/grpc"
	httpserver "mandacode.com/accounts/role/cmd/server/http"
	groupapp "mandacode.com/accounts/role/internal/app/group"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
	permissionapp "mandacode.com/accounts/role/internal/app/permission"
	serviceapp "mandacode.com/accounts/role/internal/app/service"
	"mandacode.com/accounts/role/internal/config"
	"mandacode.com/accounts/role/internal/infra/database"
	"mandacode.com/accounts/role/internal/infra/database/repository"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	defer logger.Sync()

	validator := validator.New()

	cfg, err := config.LoadConfig(validator)
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	entClient, err := database.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	groupRepo := repository.NewGroupRepository(entClient)
	groupUserRepo := repository.NewGroupUserRepository(entClient)
	serviceRepo := repository.NewServiceRepository(entClient)
	permissionRepo := repository.NewClientAccessRepository(entClient)

	groupApp := groupapp.NewGroupApp(groupRepo)
	groupUserApp := groupuserapp.NewGroupUserApp(groupUserRepo)
	serviceApp := serviceapp.NewServiceApp(serviceRepo)
	permissionApp := permissionapp.NewPermissionApp(groupUserRepo, permissionRepo, cfg.AdminGroupID)

	grpcServerRegisterer := grpcserver.NewGRPCRegisterer(groupUserApp, logger)
	clientRegisterer := httpserver.NewClientRegisterer(
		groupUserApp,
		logger,
		validator,
		cfg.UIDHeader,
	)
	adminRegisterer := httpserver.NewAdminRegisterer(
		permissionApp,
		groupApp,
		groupUserApp,
		serviceApp,
		logger,
		validator,
		cfg.UIDHeader,
	)

	grpcServer, err := grpcserver.NewGRPCServer(cfg.GRPCPort, logger, grpcServerRegisterer, []string{
		"profile.v1.ProfileService",
	})
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}
	httpServer, err := httpserver.NewHTTPServer(cfg.HTTPPort, logger, clientRegisterer, adminRegisterer)
	if err != nil {
		logger.Fatal("failed to create HTTP server", zap.Error(err))
	}

	serverManager := server.NewServerManager([]server.Server{grpcServer, httpServer}, logger)

	serverManager.Run()

	logger.Info("servers are running", zap.Int("grpc_port", cfg.GRPCPort), zap.Int("http_port", cfg.HTTPPort))
}
