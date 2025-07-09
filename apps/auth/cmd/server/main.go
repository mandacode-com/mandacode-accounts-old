package main

import (
	"context"
	"os"
	"os/signal"

	sessionredis "github.com/gin-contrib/sessions/redis"
	"github.com/go-playground/validator/v10"
	"github.com/mandacode-com/golib/server"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	httpserver "mandacode.com/accounts/auth/cmd/server/http"
	"mandacode.com/accounts/auth/config"
	"mandacode.com/accounts/auth/ent/oauthauth"
	handlerv1 "mandacode.com/accounts/auth/internal/handler/v1"
	dbinfra "mandacode.com/accounts/auth/internal/infra/database"
	"mandacode.com/accounts/auth/internal/infra/mailer"
	"mandacode.com/accounts/auth/internal/infra/oauthapi"
	tokeninfra "mandacode.com/accounts/auth/internal/infra/token"
	coderepo "mandacode.com/accounts/auth/internal/repository/code"
	dbrepository "mandacode.com/accounts/auth/internal/repository/database"
	tokenrepo "mandacode.com/accounts/auth/internal/repository/token"
	"mandacode.com/accounts/auth/internal/usecase/localauth"
	oauthusecase "mandacode.com/accounts/auth/internal/usecase/oauthauth"
	"mandacode.com/accounts/auth/internal/util"
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

	// Initialize Redis clients and session store
	loginCodeStore := redis.NewClient(&redis.Options{
		Addr:     cfg.LoginCodeStore.Address,
		Password: cfg.LoginCodeStore.Password,
		DB:       cfg.LoginCodeStore.DB,
	})
	emailCodeStore := redis.NewClient(&redis.Options{
		Addr:     cfg.EmailCodeStore.Address,
		Password: cfg.EmailCodeStore.Password,
		DB:       cfg.EmailCodeStore.DB,
	})
	sessionStore, err := sessionredis.NewStore(cfg.SessionStore.DB, "tcp", cfg.SessionStore.Address, "", cfg.SessionStore.Password, []byte(cfg.SessionStore.HashKey))
	if err != nil {
		logger.Fatal("failed to create session store", zap.Error(err))
	}

	// Initialize database and token clients
	dbClient, err := dbinfra.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	tokenClient, _, err := tokeninfra.NewTokenClient(cfg.TokenServiceAddr)
	if err != nil {
		logger.Fatal("failed to create token client", zap.Error(err))
	}

	// Initialize mailer
	mailWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.MailWriter.Address),
		Topic:                  cfg.MailWriter.Topic,
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}
	mailer := mailer.NewMailer(mailWriter)

	// Initialize OAuth APIs
	googleApi, err := oauthapi.NewGoogleAPI(cfg.GoogleOAuth.ClientID, cfg.GoogleOAuth.ClientSecret, cfg.GoogleOAuth.RedirectURL, validator)
	if err != nil {
		logger.Fatal("failed to create Google OAuth API", zap.Error(err))
	}
	naverApi, err := oauthapi.NewNaverAPI(cfg.NaverOAuth.ClientID, cfg.NaverOAuth.ClientSecret, cfg.NaverOAuth.RedirectURL, validator)
	if err != nil {
		logger.Fatal("failed to create Naver OAuth API", zap.Error(err))
	}
	kakaoApi, err := oauthapi.NewKakaoAPI(cfg.KakaoOAuth.ClientID, cfg.KakaoOAuth.ClientSecret, cfg.KakaoOAuth.RedirectURL, validator)
	if err != nil {
		logger.Fatal("failed to create Kakao OAuth API", zap.Error(err))
	}
	oauthApis := map[oauthauth.Provider]oauthapi.OAuthAPI{
		oauthauth.ProviderGoogle: googleApi,
		oauthauth.ProviderNaver:  naverApi,
		oauthauth.ProviderKakao:  kakaoApi,
	}

	// Initialize random code generators
	emailCodeGenerator := util.NewRandomGenerator(32)
	loginCodeGenerator := util.NewRandomGenerator(32)

	// Initialize repositories
	authAccountRepo := dbrepository.NewAuthAccountRepository(dbClient)
	localAuthRepo := dbrepository.NewLocalAuthRepository(dbClient)
	oauthAuthRepo := dbrepository.NewOAuthAuthRepository(dbClient)
	tokenRepo := tokenrepo.NewTokenRepository(tokenClient)

	// Initialize code managers
	loginCodeManager := coderepo.NewCodeManager(loginCodeGenerator, cfg.LoginCodeStore.Timeout, loginCodeStore, cfg.LoginCodeStore.Prefix)
	emailCodeManager := coderepo.NewCodeManager(emailCodeGenerator, cfg.EmailCodeStore.Timeout, emailCodeStore, cfg.EmailCodeStore.Prefix)

	// Initialize use cases
	localLoginUsecase := localauth.NewLoginUsecase(authAccountRepo, localAuthRepo, tokenRepo, loginCodeManager)
	localSignupUsecase := localauth.NewSignupUsecase(authAccountRepo, localAuthRepo, tokenRepo, mailer, emailCodeManager, cfg.VerifyEmailURL)
	oauthLoginUsecase := oauthusecase.NewLoginUsecase(authAccountRepo, oauthAuthRepo, tokenRepo, loginCodeManager, oauthApis)

	// Initialize handlers
	localAuthHandler, err := handlerv1.NewLocalAuthHandler(localLoginUsecase, localSignupUsecase, logger, validator)
	if err != nil {
		logger.Fatal("failed to create local auth handler", zap.Error(err))
	}
	oauthHandler, err := handlerv1.NewOAuthHandler(oauthLoginUsecase, logger, validator)
	if err != nil {
		logger.Fatal("failed to create OAuth handler", zap.Error(err))
	}

	// Initialize servers
	httpServer := httpserver.NewServer(cfg.Port, logger, localAuthHandler, oauthHandler, sessionStore)

	manager := server.NewServerManager(
		[]server.Server{httpServer},
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	go func() {
		sig := <-signalChan
		logger.Info("received signal, shutting down", zap.String("signal", sig.String()))
		cancel() // Cancel the context to stop the server
	}()

	if err := manager.Run(ctx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
