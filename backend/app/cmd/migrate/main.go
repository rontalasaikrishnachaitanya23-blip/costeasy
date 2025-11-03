// backend/app/cmd/migrate/main.go
package main

import (
	"context"
	"flag"
	"log"

	"github.com/joho/godotenv"

	"github.com/chaitu35/costeasy/backend/app/config"
	"github.com/chaitu35/costeasy/backend/database"
	"github.com/chaitu35/costeasy/backend/pkg/migrator"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	command := flag.String("command", "up", "Command: up, down, version")
	steps := flag.Int("steps", 1, "Number of steps for rollback")
	flag.Parse()

	// Load configuration
	cfg := config.LoadConfig()

	// Get database connection using ConnectDB
	pool, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create migrator
	m := migrator.NewMigrator(pool, "database/migrations")

	ctx := context.Background()

	// Execute command
	switch *command {
	case "up":
		log.Println("Running migrations up...")
		if err := m.Up(ctx); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✓ Migrations completed successfully")

	case "down":
		log.Printf("Rolling back %d migration(s)...\n", *steps)
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
