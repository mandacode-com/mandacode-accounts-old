package main

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/profile/cmd/server/grpc"
	httpserver "mandacode.com/accounts/profile/cmd/server/http"
	"mandacode.com/accounts/profile/internal/config"
	"mandacode.com/accounts/profile/internal/infra/database"
	"mandacode.com/accounts/profile/internal/infra/database/repository"
	"mandacode.com/accounts/profile/internal/infra/service"
	"mandacode.com/lib/server/server"
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
	repo := repository.NewProfileRepository(entClient)
	profileService := service.NewProfileService(repo)

	grpcServerRegisterer := grpcserver.NewGRPCRegisterer(profileService, logger)
	httpServerRegisterer := httpserver.NewHTTPRegisterer(
		profileService,
		logger,
		validator,
	)

	grpcServer, err := grpcserver.NewGRPCServer(cfg.GRPCPort, logger, grpcServerRegisterer, []string{
		"profile.v1.ProfileService",
	})
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}
	httpServer, err := httpserver.NewHTTPServer(cfg.HTTPPort, logger, httpServerRegisterer)
	if err != nil {
		logger.Fatal("failed to create HTTP server", zap.Error(err))
	}

	serverManager := &server.ServerManager{
		Servers: []server.Server{grpcServer, httpServer},
	}

	serverManager.Run()

	logger.Info("servers are running", zap.Int("grpc_port", cfg.GRPCPort), zap.Int("http_port", cfg.HTTPPort))
}
