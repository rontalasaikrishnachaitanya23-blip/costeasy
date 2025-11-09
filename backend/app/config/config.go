package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Host             string
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	Environment      string
	TrustedProxies   []string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // load .env file if exists

	return &Config{
		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSLMode:        os.Getenv("DB_SSLMODE"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
		Environment:      os.Getenv("ENVIRONMENT"),
		TrustedProxies:   []string{os.Getenv("TRUSTED_PROXIES")},
	}
}
