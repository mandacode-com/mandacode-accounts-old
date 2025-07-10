package grpcserver

import (
	localloginv1 "github.com/mandacode-com/accounts-proto/auth/login/local/v1"
	oauthloginv1 "github.com/mandacode-com/accounts-proto/auth/login/oauth/v1"
	tokenv1 "github.com/mandacode-com/accounts-proto/auth/token/v1"
	localuserv1 "github.com/mandacode-com/accounts-proto/auth/user/local/v1"
	oauthuserv1 "github.com/mandacode-com/accounts-proto/auth/user/oauth/v1"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	locallogin "mandacode.com/accounts/auth/internal/app/login/local"
	oauthlogin "mandacode.com/accounts/auth/internal/app/login/oauth"
	"mandacode.com/accounts/auth/internal/app/token"
	localuser "mandacode.com/accounts/auth/internal/app/user/local"
	oauthuser "mandacode.com/accounts/auth/internal/app/user/oauth"
	loginhandler "mandacode.com/accounts/auth/internal/handler/login"
	tokenhandler "mandacode.com/accounts/auth/internal/handler/token"
	userhandler "mandacode.com/accounts/auth/internal/handler/user"
)

type GRPCRegisterer struct {
	LocalLoginApp locallogin.LocalLoginApp
	OAuthLoginApp oauthlogin.OAuthLoginApp
	LocalUserApp  localuser.LocalUserApp
	OAuthUserApp  oauthuser.OAuthUserApp
	TokenApp      token.TokenApp
	Logger        *zap.Logger
}

func NewGRPCRegisterer(
	localLoginApp locallogin.LocalLoginApp,
	oauthLoginApp oauthlogin.OAuthLoginApp,
	localUserApp localuser.LocalUserApp,
	oauthUserApp oauthuser.OAuthUserApp,
	tokenApp token.TokenApp,
	logger *zap.Logger,
) server.GRPCRegisterer {
	return &GRPCRegisterer{
		LocalLoginApp: localLoginApp,
		OAuthLoginApp: oauthLoginApp,
		LocalUserApp:  localUserApp,
		OAuthUserApp:  oauthUserApp,
		TokenApp:      tokenApp,
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
	tokenHandler, err := tokenhandler.NewTokenHandler(r.TokenApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create token handler", zap.Error(err))
		return err
	}

	localloginv1.RegisterLocalLoginServiceServer(server, localLoginHandler)
	oauthloginv1.RegisterOAuthLoginServiceServer(server, oauthLoginHandler)
	localuserv1.RegisterLocalUserServiceServer(server, localUserHandler)
	oauthuserv1.RegisterOAuthUserServiceServer(server, oauthUserHandler)
	tokenv1.RegisterTokenServiceServer(server, tokenHandler)

	return nil
}
