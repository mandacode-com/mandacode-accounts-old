package authhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mandacode.com/accounts/auth/internal/app/auth"
	localauthv1 "mandacode.com/accounts/auth/proto/auth/local/v1"
)

type LocalAuthHandler struct {
	localauthv1.UnimplementedLocalAuthServiceServer
	app    *auth.LocalAuthApp
	logger *zap.Logger
}

// NewLocalAuthHandler returns a new local authentication service handler
func NewLocalAuthHandler(app *auth.LocalAuthApp, logger *zap.Logger) (localauthv1.LocalAuthServiceServer, error) {
	if app == nil {
		return nil, errors.New("LocalAuthApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &LocalAuthHandler{
		app:    app,
		logger: logger,
	}, nil
}

func (h *LocalAuthHandler) Login(ctx context.Context, req *localauthv1.LocalLoginRequest) (*localauthv1.LocalLoginResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Attempt to login the local user
	userID, accessToken, refreshToken, err := h.app.LoginLocalUser(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("failed to login local user", zap.Error(err), zap.String("email", req.Email))
		return nil, status.Error(codes.Unauthenticated, "failed to login local user")
	}
	if userID == nil || accessToken == nil || refreshToken == nil {
		h.logger.Error("missing tokens or user ID after local login", zap.String("email", req.Email))
		return nil, status.Error(codes.Internal, "missing tokens or user ID after local login")
	}

	return &localauthv1.LocalLoginResponse{
		UserId:       *userID,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
