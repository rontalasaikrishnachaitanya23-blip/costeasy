package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chaitu35/costeasy/backend/app/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool = *pgxpool.Pool

func ConnectDB(cfg *config.Config) (Pool, error) {
	// Build DSN manually
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)

	log.Printf("[DB] Connecting as %s@%s/%s", cfg.DBUser, cfg.DBHost, cfg.DBName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return pool, nil
}
