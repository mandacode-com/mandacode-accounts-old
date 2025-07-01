package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type GoogleOAuthConfig struct {
	ClientID     string `validate:"required"`
	ClientSecret string `validate:"required"`
	RedirectURL  string `validate:"required,url"`
}

type NaverOAuthConfig struct {
	ClientID     string `validate:"required"`
	ClientSecret string `validate:"required"`
	RedirectURL  string `validate:"required,url"`
}

type KakaoOAuthConfig struct {
	ClientID     string `validate:"required"`
	ClientSecret string `validate:"required"`
	RedirectURL  string `validate:"required,url"`
}

type RedisConfig struct {
	Address  string `validate:"required"`
	Password string `validate:"omitempty"`
	DB       int    `validate:"numeric,min=0,max=15"`
}

type Config struct {
	Env               string            `validate:"required,oneof=dev prod"`
	Port              int               `validate:"required,min=1"`
	AuthServiceURL    string            `validate:"required"`
	GoogleOAuthConfig GoogleOAuthConfig `validate:"required"`
	NaverOAuthConfig  NaverOAuthConfig  `validate:"required"`
	KakaoOAuthConfig  KakaoOAuthConfig  `validate:"required"`
	CodeStore         RedisConfig       `validate:"required"`
	SessionStore      RedisConfig       `validate:"required"`
	CodeTTL           time.Duration     `validate:"required,numeric,min=1"` // TTL for code storage
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

	codeStoreDB, err := strconv.Atoi(getEnv("CODE_STORE_DB", "0"))
	if err != nil {
		return nil, errors.New("invalid CODE_STORE_DB value: " + err.Error())
	}
	sessionStoreDB, err := strconv.Atoi(getEnv("SESSION_STORE_DB", "0"))
	if err != nil {
		return nil, errors.New("invalid SESSION_STORE_DB value: " + err.Error())
	}

	codeTTL, err := time.ParseDuration(getEnv("CODE_TTL", "1m"))

	config := &Config{
		Env:            getEnv("ENV", "dev"),
		Port:           port,
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "localhost:50000"),
		GoogleOAuthConfig: GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
		},
		NaverOAuthConfig: NaverOAuthConfig{
			ClientID:     getEnv("NAVER_CLIENT_ID", ""),
			ClientSecret: getEnv("NAVER_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("NAVER_REDIRECT_URL", "http://localhost:8080/auth/naver/callback"),
		},
		KakaoOAuthConfig: KakaoOAuthConfig{
			ClientID:     getEnv("KAKAO_CLIENT_ID", ""),
			ClientSecret: getEnv("KAKAO_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("KAKAO_REDIRECT_URL", "http://localhost:8080/auth/kakao/callback"),
		},
		CodeStore: RedisConfig{
			Address:  getEnv("CODE_STORE_ADDRESS", "localhost:6379"),
			Password: getEnv("CODE_STORE_PASSWORD", ""),
			DB:       codeStoreDB,
		},
		SessionStore: RedisConfig{
			Address:  getEnv("SESSION_STORE_ADDRESS", "localhost:6379"),
			Password: getEnv("SESSION_STORE_PASSWORD", ""),
			DB:       sessionStoreDB,
		},
		CodeTTL: codeTTL,
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
