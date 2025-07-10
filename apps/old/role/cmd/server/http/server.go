package httpserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/golib/server"
	"go.uber.org/zap"
)

type HTTPServer struct {
	http      *http.Server
	engine    *gin.Engine
	clientReg server.HTTPRegisterer
	adminReg  server.HTTPRegisterer
	logger    *zap.Logger
	port      int
}

func NewHTTPServer(
	port int,
	logger *zap.Logger,
	clientReg server.HTTPRegisterer,
	adminReg server.HTTPRegisterer,
) (server.Server, error) {
	engine := gin.Default()
	return &HTTPServer{
		http:   &http.Server{Addr: ":" + strconv.Itoa(port), Handler: engine},
		engine: engine,
		clientReg: clientReg,
		adminReg:  adminReg,
		logger: logger,
		port:   port,
	}, nil
}

func (h *HTTPServer) Start() error {
	clientRg := h.engine.Group("/client")
	if err := h.clientReg.Register(clientRg); err != nil {
		h.logger.Error("failed to register HTTP handlers", zap.Error(err))
		return err
	}
	adminRg := h.engine.Group("/admin")
	if err := h.adminReg.Register(adminRg); err != nil {
		h.logger.Error("failed to register HTTP handlers", zap.Error(err))
		return err
	}

	h.logger.Info("starting HTTP server", zap.String("address", h.http.Addr))
	if err := h.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		h.logger.Error("failed to start HTTP server", zap.Error(err))
		return err
	}
	h.logger.Info("HTTP server is running", zap.String("address", h.http.Addr))
	return nil
}

func (h *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	h.logger.Info("stopping HTTP server")
	if err := h.http.Shutdown(ctx); err != nil {
		h.logger.Error("failed to gracefully shutdown HTTP server", zap.Error(err))
		return err
	}
	<-ctx.Done()
	h.logger.Info("HTTP server stopped")
	return nil
}
