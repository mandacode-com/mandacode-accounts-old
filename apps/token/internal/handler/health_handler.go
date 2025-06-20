package handler

import (
	"context"
	"log"
	"time"

	healthv1 "mandacode.com/accounts/token/proto/health/v1"
)

// healthHandler implements the gRPC HealthService interface
type healthHandler struct {
	healthv1.UnimplementedHealthServiceServer
	statusMap map[string]healthv1.ServingStatus
}

// NewHealthHandler returns a new health service handler
func NewHealthHandler() healthv1.HealthServiceServer {
	return &healthHandler{
		statusMap: map[string]healthv1.ServingStatus{
			"":        healthv1.ServingStatus_SERVING_STATUS_SERVING, // default service
			"token":   healthv1.ServingStatus_SERVING_STATUS_SERVING,
			"grpc":    healthv1.ServingStatus_SERVING_STATUS_SERVING,
			"unknown": healthv1.ServingStatus_SERVING_STATUS_NOT_SERVING,
		},
	}
}

// Check handles a unary health check
func (h *healthHandler) Check(ctx context.Context, req *healthv1.CheckRequest) (*healthv1.CheckResponse, error) {
	log.Printf("[HealthCheck] requested for service: %s", req.Service)

	status, ok := h.statusMap[req.Service]
	if !ok {
		status = healthv1.ServingStatus_SERVING_STATUS_UNSPECIFIED
	}

	return &healthv1.CheckResponse{
		Status: status,
	}, nil
}

// Watch handles a streaming health check (e.g. Kubernetes readiness probes)
func (h *healthHandler) Watch(req *healthv1.WatchRequest, stream healthv1.HealthService_WatchServer) error {
	log.Printf("[HealthWatch] watching service: %s", req.Service)

	for {
		status, ok := h.statusMap[req.Service]
		if !ok {
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
