// backend/app/config/config.go
package config

import (
	"os"
)

type Config struct {
	// Database
	DatabaseURL string

	// Encryption
	EncryptionKey string

	// Server
	Port string
	Host string

	// Logging
	LogLevel    string
	Environment string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/costeasy?sslmode=disable"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		Port:          getEnv("PORT", "8080"),
		Host:          getEnv("HOST", "0.0.0.0"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
