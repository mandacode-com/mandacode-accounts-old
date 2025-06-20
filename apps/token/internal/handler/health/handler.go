package healthhandler

import (
	"context"
	"time"

	"go.uber.org/zap"
	healthv1 "mandacode.com/accounts/token/proto/health/v1"
)

// healthHandler implements the gRPC HealthService interface
type healthHandler struct {
	healthv1.UnimplementedHealthServiceServer
	statusMap map[string]healthv1.ServingStatus
	logger    *zap.Logger
}

// NewHealthHandler returns a new health service handler
func NewHealthHandler(logger *zap.Logger) healthv1.HealthServiceServer {
	return &healthHandler{
		logger: logger,
		statusMap: map[string]healthv1.ServingStatus{
			"token":  healthv1.ServingStatus_SERVING_STATUS_SERVING,
			"health": healthv1.ServingStatus_SERVING_STATUS_SERVING,
		},
	}
}

// Check handles a unary health check
func (h *healthHandler) Check(ctx context.Context, req *healthv1.CheckRequest) (*healthv1.CheckResponse, error) {
	status, ok := h.statusMap[req.Service]
	if !ok {
		h.logger.Error("unknown service in health check", zap.String("service", req.Service))
		status = healthv1.ServingStatus_SERVING_STATUS_UNSPECIFIED
	}

	return &healthv1.CheckResponse{
		Status: status,
	}, nil
}

// Watch handles a streaming health check (e.g. Kubernetes readiness probes)
func (h *healthHandler) Watch(req *healthv1.WatchRequest, stream healthv1.HealthService_WatchServer) error {
	for {
		status, ok := h.statusMap[req.Service]
		if !ok {
			h.logger.Error("unknown service in health watch", zap.String("service", req.Service))
			status = healthv1.ServingStatus_SERVING_STATUS_UNSPECIFIED
		}

		resp := &healthv1.WatchResponse{
			Status:    status,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(resp); err != nil {
			return err
		}

		time.Sleep(3 * time.Second)
	}
}
