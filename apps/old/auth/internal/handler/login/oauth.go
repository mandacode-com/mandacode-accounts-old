package loginhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	oauthlogin "mandacode.com/accounts/auth/internal/app/login/oauth"
	protoutil "mandacode.com/accounts/auth/internal/util/proto"
	oauthloginv1 "github.com/mandacode-com/accounts-proto/auth/login/oauth/v1"
)

type OAuthLoginHandler struct {
	oauthloginv1.UnimplementedOAuthLoginServiceServer
	app    oauthlogin.OAuthLoginApp
	logger *zap.Logger
}

// NewOAuthLoginHandler returns a new OAuth authentication service handler
func NewOAuthLoginHandler(app oauthlogin.OAuthLoginApp, logger *zap.Logger) (oauthloginv1.OAuthLoginServiceServer, error) {
	if app == nil {
		return nil, errors.New("OAuthAuthApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &OAuthLoginHandler{
		app:    app,
		logger: logger,
	}, nil
}

func (h *OAuthLoginHandler) Login(ctx context.Context, req *oauthloginv1.LoginRequest) (*oauthloginv1.LoginResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	provider, err := protoutil.FromProtoToEntProvider(req.Provider)
	if err != nil {
		h.logger.Error("invalid OAuth provider", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.InvalidArgument, "invalid OAuth provider")
	}
	// Attempt to login the user using OAuth
	loginToken, err := h.app.Login(ctx, provider, req.AccessToken)
	if err != nil {
		h.logger.Error("failed to login OAuth user", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Unauthenticated, "failed to login OAuth user")
	}
	return &oauthloginv1.LoginResponse{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}, nil
}
