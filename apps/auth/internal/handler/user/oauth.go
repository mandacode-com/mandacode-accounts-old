package userhandler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/auth/internal/app/user"
	"mandacode.com/accounts/auth/internal/util"
	oauthuserv1 "mandacode.com/accounts/proto/auth/user/oauth/v1"
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
	user, err := h.app.CreateUser(ctx, req.UserId, req.Provider, req.AccessToken, req.IsActive, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to enroll OAuth user", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Internal, "failed to enroll OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after enrollment", zap.String("user_id", req.UserId), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after enrollment")
	}

	protoOAuthUser, err := util.BuildProtoOAuthUser(user)
	if err != nil {
		h.logger.Error("failed to build proto OAuth user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return nil, status.Error(codes.Internal, "failed to build proto OAuth user")
	}

	return &oauthuserv1.EnrollUserResponse{
		User: protoOAuthUser,
	}, nil
}

func (h *OAuthUserHandler) GetUser(ctx context.Context, req *oauthuserv1.GetUserRequest) (*oauthuserv1.GetUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Retrieve the OAuth user
	user, err := h.app.GetUser(ctx, req.UserId, req.Provider)
	if err != nil {
		h.logger.Error("failed to get OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to get OAuth user")
	}
	if user == nil {
		h.logger.Warn("user not found", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	protoOAuthUser, err := util.BuildProtoOAuthUser(user)
	if err != nil {
		h.logger.Error("failed to build proto OAuth user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return nil, status.Error(codes.Internal, "failed to build proto OAuth user")
	}
	return &oauthuserv1.GetUserResponse{
		User: protoOAuthUser,
	}, nil
}

func (h *OAuthUserHandler) DeleteUser(ctx context.Context, req *oauthuserv1.DeleteUserRequest) (*oauthuserv1.DeleteUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Delete the OAuth user
	user, err := h.app.DeleteUser(ctx, req.UserId, req.Provider)
	if err != nil {
		h.logger.Error("failed to delete OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to delete OAuth user")
	}

	return &oauthuserv1.DeleteUserResponse{
		UserId:    user.ID.String(),
		Provider:  *req.Provider,
		DeletedAt: timestamppb.New(user.DeletedAt),
	}, nil
}

func (h *OAuthUserHandler) SyncUser(ctx context.Context, req *oauthuserv1.SyncUserRequest) (*oauthuserv1.SyncUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Sync the OAuth user
	user, err := h.app.SyncUser(ctx, req.UserId, req.Provider, req.AccessToken)
	if err != nil {
		h.logger.Error("failed to sync OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to sync OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after sync", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after sync")
	}

	protoOAuthUser, err := util.BuildProtoOAuthUser(user)
	if err != nil {
		h.logger.Error("failed to build proto OAuth user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return nil, status.Error(codes.Internal, "failed to build proto OAuth user")
	}
	return &oauthuserv1.SyncUserResponse{
		User: protoOAuthUser,
	}, nil
}

// Update Active Status
func (h *OAuthUserHandler) UpdateActiveStatus(ctx context.Context, req *oauthuserv1.UpdateActiveStatusRequest) (*oauthuserv1.UpdateActiveStatusResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the active status of the OAuth user
	user, err := h.app.UpdateActiveStatus(ctx, req.UserId, req.Provider, req.IsActive)
	if err != nil {
		h.logger.Error("failed to update active status of OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update active status of OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after update")
	}

	protoOAuthUser, err := util.BuildProtoOAuthUser(user)
	if err != nil {
		h.logger.Error("failed to build proto OAuth user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return nil, status.Error(codes.Internal, "failed to build proto OAuth user")
	}
	return &oauthuserv1.UpdateActiveStatusResponse{
		User: protoOAuthUser,
	}, nil
}

// Update Verified Status
func (h *OAuthUserHandler) UpdateVerifiedStatus(ctx context.Context, req *oauthuserv1.UpdateVerifiedStatusRequest) (*oauthuserv1.UpdateVerifiedStatusResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Update the verified status of the OAuth user
	user, err := h.app.UpdateVerifiedStatus(ctx, req.UserId, req.Provider, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to update verified status of OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update verified status of OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after update")
	}

	protoOAuthUser, err := util.BuildProtoOAuthUser(user)
	if err != nil {
		h.logger.Error("failed to build proto OAuth user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return nil, status.Error(codes.Internal, "failed to build proto OAuth user")
	}
	return &oauthuserv1.UpdateVerifiedStatusResponse{
		User: protoOAuthUser,
	}, nil
}
