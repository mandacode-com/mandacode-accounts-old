package main

import (
	"mandacode.com/accounts/auth/ent/oauthuser"
	"mandacode.com/accounts/auth/internal/app/login"
	"mandacode.com/accounts/auth/internal/app/user"
	"mandacode.com/accounts/auth/internal/config"
	"mandacode.com/accounts/auth/internal/database"
	"mandacode.com/accounts/auth/internal/database/repository"
	oauthdomain "mandacode.com/accounts/auth/internal/domain/service/oauth"
	loginhandler "mandacode.com/accounts/auth/internal/handler/login"
	userhandler "mandacode.com/accounts/auth/internal/handler/user"
	oauthsvc "mandacode.com/accounts/auth/internal/infra/oauth"
	"mandacode.com/accounts/auth/internal/infra/token"
	loginsvc "mandacode.com/accounts/auth/internal/service/login"
	usersvc "mandacode.com/accounts/auth/internal/service/user"
	localloginv1 "mandacode.com/accounts/proto/auth/login/local/v1"
	oauthloginv1 "mandacode.com/accounts/proto/auth/login/oauth/v1"
	localuserv1 "mandacode.com/accounts/proto/auth/user/local/v1"
	oauthuserv1 "mandacode.com/accounts/proto/auth/user/oauth/v1"

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
	localRepo := repository.NewLocalUserRepository(entClient)
	oauthRepo := repository.NewOAuthUserRepository(entClient)

	// Token service
	tokenService := token.New(tokenClient)

	localAuthSvc := loginsvc.NewLocalLoginService(localRepo)
	oauthAuthSvc := loginsvc.NewOAuthLoginService(oauthRepo)
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
	localLoginApp := login.NewLocalLoginApp(tokenService, localAuthSvc)
	oauthLoginApp := login.NewOAuthLoginApp(oauthServiceMap, tokenService, oauthAuthSvc)
	localUserApp := user.NewLocalUserApp(localUserSvc)
	oauthUserApp := user.NewOAuthUserApp(oauthServiceMap, oauthUserSvc)

	// Handlers
	localLoginHandler, err := loginhandler.NewLocalLoginHandler(localLoginApp, logger)
	if err != nil {
		return err
	}
	oauthLoginHandler, err := loginhandler.NewOAuthLoginHandler(oauthLoginApp, logger)
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
	localloginv1.RegisterLocalLoginServiceServer(server, localLoginHandler)
	oauthloginv1.RegisterOAuthLoginServiceServer(server, oauthLoginHandler)

	return nil
}
