package grpcserver

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	locallogin "mandacode.com/accounts/auth/internal/app/login/local"
	oauthlogin "mandacode.com/accounts/auth/internal/app/login/oauth"
	localuser "mandacode.com/accounts/auth/internal/app/user/local"
	oauthuser "mandacode.com/accounts/auth/internal/app/user/oauth"
	loginhandler "mandacode.com/accounts/auth/internal/handler/login"
	userhandler "mandacode.com/accounts/auth/internal/handler/user"
	localloginv1 "mandacode.com/accounts/proto/auth/login/local/v1"
	oauthloginv1 "mandacode.com/accounts/proto/auth/login/oauth/v1"
	localuserv1 "mandacode.com/accounts/proto/auth/user/local/v1"
	oauthuserv1 "mandacode.com/accounts/proto/auth/user/oauth/v1"
	"mandacode.com/lib/server/server"
)

type GRPCRegisterer struct {
	LocalLoginApp locallogin.LocalLoginApp
	OAuthLoginApp oauthlogin.OAuthLoginApp
	LocalUserApp  localuser.LocalUserApp
	OAuthUserApp  oauthuser.OAuthUserApp
	Logger        *zap.Logger
}

func NewGRPCRegisterer(
	localLoginApp locallogin.LocalLoginApp,
	oauthLoginApp oauthlogin.OAuthLoginApp,
	localUserApp localuser.LocalUserApp,
	oauthUserApp oauthuser.OAuthUserApp,
	logger *zap.Logger,
) server.GRPCRegisterer {
	return &GRPCRegisterer{
		LocalLoginApp: localLoginApp,
		OAuthLoginApp: oauthLoginApp,
		LocalUserApp:  localUserApp,
		OAuthUserApp:  oauthUserApp,
		Logger:        logger,
	}
}

func (r *GRPCRegisterer) Register(server *grpc.Server) error {
	localLoginHandler, err := loginhandler.NewLocalLoginHandler(r.LocalLoginApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create local login handler", zap.Error(err))
		return err
	}
	oauthLoginHandler, err := loginhandler.NewOAuthLoginHandler(r.OAuthLoginApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create OAuth login handler", zap.Error(err))
		return err
	}
	localUserHandler, err := userhandler.NewLocalUserHandler(r.LocalUserApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create local user handler", zap.Error(err))
		return err
	}
	oauthUserHandler, err := userhandler.NewOAuthUserHandler(r.OAuthUserApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create OAuth user handler", zap.Error(err))
		return err
	}

	localloginv1.RegisterLocalLoginServiceServer(server, localLoginHandler)
	oauthloginv1.RegisterOAuthLoginServiceServer(server, oauthLoginHandler)
	localuserv1.RegisterLocalUserServiceServer(server, localUserHandler)
	oauthuserv1.RegisterOAuthUserServiceServer(server, oauthUserHandler)

	return nil
}
