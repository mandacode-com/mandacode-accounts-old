package main

import (
	"log"

	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/token/cmd/server/grpc"
	token "mandacode.com/accounts/token/internal/app"
	"mandacode.com/accounts/token/internal/config"
	tokengen "mandacode.com/accounts/token/internal/infra/token"
	"mandacode.com/lib/server/server"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// tokenGenerator, err := tokengen.NewTokenGenerator()
	accesTokenGen, err := tokengen.NewTokenGeneratorByStr(
		cfg.AccessPrivateKey,
		cfg.AccessTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create access token generator", zap.Error(err))
	}
	refreshTokenGen, err := tokengen.NewTokenGeneratorByStr(
		cfg.RefreshPrivateKey,
		cfg.RefreshTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create refresh token generator", zap.Error(err))
	}
	emailVerificationTokenGen, err := tokengen.NewTokenGeneratorByStr(
		cfg.EmailVerificationPrivateKey,
		cfg.EmailVerificationTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create email verification token generator", zap.Error(err))
	}

	accessTokenApp := token.NewAccessTokenApp(accesTokenGen)
	refreshTokenApp := token.NewRefreshTokenApp(refreshTokenGen)
	emailVerificationTokenApp := token.NewEmailVerificationTokenApp(emailVerificationTokenGen)

	// Create the gRPC server
	registerer := grpcserver.NewGRPCRegisterer(
		accessTokenApp,
		refreshTokenApp,
		emailVerificationTokenApp,
		logger,
	)
	servingStatus := []string{
		"token.v1.TokenService",
	}
	grpcServer, err := grpcserver.NewGRPCServer(cfg.Port, logger, registerer, servingStatus)
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}

	serverManager := server.NewServerManager([]server.Server{grpcServer}, logger)

	serverManager.Run()
}
