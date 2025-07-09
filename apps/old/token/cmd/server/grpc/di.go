package grpcserver

import (
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	tokenv1 "github.com/mandacode-com/accounts-proto/token/v1"
	token "mandacode.com/accounts/token/internal/app"
	"mandacode.com/accounts/token/internal/handler"
)

type GRPCRegisterer struct {
	accessTokenApp            *token.AccessTokenApp
	refreshTokenApp           *token.RefreshTokenApp
	emailVerificationTokenApp *token.EmailVerificationTokenApp
	Logger                    *zap.Logger
}

func NewGRPCRegisterer(
	accessTokenApp *token.AccessTokenApp,
	refreshTokenApp *token.RefreshTokenApp,
	emailVerificationTokenApp *token.EmailVerificationTokenApp,
	logger *zap.Logger,
) server.GRPCRegisterer {
	return &GRPCRegisterer{
		accessTokenApp:            accessTokenApp,
		refreshTokenApp:           refreshTokenApp,
		emailVerificationTokenApp: emailVerificationTokenApp,
		Logger:                    logger,
	}
}

func (r *GRPCRegisterer) Register(server *grpc.Server) error {
	tokenHandler, err := handler.NewTokenHandler(
		r.accessTokenApp,
		r.refreshTokenApp,
		r.emailVerificationTokenApp,
		r.Logger,
	)
	if err != nil {
		r.Logger.Error("failed to create token handler", zap.Error(err))
		return err
	}

	tokenv1.RegisterTokenServiceServer(server, tokenHandler)

	return nil
}
