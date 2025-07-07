package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	httpserver "mandacode.com/accounts/user/cmd/server/http"
	localuserapp "mandacode.com/accounts/user/internal/app/user/local"
	oauthuserapp "mandacode.com/accounts/user/internal/app/user/oauth"
	userapp "mandacode.com/accounts/user/internal/app/user/user"
	"mandacode.com/accounts/user/internal/config"
	localuser "mandacode.com/accounts/user/internal/infra/auth/local"
	oauthuser "mandacode.com/accounts/user/internal/infra/auth/oauth"
	"mandacode.com/accounts/user/internal/infra/database"
	"mandacode.com/accounts/user/internal/infra/database/repository"
	"mandacode.com/accounts/user/internal/infra/mailer"
	"mandacode.com/accounts/user/internal/infra/profile"
	token "mandacode.com/accounts/user/internal/infra/token"
	userevent "mandacode.com/accounts/user/internal/infra/user_event"
	"mandacode.com/accounts/user/internal/util"
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

	entClient, err := database.NewEntClient(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to create database client", zap.Error(err))
	}
	repo := repository.NewUserRepository(entClient)

	// Kafka writers for user events
	userCreationFailedWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.ServiceAddr.Kafka),
		Topic:    "user_creation_failed",
		Balancer: &kafka.Hash{},
	}
	archiveUserWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.ServiceAddr.Kafka),
		Topic:    "archive_user",
		Balancer: &kafka.Hash{},
	}
	deleteUserWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.ServiceAddr.Kafka),
		Topic:    "delete_user",
		Balancer: &kafka.Hash{},
	}
	mailWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.ServiceAddr.Kafka),
		Topic:    "mail",
		Balancer: &kafka.Hash{},
	}

	localUserClient, _, err := localuser.NewGRPCClient(cfg.ServiceAddr.Auth)
	if err != nil {
		logger.Fatal("failed to create local user gRPC client", zap.Error(err))
	}
	localUserService := localuser.NewLocalUserService(localUserClient)

	oauthUserClient, _, err := oauthuser.NewGRPCClient(cfg.ServiceAddr.Auth)
	if err != nil {
		logger.Fatal("failed to create OAuth user gRPC client", zap.Error(err))
	}
	oauthUserService := oauthuser.NewOAuthUserService(oauthUserClient)

	profileClient, _, err := profile.NewGRPCClient(cfg.ServiceAddr.Profile)
	if err != nil {
		logger.Fatal("failed to create profile gRPC client", zap.Error(err))
	}
	profileService := profile.NewProfileService(profileClient)

	tokenClient, _, err := token.NewGRPCClient(cfg.ServiceAddr.Token)
	if err != nil {
		logger.Fatal("failed to create token gRPC client", zap.Error(err))
	}
	tokenService := token.NewTokenService(tokenClient)

	mailerService := mailer.NewMailer(mailWriter)

	codeGenerator := util.NewRandomCodeGenerator(cfg.SyncCodeLength, cfg.EmailVerificationCodeLength)

	userEventService := userevent.NewUserEventService(userCreationFailedWriter, archiveUserWriter, deleteUserWriter)

	userApp := userapp.NewUserApp(repo, localUserService, oauthUserService, profileService, userEventService, codeGenerator)
	localUserApp := localuserapp.NewLocalUserApp(repo, localUserService, profileService, tokenService, mailerService, userEventService)
	oauthUserApp := oauthuserapp.NewOAuthUserApp(repo, oauthUserService, profileService, userEventService, codeGenerator)

	httpServerRegisterer := httpserver.NewHTTPRegisterer(
		userApp,
		localUserApp,
		oauthUserApp,
		logger,
		validator,
		cfg.UIDHeader,
	)

	httpServer, err := httpserver.NewHTTPServer(cfg.Port, logger, httpServerRegisterer)
	if err != nil {
		logger.Fatal("failed to create HTTP server", zap.Error(err))
	}
	if err := httpServer.Start(); err != nil {
		logger.Fatal("failed to start HTTP server", zap.Error(err))
	}
}
