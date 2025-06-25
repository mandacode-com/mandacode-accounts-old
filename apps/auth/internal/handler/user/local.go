package userhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mandacode.com/accounts/auth/internal/app/user"
	localuserv1 "mandacode.com/accounts/auth/proto/user/local/v1"
)

type LocalUserHandler struct {
	localuserv1.UnimplementedLocalUserServiceServer
	localUserApp *user.LocalUserApp
	logger       *zap.Logger
}

// NewLocalUserHandler returns a new local user service handler
func NewLocalUserHandler(localUserApp *user.LocalUserApp, logger *zap.Logger) (localuserv1.LocalUserServiceServer, error) {
	if localUserApp == nil {
		return nil, errors.New("localUserApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	return &LocalUserHandler{
		localUserApp: localUserApp,
		logger:       logger,
	}, nil
}

func (h *LocalUserHandler) EnrollUser(ctx context.Context, req *localuserv1.EnrollUserRequest) (*localuserv1.EnrollUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Create the local user
	userID, err := h.localUserApp.CreateLocalUser(ctx, req.UserId, req.Email, req.Password, req.IsActive, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to enroll local user", zap.Error(err), zap.String("email", req.Email))
		return nil, status.Error(codes.Internal, "failed to enroll local user")
	}
	if userID == nil {
		h.logger.Error("missing user ID after local enrollment", zap.String("email", req.Email))
		return nil, status.Error(codes.Internal, "missing user ID after local enrollment")
	}

	return &localuserv1.EnrollUserResponse{
		Success: true,
	}, nil
}

func (h *LocalUserHandler) DeleteUser(ctx context.Context, req *localuserv1.DeleteUserRequest) (*localuserv1.DeleteUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Delete the local user
	err := h.localUserApp.DeleteLocalUser(ctx, req.UserId)
	if err != nil {
		h.logger.Error("failed to delete local user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to delete local user")
	}

	return &localuserv1.DeleteUserResponse{
		Success: true,
	}, nil
}

func (h *LocalUserHandler) UpdateUser(ctx context.Context, req *localuserv1.UpdateUserRequest) (*localuserv1.UpdateUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the local user
	user, err := h.localUserApp.UpdateLocalUser(ctx, req.UserId, req.Email, req.Password, req.IsActive, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to update local user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update local user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after update")
	}

	return &localuserv1.UpdateUserResponse{
		Success: true,
	}, nil
}
