package loginhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	locallogin "mandacode.com/accounts/auth/internal/app/login/local"
	localloginv1 "mandacode.com/accounts/proto/auth/login/local/v1"
)

type LocalLoginHandler struct {
	localloginv1.UnimplementedLocalLoginServiceServer
	app    locallogin.LocalLoginApp
	logger *zap.Logger
}

// NewLocalLoginHandler returns a new local authentication service handler
func NewLocalLoginHandler(app locallogin.LocalLoginApp, logger *zap.Logger) (localloginv1.LocalLoginServiceServer, error) {
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
	loginToken, err := h.app.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.logger.Error("failed to login local user", zap.Error(err), zap.String("email", req.Email))
		return nil, status.Error(codes.Unauthenticated, "failed to login local user")
	}

	return &localloginv1.LoginResponse{
		AccessToken:  loginToken.AccessToken,
		RefreshToken: loginToken.RefreshToken,
	}, nil
}
