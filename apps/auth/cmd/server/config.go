package main

import (
	"go.uber.org/zap"
	"mandacode.com/accounts/auth/internal/config"
)

func loadConfig(logger *zap.Logger) *config.Config {
	cfg, err := config.LoadConfig()
	if cfg == nil {
		logger.Fatal("failed to load config")
	}
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}
	return cfg
}
