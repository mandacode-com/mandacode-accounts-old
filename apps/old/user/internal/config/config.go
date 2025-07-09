package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type ServiceAddrConfig struct {
	Auth    string `validate:"required"`
	Profile string `validate:"required"`
	Mailer  string `validate:"required"`
	Token   string `validate:"required"`
	Kafka   string `validate:"required"`
}

type Config struct {
	Env                         string            `validate:"required,oneof=dev prod"`
	Port                        int               `validate:"required,min=1,max=65535"`
	DatabaseURL                 string            `validate:"required"`
	UIDHeader                   string            `validate:"required"`
	SyncCodeLength              int               `validate:"required,min=1"`
	EmailVerificationCodeLength int               `validate:"required,min=1"`
	ServiceAddr                 ServiceAddrConfig `validate:"required"`
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig(v *validator.Validate) (*Config, error) {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, errors.New("invalid PORT value: " + err.Error())
	}

	syncCodeLength, err := strconv.Atoi(getEnv("SYNC_CODE_LENGTH", "6"))
	if err != nil {
		return nil, errors.New("invalid SYNC_CODE_LENGTH value: " + err.Error())
	}
	emailVerificationCodeLength, err := strconv.Atoi(getEnv("EMAIL_VERIFICATION_CODE_LENGTH", "6"))
	if err != nil {
		return nil, errors.New("invalid EMAIL_VERIFICATION_CODE_LENGTH value: " + err.Error())
	}

	config := &Config{
		Env:                         getEnv("ENV", "dev"),
		Port:                        port,
		DatabaseURL:                 getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/accounts_role?sslmode=disable"),
		UIDHeader:                   getEnv("UID_HEADER", "X-User-ID"),
		SyncCodeLength:              syncCodeLength,
		EmailVerificationCodeLength: emailVerificationCodeLength,
		ServiceAddr: ServiceAddrConfig{
			Auth:    getEnv("AUTH_SERVICE_ADDR", ""),
			Profile: getEnv("PROFILE_SERVICE_ADDR", ""),
			Mailer:  getEnv("MAILER_SERVICE_ADDR", ""),
			Token:   getEnv("TOKEN_SERVICE_ADDR", ""),
			Kafka:   getEnv("KAFKA_SERVICE_ADDR", ""),
		},
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
