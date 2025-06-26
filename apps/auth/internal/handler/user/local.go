package userhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/auth/internal/app/user"
	localuserv1 "mandacode.com/accounts/proto/auth/user/local/v1"
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

func (h *LocalUserHandler) GetUser(ctx context.Context, req *localuserv1.GetUserRequest) (*localuserv1.GetUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Retrieve the local user
	user, err := h.localUserApp.GetUser(ctx, req.UserId)
	if err != nil {
		h.logger.Error("failed to get local user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to get local user")
	}
	if user == nil {
		h.logger.Warn("user not found", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &localuserv1.GetUserResponse{
		UserId:     user.ID.String(),
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *LocalUserHandler) EnrollUser(ctx context.Context, req *localuserv1.EnrollUserRequest) (*localuserv1.EnrollUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Create the local user
	user, err := h.localUserApp.CreateUser(ctx, req.UserId, req.Email, req.Password, req.IsActive, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to enroll local user", zap.Error(err), zap.String("email", req.Email))
		return nil, status.Error(codes.Internal, "failed to enroll local user")
	}
	if user == nil {
		h.logger.Error("missing user ID after local enrollment", zap.String("email", req.Email))
		return nil, status.Error(codes.Internal, "missing user ID after local enrollment")
	}

	return &localuserv1.EnrollUserResponse{
		UserId:     user.ID.String(),
		Email:      user.Email,
		IsActive:   user.IsActive,
		IsVerified: user.IsVerified,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *LocalUserHandler) DeleteUser(ctx context.Context, req *localuserv1.DeleteUserRequest) (*localuserv1.DeleteUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Delete the local user
	user, err := h.localUserApp.DeleteUser(ctx, req.UserId)
	if err != nil {
		h.logger.Error("failed to delete local user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to delete local user")
	}

	return &localuserv1.DeleteUserResponse{
		UserId:    user.ID.String(),
		DeletedAt: timestamppb.New(user.DeletedAt),
	}, nil
}

func (h *LocalUserHandler) UpdateEmail(ctx context.Context, req *localuserv1.UpdateEmailRequest) (*localuserv1.UpdateEmailResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	user, err := h.localUserApp.UpdateEmail(ctx, req.UserId, req.Email)
	if err != nil {
		h.logger.Error("failed to update local user email", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update local user email")
	}
	if user == nil {
		h.logger.Error("missing user ID after local email update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "missing user ID after local email update")
	}

	return &localuserv1.UpdateEmailResponse{
		UserId:    user.ID.String(),
		Email:     user.Email,
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *LocalUserHandler) UpdatePassword(ctx context.Context, req *localuserv1.UpdatePasswordRequest) (*localuserv1.UpdatePasswordResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the local user's password
	user, err := h.localUserApp.UpdatePassword(ctx, req.UserId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		h.logger.Error("failed to update local user password", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update local user password")
	}
	if user == nil {
		h.logger.Error("missing user ID after local password update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "missing user ID after local password update")
	}

	return &localuserv1.UpdatePasswordResponse{
		UserId:    user.ID.String(),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *LocalUserHandler) UpdateActiveStatus(ctx context.Context, req *localuserv1.UpdateActiveStatusRequest) (*localuserv1.UpdateActiveStatusResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the local user's active status
	user, err := h.localUserApp.UpdateActiveStatus(ctx, req.UserId, req.IsActive)
	if err != nil {
		h.logger.Error("failed to update local user active status", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update local user active status")
	}
	if user == nil {
		h.logger.Error("missing user ID after local active status update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "missing user ID after local active status update")
	}

	return &localuserv1.UpdateActiveStatusResponse{
		UserId:    user.ID.String(),
		IsActive:  user.IsActive,
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, nil
}

func (h *LocalUserHandler) UpdateVerifiedStatus(ctx context.Context, req *localuserv1.UpdateVerifiedStatusRequest) (*localuserv1.UpdateVerifiedStatusResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the local user's verified status
	user, err := h.localUserApp.UpdateVerifiedStatus(ctx, req.UserId, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to update local user verified status", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update local user verified status")
	}
	if user == nil {
		h.logger.Error("missing user ID after local verified status update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "missing user ID after local verified status update")
	}

	return &localuserv1.UpdateVerifiedStatusResponse{
		UserId:     user.ID.String(),
		IsVerified: user.IsVerified,
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}, nil
}
