package userhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mandacode.com/accounts/auth/internal/app/user"
	oauthuserv1 "mandacode.com/accounts/auth/proto/user/oauth/v1"
)

type OAuthUserHandler struct {
	oauthuserv1.UnimplementedOAuthUserServiceServer
	app    *user.OAuthUserApp
	logger *zap.Logger
}

// NewOAuthUserHandler returns a new OAuth user service handler
func NewOAuthUserHandler(app *user.OAuthUserApp, logger *zap.Logger) (oauthuserv1.OAuthUserServiceServer, error) {
	if app == nil {
		return nil, errors.New("OAuthUserApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &OAuthUserHandler{
		app:    app,
		logger: logger,
	}, nil
}

func (h *OAuthUserHandler) EnrollUser(ctx context.Context, req *oauthuserv1.EnrollUserRequest) (*oauthuserv1.EnrollUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Create the OAuth user
	userID, err := h.app.CreateOAuthUser(ctx, req.UserId, req.Provider, req.AccessToken, req.IsActive, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to enroll OAuth user", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Internal, "failed to enroll OAuth user")
	}
	if userID == nil {
		h.logger.Error("missing user ID after OAuth enrollment", zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Internal, "missing user ID after OAuth enrollment")
	}

	return &oauthuserv1.EnrollUserResponse{
		Success: true,
	}, nil
}

func (h *OAuthUserHandler) DeleteUser(ctx context.Context, req *oauthuserv1.DeleteUserRequest) (*oauthuserv1.DeleteUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Delete the OAuth user
	err := h.app.DeleteOAuthUser(ctx, req.UserId, req.Provider)
	if err != nil {
		h.logger.Error("failed to delete OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to delete OAuth user")
	}

	return &oauthuserv1.DeleteUserResponse{
		Success: true,
	}, nil
}

func (h *OAuthUserHandler) UpdateUser(ctx context.Context, req *oauthuserv1.UpdateUserRequest) (*oauthuserv1.UpdateUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the OAuth user
	user, err := h.app.UpdateOAuthUser(ctx, req.UserId, req.Provider, req.ProviderId, req.Email, req.IsActive, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to update OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after update")
	}

	return &oauthuserv1.UpdateUserResponse{
		Success: true,
	}, nil
}
