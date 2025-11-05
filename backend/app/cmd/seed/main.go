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

	cfg := config.LoadConfig()
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	log.Println("Connected to database")

	// Seed in order
	if err := seedOrganizations(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed organizations: %v", err)
	}

	if err := seedUsers(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}

	if err := seedRoles(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed roles: %v", err)
	}

	if err := seedPermissions(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed permissions: %v", err)
	}

	if err := seedRolePermissions(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed role permissions: %v", err)
	}

	if err := seedGLAccounts(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed GL accounts: %v", err)
	}

	if err := seedShafafiyaSettings(ctx, dbPool); err != nil {
		log.Fatalf("Failed to seed Shafafiya settings: %v", err)
	}

	log.Println("✓ Database seeding completed successfully")
	log.Println("\n📝 Default Credentials:")
	log.Println("   Username: admin")
	log.Println("   Password: password123")
}
func seedShafafiyaSettings(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding Shafafiya settings...")

	// Create healthcare organization in Abu Dhabi
	healthcareOrgID := uuid.MustParse("00000000-0000-0000-0000-000000000002")

	orgQuery := `
        INSERT INTO organizations (
            id, name, display_name, email, phone, address, city, 
            state, country, postal_code, is_active
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (id) DO NOTHING
    `

	_, err := db.Exec(ctx, orgQuery,
		healthcareOrgID,
		"Abu Dhabi Health Clinic",
		"Abu Dhabi Health Clinic LLC",
		"contact@adhealthclinic.ae",
		"+971-2-1234567",
		"Corniche Road, Building 23",
		"Abu Dhabi",
		"Abu Dhabi",
		"UAE",
		"12345",
		true,
	)
	if err != nil {
		return err
	}

	// Add Shafafiya settings
	shafafiyaQuery := `
        INSERT INTO shafafiya_org_settings (
            id, organization_id, username, password_encrypted, provider_code,
            default_currency_code, default_language, include_sensitive_data,
            costing_method, allocation_method, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        ON CONFLICT (organization_id) DO NOTHING
    `

	// Note: Use actual crypto service in production
	_, err = db.Exec(ctx, shafafiyaQuery,
		uuid.New(),
		healthcareOrgID,
		"demo_clinic_user",
		"ENCRYPTED_PASSWORD_PLACEHOLDER", // Replace with actual encrypted password
		"PROV001",
		"AED",
		"en",
		true,
		"departmental",
		"weighted",
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	log.Println("✓ Shafafiya settings seeded")
	return nil
}

func seedRolePermissions(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding role permissions...")

	adminRoleID := uuid.MustParse("00000000-0000-0000-0000-000000000020")

	// Get all permissions
	permQuery := `SELECT id FROM permissions`
	rows, err := db.Query(ctx, permQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var permissionIDs []uuid.UUID
	for rows.Next() {
		var permID uuid.UUID
		if err := rows.Scan(&permID); err != nil {
			return err
		}
		permissionIDs = append(permissionIDs, permID)
	}

	// Assign all permissions to admin role
	rolePermQuery := `
        INSERT INTO role_permissions (role_id, permission_id)
        VALUES ($1, $2)
        ON CONFLICT (role_id, permission_id) DO NOTHING
    `

	for _, permID := range permissionIDs {
		_, err := db.Exec(ctx, rolePermQuery, adminRoleID, permID)
		if err != nil {
			return err
		}
	}

	log.Println("✓ Role permissions seeded")
	return nil
}
func seedGLAccounts(ctx context.Context, db *pgxpool.Pool) error {
	log.Println("Seeding GL accounts...")

	accounts := []struct {
		code               string
		name               string
		accType            string
		parentCode         *string
		costCenterRequired bool
	}{
		{"AST-1000", "Cash", "ASSET", nil, false},
		{"AST-1100", "Accounts Receivable", "ASSET", nil, false},
		{"AST-1200", "Inventory", "ASSET", nil, true},
		{"LIA-2000", "Accounts Payable", "LIABILITY", nil, false},
		{"LIA-2100", "Accrued Expenses", "LIABILITY", nil, false},
		{"EQU-3000", "Capital", "EQUITY", nil, false},
		{"EQU-3100", "Retained Earnings", "EQUITY", nil, false},
		{"REV-4000", "Sales Revenue", "REVENUE", nil, false},
		{"REV-4100", "Service Revenue", "REVENUE", nil, false},
		{"EXP-5000", "Cost of Goods Sold", "EXPENSE", nil, true},
		{"EXP-6000", "Salaries Expense", "EXPENSE", nil, false},
		{"EXP-6100", "Rent Expense", "EXPENSE", nil, false},
		{"EXP-6200", "Utilities Expense", "EXPENSE", nil, false},
	}

	query := `
        INSERT INTO gl_accounts (
            id, code, name, type, parent_code, cost_center_required,
            is_active, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (code) DO NOTHING
    `

	for _, acc := range accounts {
		_, err := db.Exec(ctx, query,
			uuid.New(),
			acc.code,
			acc.name,
			acc.accType,
			acc.parentCode,
			acc.costCenterRequired,
			true,
			time.Now(),
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	log.Println("✓ GL accounts seeded")
	return nil
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
