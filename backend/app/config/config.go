package config

import (
	"os"
)

type Config struct {
	// Server
	Host        string
	Port        string
	Environment string
	LogLevel    string

	// Database
	DatabaseURL string

	// JWT
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string

	// Add other existing fields...
}

func LoadConfig() *Config {
	return &Config{
		Host:             getEnv("HOST", "0.0.0.0"),
		Port:             getEnv("PORT", "8080"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "debug"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgresql://user:pass@localhost:5432/costeasy?sslmode=disable"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "7d"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
