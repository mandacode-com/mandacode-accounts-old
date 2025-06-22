package main

import (
	"log"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"mandacode.com/accounts/token/internal/app/token"
	"mandacode.com/accounts/token/internal/config"
	"mandacode.com/accounts/token/internal/handler/health"
	"mandacode.com/accounts/token/internal/handler/token"
	"mandacode.com/accounts/token/internal/infra"
	healthProto "mandacode.com/accounts/token/proto/health/v1"
	tokenProto "mandacode.com/accounts/token/proto/token/v1"
)

func main() {

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	cfg := config.LoadConfig()

	// Load RSA keys from PEM files
	accessTokenGen, err := infra.NewTokenGeneratorByStr(
		cfg.AccessPrivateKey,
		cfg.AccessTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create access token generator", zap.Error(err))
	}
	refreshTokenGen, err := infra.NewTokenGeneratorByStr(
		cfg.RefreshPrivateKey,
		cfg.RefreshTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create refresh token generator", zap.Error(err))
	}
	emailVerificationTokenGen, err := infra.NewTokenGeneratorByStr(
		cfg.EmailVerificationPrivateKey,
		cfg.EmailVerificationTokenDuration,
	)
	if err != nil {
		logger.Fatal("failed to create email verification token generator", zap.Error(err))
	}

	// Create the token service with the JWT generator
	tokenService := token.NewTokenService(accessTokenGen, refreshTokenGen, emailVerificationTokenGen)

	// Set up the gRPC server and register the JWT service handler
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Fatal("failed to listen on port 50051", zap.Error(err))
	}

	// Create a new gRPC server and register the JWT service
	grpcServer := grpc.NewServer()

	// Register the token service
	tokenHandler := tokenhandler.NewTokenHandler(tokenService, logger)
	tokenProto.RegisterTokenServiceServer(grpcServer, tokenHandler)

	// Register the health service
	healthHandler := healthhandler.NewHealthHandler(logger)
	healthProto.RegisterHealthServiceServer(grpcServer, healthHandler)

	logger.Info("Token gRPC service running", zap.String("address", ":50051"))
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve gRPC server", zap.Error(err))
	}
}
