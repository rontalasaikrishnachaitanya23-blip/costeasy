package migrator

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrator handles file-based SQL migrations using pgxpool
type Migrator struct {
	pool           *pgxpool.Pool
	migrationsPath string
}

// NewMigrator initializes the migrator with a DB pool and migrations directory
func NewMigrator(pool *pgxpool.Pool, migrationsPath string) *Migrator {
	return &Migrator{
		pool:           pool,
		migrationsPath: migrationsPath,
	}
}

// Up applies all pending .up.sql migrations in order
func (m *Migrator) Up(ctx context.Context) error {
	log.Println("🚀 Starting database migrations...")

	// Ensure tracking table exists
	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create schema_migrations: %w", err)
	}

	// Load migration files
	files, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, a := range applied {
		appliedMap[a] = true
	}

	appliedCount := 0
	for _, file := range files {
		if !strings.HasSuffix(file, ".up.sql") {
			continue
		}

		name := strings.TrimSuffix(file, ".up.sql")
		if appliedMap[name] {
			log.Printf("⏩ Skipping already applied migration: %s", name)
			continue
		}

		if err := m.runMigration(ctx, name, file, true); err != nil {
			return err
		}
		appliedCount++
	}

	if appliedCount == 0 {
		log.Println("✅ No pending migrations.")
	} else {
		log.Printf("✅ Applied %d new migration(s).", appliedCount)
	}
	return nil
}

// Down rolls back the most recent migrations (default: 1 step)
func (m *Migrator) Down(ctx context.Context, steps int) error {
	if steps <= 0 {
		steps = 1
	}

	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	if len(applied) == 0 {
		log.Println("⚠️  No migrations to rollback.")
		return nil
	}

	rollbackCount := 0
	for i := 0; i < steps && i < len(applied); i++ {
		name := applied[len(applied)-1-i]
		file := name + ".down.sql"

		if err := m.runMigration(ctx, name, file, false); err != nil {
			return err
		}
		rollbackCount++
	}

	log.Printf("🌀 Rolled back %d migration(s).", rollbackCount)
	return nil
}

// Version prints the latest applied migration version
func (m *Migrator) Version(ctx context.Context) (string, error) {
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch applied migrations: %w", err)
	}
	if len(applied) == 0 {
		return "none", nil
	}
	return applied[len(applied)-1], nil
}

// ---------------------------------------------------------------------
// 🧩 Internal Helper Methods
// ---------------------------------------------------------------------

// runMigration executes one migration inside a transaction
func (m *Migrator) runMigration(ctx context.Context, name, file string, isUp bool) error {
	filePath := filepath.Join(m.migrationsPath, file)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", file, err)
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction for %s: %w", name, err)
	}
	defer tx.Rollback(ctx)

	log.Printf("➡️  Executing %s ...", file)
	start := time.Now()

	if _, err := tx.Exec(ctx, string(content)); err != nil {
		return fmt.Errorf("migration %s failed: %w", name, err)
	}

	if isUp {
		if err := m.recordMigration(ctx, name); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}
	} else {
		if err := m.removeMigration(ctx, name); err != nil {
			return fmt.Errorf("failed to remove migration %s: %w", name, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", name, err)
	}

	log.Printf("✅ Completed: %s (in %v)", file, time.Since(start))
	return nil
}

// createMigrationsTable ensures the schema_migrations table exists
func (m *Migrator) createMigrationsTable(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		id SERIAL PRIMARY KEY,
		version VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP DEFAULT NOW()
	);
	`
	_, err := m.pool.Exec(ctx, query)
	return err
}

// getMigrationFiles returns ordered .sql files from the directory
func (m *Migrator) getMigrationFiles() ([]string, error) {
	path, err := filepath.Abs(m.migrationsPath)
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, f := range files {
		if !f.IsDir() && (strings.HasSuffix(f.Name(), ".up.sql") || strings.HasSuffix(f.Name(), ".down.sql")) {
			migrations = append(migrations, f.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

// getAppliedMigrations fetches all versions already in schema_migrations
func (m *Migrator) getAppliedMigrations(ctx context.Context) ([]string, error) {
	rows, err := m.pool.Query(ctx, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied = append(applied, v)
	}
	return applied, rows.Err()
}

// recordMigration inserts a version record after success
func (m *Migrator) recordMigration(ctx context.Context, name string) error {
	_, err := m.pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", name)
	return err
}

// removeMigration deletes a version record during rollback
func (m *Migrator) removeMigration(ctx context.Context, name string) error {
	_, err := m.pool.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", name)
	return err
}
