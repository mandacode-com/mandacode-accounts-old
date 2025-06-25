package main

import (
	"net"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	port   int
	logger *zap.Logger
}

func newGRPCServer(logger *zap.Logger) (*GRPCServer, error) {
	cfg := loadConfig(logger)

	server := grpc.NewServer()
	if err := registerHandlers(server, cfg, logger); err != nil {
		return nil, err
	}

	// Register health check service
	healthServer := health.NewServer()
	healthServer.SetServingStatus("token.v1.TokenService", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service on gRPC server
	reflection.Register(server)

	return &GRPCServer{
		server: server,
		port:   cfg.Port,
		logger: logger,
	}, nil
}

func (g *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(g.port))
	if err != nil {
		return err
	}
	g.logger.Info("gRPC server is running", zap.String("address", ":"+strconv.Itoa(g.port)))
	return g.server.Serve(lis)
}
