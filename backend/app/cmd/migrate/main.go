// backend/app/cmd/migrate/main.go
package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/chaitu35/costeasy/backend/app/config"
	"github.com/chaitu35/costeasy/backend/database"
	"github.com/chaitu35/costeasy/backend/pkg/migrator"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	command := flag.String("command", "up", "Command: up, down, version")
	steps := flag.Int("steps", 1, "Number of steps for rollback")
	path := flag.String("path", "database/migrations", "Relative path to migrations folder")
	flag.Parse()

	// Allow override via env if needed
	if envPath := os.Getenv("MIGRATIONS_PATH"); envPath != "" {
		*path = envPath
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Get database connection using ConnectDB
	pool, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create migrator using provided path
	m := migrator.NewMigrator(pool, *path)

	ctx := context.Background()

	switch *command {
	case "up":
		log.Printf("Running migrations up (path=%s)...", *path)
		if err := m.Up(ctx); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✓ Migrations completed successfully")

	case "down":
		log.Printf("Rolling back %d migration(s) (path=%s)...", *steps, *path)
		if err := m.Down(ctx, *steps); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✓ Rollback completed successfully")

	case "version":
		version, err := m.Version(ctx)
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current migration version: %s", version)

	default:
		log.Fatalf("Unknown command: %s. Use: up, down, or version", *command)
	}
}
