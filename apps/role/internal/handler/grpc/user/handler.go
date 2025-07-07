package usergrpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	userv1 "github.com/mandacode-com/accounts-proto/role/user/v1"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
)

type UserGRPCHandler struct {
	userv1.UnimplementedUserServiceServer
	groupUserApp groupuserapp.GroupUserApp
	logger       *zap.Logger
}

// CleanupRole implements userv1.UserServiceServer.
func (u *UserGRPCHandler) CleanupRole(ctx context.Context, req *userv1.CleanupRoleRequest) (*userv1.CleanupRoleResponse, error) {
	if err := req.ValidateAll(); err != nil {
		u.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		u.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	if err := u.groupUserApp.DeleteGroupUserByUserID(userID); err != nil {
		u.logger.Error("failed to delete group user by user ID", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.Internal, "failed to delete group user by user ID: %v", err)
	}

	return &userv1.CleanupRoleResponse{
		UserId:      req.UserId,
		CleanupTime: timestamppb.Now(),
	}, nil
}

func NewUserGRPCHandler(groupUserApp groupuserapp.GroupUserApp, logger *zap.Logger) (userv1.UserServiceServer, error) {
	if groupUserApp == nil {
		return nil, errors.New("groupUserApp cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &UserGRPCHandler{
		groupUserApp: groupUserApp,
		logger:       logger,
	}, nil
}
