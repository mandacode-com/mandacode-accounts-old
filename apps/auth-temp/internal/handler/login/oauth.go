package loginhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mandacode.com/accounts/auth/internal/app/login"
	oauthloginv1 "mandacode.com/accounts/proto/auth/login/oauth/v1"
)

type OAuthLoginHandler struct {
	oauthloginv1.UnimplementedOAuthLoginServiceServer
	app    *login.OAuthLoginApp
	logger *zap.Logger
}

// NewOAuthLoginHandler returns a new OAuth authentication service handler
func NewOAuthLoginHandler(app *login.OAuthLoginApp, logger *zap.Logger) (oauthloginv1.OAuthLoginServiceServer, error) {
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

	// Attempt to login the user using OAuth
	userID, accessToken, refreshToken, err := h.app.LoginOAuthUser(ctx, req.Provider, req.AccessToken)
	if err != nil {
		h.logger.Error("failed to login OAuth user", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Unauthenticated, "failed to login OAuth user")
	}
	if userID == nil || accessToken == nil || refreshToken == nil {
		h.logger.Error("missing tokens or user ID after OAuth login", zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Internal, "missing tokens or user ID after OAuth login")
	}

	return &oauthloginv1.LoginResponse{
		UserId:       *userID,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
