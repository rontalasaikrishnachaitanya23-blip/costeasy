// backend/database/db.go
package database

import (
    "context"
    "fmt"

    "github.com/chaitu35/costeasy/backend/app/config"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Pool is an alias for pgxpool.Pool
type Pool = *pgxpool.Pool

// ConnectDB establishes a connection to the database
func ConnectDB(cfg *config.Config) (Pool, error) {
    // Use DATABASE_URL directly
    connString := cfg.DatabaseURL

    pool, err := pgxpool.New(context.Background(), connString)
    if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }

    // Test connection
    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("unable to ping database: %w", err)
    }

    return pool, nil
}
