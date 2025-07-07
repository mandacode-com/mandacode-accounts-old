package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string `validate:"required,oneof=dev prod"`
	HTTPPort    int    `validate:"required,min=1"`
	GRPCPort    int    `validate:"required,min=1"`
	DatabaseURL string `validate:"required"`
	UIDHeader   string `validate:"required"` // Header name for user ID
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(v *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "prod" {
		_ = godotenv.Load()
	}

	httpPort, err := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	if err != nil {
		return nil, err
	}
	grpcPort, err := strconv.Atoi(getEnv("GRPC_PORT", "50051"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Env:         getEnv("ENV", "dev"),
		HTTPPort:    httpPort,
		GRPCPort:    grpcPort,
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/accounts_profile?sslmode=disable"),
		UIDHeader:   getEnv("UID_HEADER", "X-User-ID"),
	}

	if err := v.Struct(config); err != nil {
		return nil, errors.New("invalid configuration: " + err.Error())
	}

	return config, nil
}

// getEnv returns env value or fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
