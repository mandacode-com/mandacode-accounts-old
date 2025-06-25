package authhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mandacode.com/accounts/auth/internal/app/auth"
	oauthauthv1 "mandacode.com/accounts/auth/proto/auth/oauth/v1"
)

type OAuthAuthHandler struct {
	oauthauthv1.UnimplementedOAuthAuthServiceServer
	app    *auth.OAuthAuthApp
	logger *zap.Logger
}

// NewOAuthAuthHandler returns a new OAuth authentication service handler
func NewOAuthAuthHandler(app *auth.OAuthAuthApp, logger *zap.Logger) (oauthauthv1.OAuthAuthServiceServer, error) {
	if app == nil {
		return nil, errors.New("OAuthAuthApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &OAuthAuthHandler{
		app:    app,
		logger: logger,
	}, nil
}

func (h *OAuthAuthHandler) Login(ctx context.Context, req *oauthauthv1.OAuthLoginRequest) (*oauthauthv1.OAuthLoginResponse, error) {
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

	return &oauthauthv1.OAuthLoginResponse{
		UserId:       *userID,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
