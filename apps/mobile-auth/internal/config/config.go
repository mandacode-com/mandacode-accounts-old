package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `validate:"required,oneof=dev prod"`
	Port           int    `validate:"required,min=1"`
	AuthServiceURL string `validate:"required"`
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(v *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	config := &Config{
		Env:            getEnv("ENV", "dev"),
		Port:           port,
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", ""),
	}

	if err := v.Struct(config); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			return nil, errors.New("config validation failed: " + validationErrors.Error())
		}
		return nil, errors.New("config validation failed: " + err.Error())
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
