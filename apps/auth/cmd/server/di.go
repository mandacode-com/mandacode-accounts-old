package main

import (
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/app/auth"
	"mandacode.com/accounts/auth/internal/app/user"
	"mandacode.com/accounts/auth/internal/config"
	"mandacode.com/accounts/auth/internal/database"
	authrepository "mandacode.com/accounts/auth/internal/database/repository"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	authhandler "mandacode.com/accounts/auth/internal/handler/auth"
	userhandler "mandacode.com/accounts/auth/internal/handler/user"
	oauthsvc "mandacode.com/accounts/auth/internal/infra/oauth"
	"mandacode.com/accounts/auth/internal/infra/token"
	authsvc "mandacode.com/accounts/auth/internal/service/auth"
	usersvc "mandacode.com/accounts/auth/internal/service/user"
	localauthv1 "mandacode.com/accounts/auth/proto/auth/local/v1"
	oauthauthv1 "mandacode.com/accounts/auth/proto/auth/oauth/v1"
	localuserv1 "mandacode.com/accounts/auth/proto/user/local/v1"
	oauthuserv1 "mandacode.com/accounts/auth/proto/user/oauth/v1"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func registerHandlers(server *grpc.Server, cfg *config.Config, logger *zap.Logger) error {
	// Initialize Ent client
	entClient, err := database.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		return err
	}

	// Token gRPC client
	tokenClient, _, err := token.NewGRPCClient(cfg.TokenServiceURL)
	if err != nil {
		return err
	}

	// Repositories
	localRepo := authrepository.NewLocalUserRepository(entClient)
	oauthRepo := authrepository.NewOAuthUserRepository(entClient)

	// Token service
	tokenService := token.New(tokenClient)

	localAuthSvc := authsvc.NewLocalAuthService(localRepo)
	oauthAuthSvc := authsvc.NewOAuthAuthService(oauthRepo)
	localUserSvc := usersvc.NewLocalUserService(localRepo)
	oauthUserSvc := usersvc.NewOAuthUserService(oauthRepo)

	// OAuth services
	googleOAuthSvc := oauthsvc.NewGoogleOAuthService()
	kakaoOAuthSvc := oauthsvc.NewKakaoOAuthService()
	naverOAuthSvc := oauthsvc.NewNaverOAuthService()

	oauthServiceMap := &map[oauthuser.Provider]oauthdomain.OAuthService{
		oauthuser.ProviderGoogle: googleOAuthSvc,
		oauthuser.ProviderKakao:  kakaoOAuthSvc,
		oauthuser.ProviderNaver:  naverOAuthSvc,
	}

	// App layer
	localAuthApp := auth.NewLocalAuthApp(tokenService, localAuthSvc)
	oauthAuthApp := auth.NewOAuthAuthApp(oauthServiceMap, tokenService, oauthAuthSvc)
	localUserApp := user.NewLocalUserApp(localUserSvc)
	oauthUserApp := user.NewOAuthUserApp(oauthServiceMap, oauthUserSvc)

	// Handlers
	localAuthHandler, err := authhandler.NewLocalAuthHandler(localAuthApp, logger)
	if err != nil {
		return err
	}
	oauthAuthHandler, err := authhandler.NewOAuthAuthHandler(oauthAuthApp, logger)
	if err != nil {
		return err
	}
	localUserHandler, err := userhandler.NewLocalUserHandler(localUserApp, logger)
	if err != nil {
		return err
	}
	oauthUserHandler, err := userhandler.NewOAuthUserHandler(oauthUserApp, logger)
	if err != nil {
		return err
	}

	// Register gRPC
	localuserv1.RegisterLocalUserServiceServer(server, localUserHandler)
	oauthuserv1.RegisterOAuthUserServiceServer(server, oauthUserHandler)
	localauthv1.RegisterLocalAuthServiceServer(server, localAuthHandler)
	oauthauthv1.RegisterOAuthAuthServiceServer(server, oauthAuthHandler)

	return nil
}
