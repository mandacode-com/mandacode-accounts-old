package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                         string
	Port                        string
	AccessPublicKey             string
	AccessPrivateKey            string
	RefreshPublicKey            string
	RefreshPrivateKey           string
	EmailVerificationPublicKey  string
	EmailVerificationPrivateKey string
}

// LoadConfig loads env vars from .env (if exists) and returns structured config
func LoadConfig() *Config {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	return &Config{
		Env:                         getEnv("ENV", "local"),
		Port:                        getEnv("PORT", "50051"),
		AccessPublicKey:             getEnv("ACCESS_PUBLIC_KEY", ""),
		AccessPrivateKey:            getEnv("ACCESS_PRIVATE_KEY", ""),
		RefreshPublicKey:            getEnv("REFRESH_PUBLIC_KEY", ""),
		RefreshPrivateKey:           getEnv("REFRESH_PRIVATE_KEY", ""),
		EmailVerificationPublicKey:  getEnv("EMAIL_VERIFICATION_PUBLIC_KEY", ""),
		EmailVerificationPrivateKey: getEnv("EMAIL_VERIFICATION_PRIVATE_KEY", ""),
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
