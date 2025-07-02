package httpserver

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mandacode.com/lib/server/server"
)

type HTTPServer struct {
	http      *http.Server
	engine    *gin.Engine
	auth_reg  server.HTTPRegisterer
	token_reg server.HTTPRegisterer
	logger    *zap.Logger
	port      int
}

func NewHTTPServer(port int, logger *zap.Logger, auth_reg server.HTTPRegisterer, token_reg server.HTTPRegisterer) (*HTTPServer, error) {
	engine := gin.Default()

	return &HTTPServer{
		http:      &http.Server{Addr: ":" + strconv.Itoa(port), Handler: engine},
		engine:    engine,
		auth_reg:  auth_reg,
		token_reg: token_reg,
		logger:    logger,
		port:      port,
	}, nil
}

func (h *HTTPServer) Start() error {
	auth_rg := h.engine.Group("/")
	if err := h.auth_reg.Register(auth_rg); err != nil {
		h.logger.Error("failed to register HTTP handlers", zap.Error(err))
		return err
	}
	token_rg := h.engine.Group("/token")
	if err := h.token_reg.Register(token_rg); err != nil {
		h.logger.Error("failed to register token handlers", zap.Error(err))
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
