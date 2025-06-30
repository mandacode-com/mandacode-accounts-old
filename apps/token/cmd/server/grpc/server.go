package grpcserver

import (
	"net"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"mandacode.com/lib/server/server"
)

type GRPCServer struct {
	server     *grpc.Server
	registerer server.GRPCRegisterer
	logger     *zap.Logger
	port       int
}

func NewGRPCServer(port int, logger *zap.Logger, registerer server.GRPCRegisterer, servingServices []string) (server.Server, error) {
	server := grpc.NewServer()

	// Register health check service
	healthServer := health.NewServer()
	for _, service := range servingServices {
		healthServer.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)
	}
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Register reflection service on gRPC server
	reflection.Register(server)

	return &GRPCServer{
		server:     server,
		registerer: registerer,
		logger:     logger,
		port:       port,
	}, nil
}

func (g *GRPCServer) Start() error {
	if err := g.registerer.Register(g.server); err != nil {
		g.logger.Error("failed to register gRPC handlers", zap.Error(err))
		return err
	}

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(g.port))
	if err != nil {
		g.logger.Error("failed to listen on port", zap.Int("port", g.port), zap.Error(err))
		return err
	}

	g.logger.Info("gRPC server is running", zap.String("address", ":"+strconv.Itoa(g.port)))
	return g.server.Serve(lis)
}

func (g *GRPCServer) Stop() error {
	g.logger.Info("stopping gRPC server")
	g.server.GracefulStop()
	g.logger.Info("gRPC server stopped")
	return nil
}
