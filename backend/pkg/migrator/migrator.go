// backend/pkg/migrator/migrator.go
package migrator

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrator handles database migrations
type Migrator struct {
	pool           *pgxpool.Pool
	migrationsPath string
}

// NewMigrator creates a new migrator instance
func NewMigrator(pool *pgxpool.Pool, migrationsPath string) *Migrator {
	return &Migrator{
		pool:           pool,
		migrationsPath: migrationsPath,
	}
}

// Up runs all pending migrations
func (m *Migrator) Up(ctx context.Context) error {
	log.Println("Running migrations...")

	// Create migrations table if not exists
	if err := m.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrations, err := m.getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedMap := make(map[string]bool)
	for _, a := range applied {
		appliedMap[a] = true
	}

	// Run pending migrations
	count := 0
	for _, migration := range migrations {
		if !strings.HasSuffix(migration, ".up.sql") {
			continue
		}

		name := migration[:len(migration)-7] // Remove .up.sql
		if appliedMap[name] {
			log.Printf("✓ Already applied: %s", name)
			continue
		}

		// Read migration file
		filePath := filepath.Join(m.migrationsPath, migration)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", migration, err)
		}

		// Execute migration
		if _, err := m.pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", name, err)
		}

		// Record migration
		if err := m.recordMigration(ctx, name); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}

		log.Printf("✓ Applied: %s", name)
		count++
	}

	if count == 0 {
		log.Println("✓ No pending migrations")
	} else {
		log.Printf("✓ Applied %d migrations", count)
	}

	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down(ctx context.Context, steps int) error {
	if steps == 0 {
		steps = 1
	}

	log.Printf("Rolling back %d migration(s)...", steps)

	// Get applied migrations in reverse order
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		log.Println("✓ No migrations to rollback")
		return nil
	}

	// Rollback migrations
	for i := 0; i < steps && i < len(applied); i++ {
		name := applied[len(applied)-1-i]

		// Find corresponding .down.sql file
		downFile := name + ".down.sql"
		filePath := filepath.Join(m.migrationsPath, downFile)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read rollback file %s: %w", downFile, err)
		}

		// Execute rollback
		if _, err := m.pool.Exec(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute rollback %s: %w", name, err)
		}

		// Remove from history
		if err := m.removeMigration(ctx, name); err != nil {
			return fmt.Errorf("failed to remove migration record %s: %w", name, err)
		}

		log.Printf("✓ Rolled back: %s", name)
	}

	return nil
}

// Version returns current migration version
func (m *Migrator) Version(ctx context.Context) (string, error) {
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get applied migrations: %w", err)
	}

	if len(applied) == 0 {
		return "0", nil
	}

	return applied[len(applied)-1], nil
}

// Helper functions

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

func (m *Migrator) getMigrationFiles() ([]string, error) {
	files, err := os.ReadDir(m.migrationsPath)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if !file.IsDir() {
			migrations = append(migrations, file.Name())
		}
	}

	sort.Strings(migrations)
	return migrations, nil
}

func (m *Migrator) getAppliedMigrations(ctx context.Context) ([]string, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := m.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		migrations = append(migrations, version)
	}

	return migrations, rows.Err()
}

func (m *Migrator) recordMigration(ctx context.Context, name string) error {
	query := "INSERT INTO schema_migrations (version) VALUES ($1)"
	_, err := m.pool.Exec(ctx, query, name)
	return err
}

func (m *Migrator) removeMigration(ctx context.Context, name string) error {
	query := "DELETE FROM schema_migrations WHERE version = $1"
	_, err := m.pool.Exec(ctx, query, name)
	return err
}
