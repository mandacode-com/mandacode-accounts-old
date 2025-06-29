package grpchandler

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mandacode.com/accounts/profile/internal/app/profile"
	profilev1 "mandacode.com/accounts/proto/profile/v1"
)

type ProfileGRPCHandler struct {
	profilev1.UnimplementedProfileServiceServer
	initializeUC profile.InitializeProfileUsecase
	deleteUC     profile.DeleteProfileUsecase
	logger       *zap.Logger
}

func NewProfileGRPCHandler(initializeUC profile.InitializeProfileUsecase, deleteUC profile.DeleteProfileUsecase, logger *zap.Logger) (profilev1.ProfileServiceServer, error) {
	if initializeUC == nil {
		return nil, errors.New("initializeUC cannot be nil")
	}
	if deleteUC == nil {
		return nil, errors.New("deleteUC cannot be nil")
	}
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &ProfileGRPCHandler{
		initializeUC: initializeUC,
		deleteUC:     deleteUC,
		logger:       logger,
	}, nil
}

func (h *ProfileGRPCHandler) CreateProfile(context context.Context, req *profilev1.CreateProfileRequest) (*profilev1.CreateProfileResponse, error) {
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	profile, err := h.initializeUC.InitializeProfile(userID)
	if err != nil {
		h.logger.Error("failed to initialize profile", zap.Error(err))
		return nil, err
	}
	protoProfile, err := ToProtoProfile(profile)
	if err != nil {
		h.logger.Error("failed to convert profile to proto", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to convert profile to proto: %v", err)
	}
	if protoProfile == nil {
		h.logger.Error("converted profile is nil")
		return nil, status.Error(codes.Internal, "converted profile is nil")
	}

	return &profilev1.CreateProfileResponse{
		Profile: protoProfile,
	}, nil
}

func (h *ProfileGRPCHandler) DeleteProfile(context context.Context, req *profilev1.DeleteProfileRequest) (*profilev1.DeleteProfileResponse, error) {
	if err := req.ValidateAll(); err != nil {
		h.logger.Error("invalid request parameters", zap.Error(err))
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters: %v", err)
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		h.logger.Error("invalid user ID format", zap.Error(err), zap.String("user_id", req.UserId))
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	err = h.deleteUC.DeleteProfile(userID)
	if err != nil {
		h.logger.Error("failed to delete profile", zap.Error(err))
		return nil, err
	}

	return &profilev1.DeleteProfileResponse{
		UserId:    userID.String(),
		DeletedAt: timestamppb.Now(),
	}, nil
}
