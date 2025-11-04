// backend/app/cmd/seed/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/chaitu35/costeasy/backend/app/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("Starting database seeding...")

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Connected to database")

	// Seed organizations
	if err := seedOrganizations(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed organizations: %v", err)
	}

	// Seed users
	if err := seedUsers(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	// Seed roles
	if err := seedRoles(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed roles: %v", err)
	}

	// Seed permissions
	if err := seedPermissions(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed permissions: %v", err)
	}

	log.Println("✓ Database seeding completed successfully")
}

func seedOrganizations(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding organizations...")

	query := `
		INSERT INTO organizations (id, name, is_active)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO NOTHING
	`

	orgID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	_, err := db.Exec(ctx, query, orgID, "Demo Organization", true)
	if err != nil {
		return err
	}

	log.Println("✓ Organizations seeded")
	return nil
}

func seedUsers(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding users...")

	// Hash password: "password123"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	orgID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000010")

	query := `
		INSERT INTO users (
			id, organization_id, username, email, password_hash,
			first_name, last_name, is_active, is_verified,
			password_changed_at, allow_remote_access
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (username) DO NOTHING
	`

	_, err = db.Exec(ctx, query,
		adminID,
		orgID,
		"admin",
		"admin@demo.com",
		string(passwordHash),
		"Admin",
		"User",
		true,
		true,
		time.Now(),
		true, // Allow remote access for admin
	)

	if err != nil {
		return err
	}

	log.Println("✓ Users seeded (username: admin, password: password123)")
	return nil
}

func seedRoles(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding roles...")

	orgID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	roles := []struct {
		id          uuid.UUID
		name        string
		displayName string
		description string
		isSystem    bool
	}{
		{
			id:          uuid.MustParse("00000000-0000-0000-0000-000000000020"),
			name:        "admin",
			displayName: "Administrator",
			description: "Full system access",
			isSystem:    true,
		},
		{
			id:          uuid.MustParse("00000000-0000-0000-0000-000000000021"),
			name:        "accountant",
			displayName: "Accountant",
			description: "GL and accounting access",
			isSystem:    true,
		},
		{
			id:          uuid.MustParse("00000000-0000-0000-0000-000000000022"),
			name:        "viewer",
			displayName: "Viewer",
			description: "Read-only access",
			isSystem:    true,
		},
	}

	query := `
		INSERT INTO roles (id, organization_id, name, display_name, description, is_system_role, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (organization_id, name) DO NOTHING
	`

	for _, role := range roles {
		_, err := db.Exec(ctx, query,
			role.id,
			orgID,
			role.name,
			role.displayName,
			role.description,
			role.isSystem,
			true,
		)
		if err != nil {
			return err
		}
	}

	// Assign admin role to admin user
	adminID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	adminRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000020")

	userRoleQuery := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := db.Exec(ctx, userRoleQuery, adminID, adminRoleID)
	if err != nil {
		return err
	}

	log.Println("✓ Roles seeded")
	return nil
}

func seedPermissions(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding permissions...")

	permissions := []struct {
		module      string
		resource    string
		action      string
		displayName string
		description string
	}{
		// GL permissions
		{"gl", "accounts", "view", "View Accounts", "View chart of accounts"},
		{"gl", "accounts", "create", "Create Accounts", "Create new accounts"},
		{"gl", "accounts", "edit", "Edit Accounts", "Edit existing accounts"},
		{"gl", "accounts", "delete", "Delete Accounts", "Delete accounts"},
		{"gl", "journal_entries", "view", "View Journal Entries", "View journal entries"},
		{"gl", "journal_entries", "create", "Create Journal Entries", "Create journal entries"},
		{"gl", "journal_entries", "edit", "Edit Journal Entries", "Edit journal entries"},
		{"gl", "journal_entries", "delete", "Delete Journal Entries", "Delete journal entries"},

		// Auth permissions
		{"auth", "users", "view", "View Users", "View user list"},
		{"auth", "users", "create", "Create Users", "Create new users"},
		{"auth", "users", "edit", "Edit Users", "Edit user details"},
		{"auth", "users", "delete", "Delete Users", "Delete users"},
		{"auth", "roles", "view", "View Roles", "View roles"},
		{"auth", "roles", "create", "Create Roles", "Create roles"},
		{"auth", "roles", "edit", "Edit Roles", "Edit roles"},
		{"auth", "roles", "delete", "Delete Roles", "Delete roles"},

		// Settings permissions
		{"settings", "organization", "view", "View Organization", "View organization settings"},
		{"settings", "organization", "edit", "Edit Organization", "Edit organization settings"},
	}

	query := `
		INSERT INTO permissions (id, module, resource, action, display_name, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (module, resource, action) DO NOTHING
	`

	for _, perm := range permissions {
		_, err := db.Exec(ctx, query,
			uuid.New(),
			perm.module,
			perm.resource,
			perm.action,
			perm.displayName,
			perm.description,
		)
		if err != nil {
			return err
		}
	}

	log.Println("✓ Permissions seeded")
	return nil
}
