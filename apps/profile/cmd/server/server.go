package servers

import (
	"os"
	"os/signal"

	"go.uber.org/zap"
)

type Server interface {
	Start() error
	Stop() error
}

type ServerManager struct {
	Servers []Server
	logger  *zap.Logger
}

func (sm *ServerManager) Run() {
	for _, server := range sm.Servers {
		go func(s Server) {
			if err := s.Start(); err != nil {
				sm.logger.Error("failed to start server", zap.Error(err))
			}
		}(server)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for {
		select {
		case <-signalChan:
			sm.logger.Info("received shutdown signal, stopping servers")
			for _, server := range sm.Servers {
				if err := server.Stop(); err != nil {
					sm.logger.Error("failed to stop server", zap.Error(err))
				}
			}
			return
		}
	}
}
