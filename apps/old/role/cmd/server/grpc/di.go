package grpcserver

import (
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	userv1 "github.com/mandacode-com/accounts-proto/role/user/v1"
	groupuserapp "mandacode.com/accounts/role/internal/app/groupuser"
	usergrpc "mandacode.com/accounts/role/internal/handler/grpc/user"
)

type GRPCRegisterer struct {
	GroupUserApp groupuserapp.GroupUserApp
	Logger       *zap.Logger
}

func NewGRPCRegisterer(groupUserApp groupuserapp.GroupUserApp, logger *zap.Logger) server.GRPCRegisterer {
	return &GRPCRegisterer{
		GroupUserApp: groupUserApp,
		Logger:       logger,
	}
}

func (r *GRPCRegisterer) Register(server *grpc.Server) error {
	userGRPCHandler, err := usergrpc.NewUserGRPCHandler(r.GroupUserApp, r.Logger)
	if err != nil {
		r.Logger.Error("failed to create gRPC handler", zap.Error(err))
		return err
	}

	userv1.RegisterUserServiceServer(server, userGRPCHandler)

	return nil
}
