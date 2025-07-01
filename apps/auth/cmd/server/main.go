package main

import (
	"log"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	grpcserver "mandacode.com/accounts/auth/cmd/server/grpc"
	"mandacode.com/accounts/auth/ent/oauthuser"
	locallogin "mandacode.com/accounts/auth/internal/app/login/local"
	oauthlogin "mandacode.com/accounts/auth/internal/app/login/oauth"
	"mandacode.com/accounts/auth/internal/app/token"
	localuser "mandacode.com/accounts/auth/internal/app/user/local"
	oauthuserapp "mandacode.com/accounts/auth/internal/app/user/oauth"
	"mandacode.com/accounts/auth/internal/config"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/oauth"
	"mandacode.com/accounts/auth/internal/infra/database"
	"mandacode.com/accounts/auth/internal/infra/database/repository"
	oauthprovider "mandacode.com/accounts/auth/internal/infra/oauth"
	tokenprovider "mandacode.com/accounts/auth/internal/infra/token"
	"mandacode.com/lib/server/server"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}
	validate := validator.New()

	entClient, err := database.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create Ent client", zap.Error(err))
	}

	localUserRepo := repository.NewLocalUserRepository(entClient)
	oauthUserRepo := repository.NewOAuthUserRepository(entClient)

	// Token Provider
	tokenClient, _, err := tokenprovider.NewGRPCClient(cfg.TokenServiceURL)
	if err != nil {
		logger.Fatal("failed to create token gRPC client", zap.Error(err))
	}
	tokenProvider := tokenprovider.NewTokenProvider(tokenClient)

	// OAuth Provider
	googleProvider := oauthprovider.NewGoogleOAuthProvider(validate)
	naverProvider := oauthprovider.NewNaverOAuthProvider(validate)
	kakaoProvider := oauthprovider.NewKakaoOAuthProvider(validate)
	oauthProviders := &map[oauthuser.Provider]oauthdomain.OAuthProvider{
		oauthuser.ProviderGoogle: googleProvider,
		oauthuser.ProviderNaver:  naverProvider,
		oauthuser.ProviderKakao:  kakaoProvider,
	}

	// Login App
	localLoginApp := locallogin.NewLocalLoginApp(
		tokenProvider,
		localUserRepo,
	)
	oauthLoginApp := oauthlogin.NewOAuthLoginApp(
		tokenProvider,
		oauthProviders,
		oauthUserRepo,
	)

	// User App
	localUserApp := localuser.NewLocalUserApp(
		localUserRepo,
	)
	oauthUserApp := oauthuserapp.NewOAuthUserApp(
		oauthProviders,
		oauthUserRepo,
	)

	// Token App
	tokenApp := token.NewTokenApp(
		tokenProvider,
	)

	// Register the applications
	registerer := grpcserver.NewGRPCRegisterer(
		localLoginApp,
		oauthLoginApp,
		localUserApp,
		oauthUserApp,
		tokenApp,
		logger,
	)
	// Create the gRPC server
	servingServices := []string{
		"auth.v1.LocalLoginService",
		"auth.v1.OAuthLoginService",
		"auth.v1.LocalUserService",
		"auth.v1.OAuthUserService",
	}
	grpcServer, err := grpcserver.NewGRPCServer(cfg.Port, logger, registerer, servingServices)
	if err != nil {
		logger.Fatal("failed to create gRPC server", zap.Error(err))
	}

	// Run the gRPC server
	serverManager := server.NewServerManager([]server.Server{grpcServer}, logger)

	serverManager.Run()

	logger.Info("gRPC server is running", zap.Int("port", cfg.Port))
}
