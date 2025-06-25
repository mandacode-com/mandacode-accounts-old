package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env             string
	Port            int
	TokenServiceURL string
	DatabaseURL     string
}

// Validate checks if the config is valid
func (c *Config) Validate() error {
	if c.Env != "dev" && c.Env != "prod" {
		return InvalidConfigError("Environment must be either 'dev' or 'prod'")
	}
	if c.Port <= 0 {
		return InvalidConfigError("Port must be a positive integer")
	}
	if c.TokenServiceURL == "" {
		return InvalidConfigError("Token Service URL must be set")
	}
	if c.DatabaseURL == "" {
		return InvalidConfigError("Database URL must be set")
	}
	return nil
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig() (*Config, error) {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	port, err := strconv.Atoi(getEnv("PORT", "50051"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Env:             getEnv("ENV", "dev"),
		Port:            port,
		TokenServiceURL: getEnv("TOKEN_SERVICE_URL", "localhost:50051"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/accounts_auth?sslmode=disable"),
	}

	if err := config.Validate(); err != nil {
		return nil, err
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
