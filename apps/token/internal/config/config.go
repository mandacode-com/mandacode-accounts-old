package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Env  string
	Port string

	AccessPrivateKey    string
	AccessTokenDuration time.Duration

	RefreshPrivateKey    string
	RefreshTokenDuration time.Duration

	EmailVerificationPrivateKey    string
	EmailVerificationTokenDuration time.Duration
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig() *Config {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	accessTokenDuration, err := time.ParseDuration(getEnv("ACCESS_TOKEN_DURATION", "15m"))
	if err != nil {
		accessTokenDuration = 15 * time.Minute // default to 15 minutes
	}
	refreshTokenDuration, err := time.ParseDuration(getEnv("REFRESH_TOKEN_DURATION", "720h"))
	if err != nil {
		refreshTokenDuration = 720 * time.Hour // default to 30 days
	}
	emailVerificationTokenDuration, err := time.ParseDuration(getEnv("EMAIL_VERIFICATION_TOKEN_DURATION", "168h"))
	if err != nil {
		emailVerificationTokenDuration = 168 * time.Hour // default to 7 days
	}

	return &Config{
		Env:                            getEnv("ENV", "local"),
		Port:                           getEnv("PORT", "50051"),
		AccessPrivateKey:               getEnv("ACCESS_PRIVATE_KEY", ""),
		AccessTokenDuration:            accessTokenDuration,
		RefreshPrivateKey:              getEnv("REFRESH_PRIVATE_KEY", ""),
		RefreshTokenDuration:           refreshTokenDuration,
		EmailVerificationPrivateKey:    getEnv("EMAIL_VERIFICATION_PRIVATE_KEY", ""),
		EmailVerificationTokenDuration: emailVerificationTokenDuration,
	}
}

// getEnv returns env value or fallback
func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
