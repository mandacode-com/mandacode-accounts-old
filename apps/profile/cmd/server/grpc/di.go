package grpcserver

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	grpcuc "mandacode.com/accounts/profile/internal/app/profile/grpc"
	svcdomain "mandacode.com/accounts/profile/internal/domain/service"
	grpchandler "mandacode.com/accounts/profile/internal/handler/grpc"
	profilev1 "mandacode.com/accounts/proto/profile/v1"
	"mandacode.com/lib/server/server"
)

type GRPCRegisterer struct {
	ProfileService svcdomain.ProfileService
	Logger         *zap.Logger
}

func NewGRPCRegisterer(profileService svcdomain.ProfileService, logger *zap.Logger) server.GRPCRegisterer {
	return &GRPCRegisterer{
		ProfileService: profileService,
		Logger:         logger,
	}
}

func (r *GRPCRegisterer) Register(server *grpc.Server) error {
	initializeProfileUC := grpcuc.NewInitializeProfileUsecase(r.ProfileService)
	deleteProfileUC := grpcuc.NewDeleteProfileUsecase(r.ProfileService)

	profileHandler, err := grpchandler.NewProfileGRPCHandler(initializeProfileUC, deleteProfileUC, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create gRPC handler", zap.Error(err))
		return err
	}

	profilev1.RegisterProfileServiceServer(server, profileHandler)

	return nil
}
