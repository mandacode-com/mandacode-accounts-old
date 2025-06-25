package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := NewLogger()
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	server, err := newGRPCServer(logger)
	if err != nil {
		logger.Fatal("failed to initialize gRPC server", zap.Error(err))
	}

	if err := server.Start(); err != nil {
		logger.Fatal("gRPC server stopped unexpectedly", zap.Error(err))
	}
}
