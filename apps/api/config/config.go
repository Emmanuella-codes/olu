package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RedisURL       string
	JWTSecret      string
	AdminJWTSecret string
	Environment    string
	SMSProvider    string
	SMSBaseURL     string
	SMSFrom        string
	SMSTimeoutSec  int
}

func Load() (*Config, error) {
	smsTimeoutSec, err := getEnvInt("SMS_TIMEOUT", 5)
	if err != nil {
		return nil, fmt.Errorf("SMS_TIMEOUT_SECONDS is invalid: %w", err)
	}

	cfg := &Config{
		Port:           getEnv("PORT", "4006"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		RedisURL:       getEnv("REDIS_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		AdminJWTSecret: getEnv("ADMIN_JWT_SECRET", ""),
		Environment:    getEnv("ENVIRONMENT", "development"),
		SMSProvider:    getEnv("SMS_PROVIDER", "mock"),
		SMSBaseURL:     getEnv("SMS_BASE_URL", "http://localhost:3001"),
		SMSFrom:        getEnv("SMS_FROM", "OLU"),
		SMSTimeoutSec:  smsTimeoutSec,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("REDIS_URL is not set")
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not set")
	}

	if cfg.AdminJWTSecret == "" {
		return nil, fmt.Errorf("ADMIN_JWT_SECRET is not set")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) (int, error) {
	if v := os.Getenv(key); v != "" {
		return strconv.Atoi(v)
	}
	return fallback, nil
}
