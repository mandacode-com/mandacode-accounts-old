package main

import (
	sessionredis "github.com/gin-contrib/sessions/redis"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	httpserver "mandacode.com/accounts/web-auth/cmd/server/http"
	locallogin "mandacode.com/accounts/web-auth/internal/app/login/local"
	oauthlogin "mandacode.com/accounts/web-auth/internal/app/login/oauth"
	"mandacode.com/accounts/web-auth/internal/app/token"
	"mandacode.com/accounts/web-auth/internal/config"
	"mandacode.com/accounts/web-auth/internal/domain/model/provider"
	oauthcodedomain "mandacode.com/accounts/web-auth/internal/domain/oauthcode"
	"mandacode.com/accounts/web-auth/internal/infra/auth"
	"mandacode.com/accounts/web-auth/internal/infra/oauthcode"
	tokenmgr "mandacode.com/accounts/web-auth/internal/infra/token"
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

	codeStore := redis.NewClient(&redis.Options{
		Addr:     cfg.CodeStore.Address,
		Password: cfg.CodeStore.Password,
		DB:       cfg.CodeStore.DB,
	})
	sessionStore, err := sessionredis.NewStore(cfg.SessionStore.DB, "tcp", cfg.SessionStore.Address, "", cfg.CodeStore.Password)
	if err != nil {
		logger.Fatal("failed to create session store", zap.Error(err))
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

	googleCode := oauthcode.NewGoogleOAuthCode(cfg.GoogleOAuthConfig.ClientID, cfg.GoogleOAuthConfig.ClientSecret, cfg.GoogleOAuthConfig.RedirectURL)
	kakaoCode := oauthcode.NewKakaoOAuthCode(cfg.KakaoOAuthConfig.ClientID, cfg.KakaoOAuthConfig.ClientSecret, cfg.KakaoOAuthConfig.RedirectURL)
	naverCode := oauthcode.NewNaverOAuthCode(cfg.NaverOAuthConfig.ClientID, cfg.NaverOAuthConfig.ClientSecret, cfg.NaverOAuthConfig.RedirectURL)
	oauthCodeMap := map[provider.Provider]oauthcodedomain.OAuthCode{
		provider.ProviderGoogle: googleCode,
		provider.ProviderKakao:  kakaoCode,
		provider.ProviderNaver:  naverCode,
	}

	localLoginApp := locallogin.NewLocalLoginApp(codeStore, authenticator, cfg.CodeTTL)
	oauthLoginApp := oauthlogin.NewOAuthLoginApp(oauthCodeMap, authenticator)
	tokenApp := token.NewTokenApp(tokenManager)

	authHandlerRegisterer := httpserver.NewAuthHandlerRegisterer(localLoginApp, oauthLoginApp, logger, validator)
	tokenHandlerRegisterer := httpserver.NewTokenHandlerRegisterer(tokenApp, logger, validator)

	httpServer, err := httpserver.NewHTTPServer(cfg.Port, logger, authHandlerRegisterer, tokenHandlerRegisterer, sessionStore)
	if err != nil {
		logger.Fatal("failed to create HTTP server", zap.Error(err))
	}

	serverManager := server.NewServerManager([]server.Server{httpServer}, logger)

	serverManager.Run()
}
