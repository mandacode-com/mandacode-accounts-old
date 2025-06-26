package loginhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mandacode.com/accounts/auth/internal/app/login"
	localloginv1 "mandacode.com/accounts/proto/auth/login/local/v1"
)

type LocalLoginHandler struct {
	localloginv1.UnimplementedLocalLoginServiceServer
	app    *login.LocalLoginApp
	logger *zap.Logger
}

// NewLocalLoginHandler returns a new local authentication service handler
func NewLocalLoginHandler(app *login.LocalLoginApp, logger *zap.Logger) (localloginv1.LocalLoginServiceServer, error) {
	if app == nil {
		return nil, errors.New("LocalAuthApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &LocalLoginHandler{
		app:    app,
		logger: logger,
	}, nil
}

func (h *LocalLoginHandler) Login(ctx context.Context, req *localloginv1.LoginRequest) (*localloginv1.LoginResponse, error) {
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

	return &localloginv1.LoginResponse{
		UserId:       *userID,
		AccessToken:  *accessToken,
		RefreshToken: *refreshToken,
	}, nil
}
