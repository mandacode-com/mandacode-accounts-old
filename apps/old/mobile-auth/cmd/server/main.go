package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	httpserver "mandacode.com/accounts/mobile-auth/cmd/server/http"
	locallogin "mandacode.com/accounts/mobile-auth/internal/app/login/local"
	oauthlogin "mandacode.com/accounts/mobile-auth/internal/app/login/oauth"
	"mandacode.com/accounts/mobile-auth/internal/app/token"
	"mandacode.com/accounts/mobile-auth/internal/config"
	"mandacode.com/accounts/mobile-auth/internal/infra/auth"
	tokenmgr "mandacode.com/accounts/mobile-auth/internal/infra/token"
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

	// Initialize auth clients
	localAuthClient, _, err := auth.NewLocalLoginClient(cfg.AuthServiceURL)
	if err != nil {
		logger.Fatal("failed to create local auth client", zap.Error(err))
	}
	oauthAuthClient, _, err := auth.NewOAuthLoginClient(cfg.AuthServiceURL)
	if err != nil {
		logger.Fatal("failed to create OAuth auth client", zap.Error(err))
	}
	tokenClient, _, err := tokenmgr.NewTokenClient(cfg.AuthServiceURL)
	if err != nil {
		logger.Fatal("failed to create token client", zap.Error(err))
	}

	authenticator := auth.NewAuthenticator(localAuthClient, oauthAuthClient)
	tokenManager := tokenmgr.NewTokenManager(tokenClient)

	localLoginApp := locallogin.NewLocalLoginApp(authenticator)
	oauthLoginApp := oauthlogin.NewOAuthLoginApp(authenticator)
	tokenApp := token.NewTokenApp(tokenManager)

	authHandlerRegisterer := httpserver.NewAuthHandlerRegisterer(localLoginApp, oauthLoginApp, logger, validator)
	tokenHandlerRegisterer := httpserver.NewTokenHandlerRegisterer(tokenApp, logger, validator)

	httpServer, err := httpserver.NewHTTPServer(cfg.Port, logger, authHandlerRegisterer, tokenHandlerRegisterer)
	if err != nil {
		logger.Fatal("failed to create HTTP server", zap.Error(err))
	}

	serverManager := server.NewServerManager([]server.Server{httpServer}, logger)

	serverManager.Run()
}
