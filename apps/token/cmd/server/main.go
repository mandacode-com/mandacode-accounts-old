package main

import (
	"log"
	"net"
	"time"

	"google.golang.org/grpc"

	"mandacode.com/accounts/token/internal/app"
	"mandacode.com/accounts/token/internal/config"
	"mandacode.com/accounts/token/internal/handler"
	"mandacode.com/accounts/token/internal/infra"
	healthProto "mandacode.com/accounts/token/proto/health/v1"
	tokenProto "mandacode.com/accounts/token/proto/token/v1"
)

func main() {
	cfg := config.LoadConfig()

	// Load RSA keys from PEM files
	accessTokenGen, err := infra.NewTokenGeneratorByStr(
		cfg.AccessPublicKey,
		cfg.AccessPrivateKey,
		"mandacode.com/accounts",
		time.Hour*24)
	if err != nil {
		log.Fatalf("failed to load access token private key: %v", err)
	}
	refreshTokenGen, err := infra.NewTokenGeneratorByStr(
		cfg.RefreshPublicKey,
		cfg.RefreshPrivateKey,
		"mandacode.com/accounts",
		time.Hour*24*30)
	if err != nil {
		log.Fatalf("failed to load refresh token private key: %v", err)
	}
	emailVerificationTokenGen, err := infra.NewTokenGeneratorByStr(
		cfg.EmailVerificationPublicKey,
		cfg.EmailVerificationPrivateKey,
		"mandacode.com/accounts/email-verification",
		time.Hour*24*7)
	if err != nil {
		log.Fatalf("failed to load email verification token private key: %v", err)
	}

	// Create the token service with the JWT generator
	tokenService := app.NewTokenService(accessTokenGen, refreshTokenGen, emailVerificationTokenGen)

	// Set up the gRPC server and register the JWT service handler
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server and register the JWT service
	grpcServer := grpc.NewServer()

	// Register the token service
	tokenHandler := handler.NewTokenHandler(tokenService)
	tokenProto.RegisterTokenServiceServer(grpcServer, tokenHandler)

	// Register the health service
	healthHandler := handler.NewHealthHandler()
	healthProto.RegisterHealthServiceServer(grpcServer, healthHandler)

	log.Println("JWT gRPC service running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
