package userhandler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	oauthuser "mandacode.com/accounts/auth/internal/app/user/oauth"
	protoutil "mandacode.com/accounts/auth/internal/util/proto"
	oauthuserv1 "github.com/mandacode-com/accounts-proto/auth/user/oauth/v1"
)

type OAuthUserHandler struct {
	oauthuserv1.UnimplementedOAuthUserServiceServer
	app    oauthuser.OAuthUserApp
	logger *zap.Logger
}

// NewOAuthUserHandler returns a new OAuth user service handler
func NewOAuthUserHandler(app oauthuser.OAuthUserApp, logger *zap.Logger) (oauthuserv1.OAuthUserServiceServer, error) {
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
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}
	provider, err := protoutil.FromProtoToEntProvider(req.Provider)

	user, err := h.app.CreateUser(userUUID, provider, req.AccessToken)
	if err != nil {
		h.logger.Error("failed to enroll OAuth user", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.Internal, "failed to enroll OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after enrollment", zap.String("user_id", req.UserId), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after enrollment")
	}

	protoOAuthUser, err := protoutil.NewProtoOAuthUser(user)
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
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}
	provider, err := protoutil.FromProtoToEntProvider(req.Provider)

	user, err := h.app.GetUser(userUUID, provider)
	if err != nil {
		h.logger.Error("failed to get OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to get OAuth user")
	}
	if user == nil {
		h.logger.Warn("user not found", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	protoOAuthUser, err := protoutil.NewProtoOAuthUser(user)
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
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	provider, err := protoutil.FromProtoToEntProvider(req.Provider)
	if err != nil {
		h.logger.Error("invalid OAuth provider", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.InvalidArgument, "invalid OAuth provider")
	}

	err = h.app.DeleteUser(userUUID, provider)
	if err != nil {
		h.logger.Error("failed to delete OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to delete OAuth user")
	}

	return &oauthuserv1.DeleteUserResponse{
		UserId:    req.UserId,
		Provider:  req.Provider,
		DeletedAt: timestamppb.Now(),
	}, nil
}

func (h *OAuthUserHandler) DeleteAllProviders(ctx context.Context, req *oauthuserv1.DeleteAllProvidersRequest) (*oauthuserv1.DeleteAllProvidersResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Delete all OAuth providers for the user
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}

	err = h.app.DeleteAllProviders(userUUID)
	if err != nil {
		h.logger.Error("failed to delete all OAuth providers for user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to delete all OAuth providers for user")
	}

	return &oauthuserv1.DeleteAllProvidersResponse{
		UserId: req.UserId,
	}, nil
}

func (h *OAuthUserHandler) SyncUser(ctx context.Context, req *oauthuserv1.SyncUserRequest) (*oauthuserv1.SyncUserResponse, error) {
	// Validate the request parameters
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	// Sync the OAuth user
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}
	provider, err := protoutil.FromProtoToEntProvider(req.Provider)
	if err != nil {
		h.logger.Error("invalid OAuth provider", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.InvalidArgument, "invalid OAuth provider")
	}

	user, err := h.app.SyncUser(userUUID, provider, req.AccessToken)
	if err != nil {
		h.logger.Error("failed to sync OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to sync OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after sync", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after sync")
	}

	protoOAuthUser, err := protoutil.NewProtoOAuthUser(user)
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
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}
	provider, err := protoutil.FromProtoToEntProvider(req.Provider)
	if err != nil {
		h.logger.Error("invalid OAuth provider", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.InvalidArgument, "invalid OAuth provider")
	}

	user, err := h.app.UpdateActiveStatus(userUUID, provider, req.IsActive)
	if err != nil {
		h.logger.Error("failed to update active status of OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update active status of OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after update")
	}

	protoOAuthUser, err := protoutil.NewProtoOAuthUser(user)
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
	userUUID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.InvalidArgument, "invalid user ID format")
	}
	provider, err := protoutil.FromProtoToEntProvider(req.Provider)
	if err != nil {
		h.logger.Error("invalid OAuth provider", zap.Error(err), zap.String("provider", req.Provider.String()))
		return nil, status.Error(codes.InvalidArgument, "invalid OAuth provider")
	}

	user, err := h.app.UpdateVerificationStatus(userUUID, provider, req.IsVerified)
	if err != nil {
		h.logger.Error("failed to update verified status of OAuth user", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Error(codes.Internal, "failed to update verified status of OAuth user")
	}
	if user == nil {
		h.logger.Error("user not found or invalid data after update", zap.String("user_id", req.UserId))
		return nil, status.Error(codes.NotFound, "user not found or invalid data after update")
	}

	protoOAuthUser, err := protoutil.NewProtoOAuthUser(user)
	if err != nil {
		h.logger.Error("failed to build proto OAuth user", zap.Error(err), zap.String("user_id", user.ID.String()))
		return nil, status.Error(codes.Internal, "failed to build proto OAuth user")
	}
	return &oauthuserv1.UpdateVerifiedStatusResponse{
		User: protoOAuthUser,
	}, nil
}
