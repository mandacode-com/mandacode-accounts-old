package grpcserver

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mandacode.com/accounts/profile/internal/app/profile"
	grpchandler "mandacode.com/accounts/profile/internal/handler/grpc"
	profilev1 "mandacode.com/accounts/proto/profile/v1"
	"mandacode.com/lib/server/server"
)

type GRPCRegisterer struct {
	ProfileApp profile.ProfileApp
	Logger         *zap.Logger
}

func NewGRPCRegisterer(profileApp profile.ProfileApp, logger *zap.Logger) server.GRPCRegisterer {
	return &GRPCRegisterer{
		ProfileApp: profileApp,
		Logger:     logger,
	}
}

func (r *GRPCRegisterer) Register(server *grpc.Server) error {
	profileHandler, err := grpchandler.NewProfileGRPCHandler(r.ProfileApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create gRPC handler", zap.Error(err))
		return err
	}

	profilev1.RegisterProfileServiceServer(server, profileHandler)

	return nil
}
