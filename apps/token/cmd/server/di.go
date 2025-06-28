package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	tokenv1 "mandacode.com/accounts/proto/token/v1"
	token "mandacode.com/accounts/token/internal/app"
	"mandacode.com/accounts/token/internal/config"
	"mandacode.com/accounts/token/internal/handler"
	"mandacode.com/accounts/token/internal/infra/service"
)

func registerHandlers(
	server *grpc.Server,
	cfg *config.Config,
	logger *zap.Logger,
) error {
	accessTokenGen, err := service.NewTokenGeneratorByStr(
		cfg.AccessPrivateKey,
		cfg.AccessTokenDuration,
	)
	if err != nil {
		return err
	}
	refreshTokenGen, err := service.NewTokenGeneratorByStr(
		cfg.RefreshPrivateKey,
		cfg.RefreshTokenDuration,
	)
	if err != nil {
		return err
	}
	emailVerificationTokenGen, err := service.NewTokenGeneratorByStr(
		cfg.EmailVerificationPrivateKey,
		cfg.EmailVerificationTokenDuration,
	)
	if err != nil {
		return err
	}

	// Create the token service with the JWT generator
	accessTokenApp := token.NewAccessTokenApp(accessTokenGen)
	refreshTokenApp := token.NewRefreshTokenApp(refreshTokenGen)
	emailVerificationTokenApp := token.NewEmailVerificationTokenApp(emailVerificationTokenGen)

	tokenHandler, err := handler.NewTokenHandler(
		accessTokenApp,
		refreshTokenApp,
		emailVerificationTokenApp,
		logger,
	)
	if err != nil {
		return err
	}

	tokenv1.RegisterTokenServiceServer(server, tokenHandler)

	return nil

}
