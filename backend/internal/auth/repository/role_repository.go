// backend/internal/auth/repository/role_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoleRepository implements RoleRepositoryInterface using pgxpool
type RoleRepository struct {
	db *pgxpool.Pool
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *pgxpool.Pool) RoleRepositoryInterface {
	return &RoleRepository{db: db}
}

// Create creates a new role
func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `
		INSERT INTO roles (
			id, organization_id, name, display_name,
			description, is_system_role, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		role.ID,
		role.OrganizationID,
		role.Name,
		role.DisplayName,
		role.Description,
		role.IsSystemRole,
		role.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `
		SELECT 
			id, organization_id, name, display_name,
			description, is_system_role, is_active,
			created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &domain.Role{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.OrganizationID,
		&role.Name,
		&role.DisplayName,
		&role.Description,
		&role.IsSystemRole,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// GetByName retrieves a role by name within an organization
func (r *RoleRepository) GetByName(ctx context.Context, orgID uuid.UUID, name string) (*domain.Role, error) {
	query := `
		SELECT 
			id, organization_id, name, display_name,
			description, is_system_role, is_active,
			created_at, updated_at
		FROM roles
		WHERE organization_id = $1 AND name = $2
	`

	role := &domain.Role{}
	err := r.db.QueryRow(ctx, query, orgID, name).Scan(
		&role.ID,
		&role.OrganizationID,
		&role.Name,
		&role.DisplayName,
		&role.Description,
		&role.IsSystemRole,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRoleNotFound
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

// Update updates an existing role
func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `
		UPDATE roles SET
			name = $2,
			display_name = $3,
			description = $4,
			is_active = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		role.ID,
		role.Name,
		role.DisplayName,
		role.Description,
		role.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}

// Delete deletes a role (only if not system role)
func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if it's a system role
	role, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if role.IsSystemRole {
		return domain.ErrSystemRoleModification
	}

	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}

// List retrieves roles with pagination
func (r *RoleRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Role, error) {
	query := `
		SELECT 
			id, organization_id, name, display_name,
			description, is_system_role, is_active,
			created_at, updated_at
		FROM roles
		WHERE organization_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(
			&role.ID,
			&role.OrganizationID,
			&role.Name,
			&role.DisplayName,
			&role.Description,
			&role.IsSystemRole,
			&role.IsActive,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// Count returns total number of roles in organization
func (r *RoleRepository) Count(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM roles WHERE organization_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count roles: %w", err)
	}

	return count, nil
}

// GetWithPermissions retrieves a role with its permissions
func (r *RoleRepository) GetWithPermissions(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	// Get role
	role, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get permissions
	permissions, err := r.GetPermissions(ctx, id)
	if err != nil {
		return nil, err
	}

	role.Permissions = permissions
	return role, nil
}

// AssignPermission assigns a permission to role
func (r *RoleRepository) AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID, grantedBy *uuid.UUID) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id, granted_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, roleID, permissionID, grantedBy)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}

	return nil
}

// RemovePermission removes a permission from role
func (r *RoleRepository) RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	result, err := r.db.Exec(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("permission assignment not found")
	}

	return nil
}

// GetPermissions gets all permissions for a role
func (r *RoleRepository) GetPermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error) {
	query := `
		SELECT p.id, p.module, p.resource, p.action,
		       p.display_name, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(
			&permission.ID,
			&permission.Module,
			&permission.Resource,
			&permission.Action,
			&permission.DisplayName,
			&permission.Description,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetSystemRoles gets all system roles
func (r *RoleRepository) GetSystemRoles(ctx context.Context, orgID uuid.UUID) ([]*domain.Role, error) {
	query := `
		SELECT 
			id, organization_id, name, display_name,
			description, is_system_role, is_active,
			created_at, updated_at
		FROM roles
		WHERE organization_id = $1 AND is_system_role = true
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get system roles: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(
			&role.ID,
			&role.OrganizationID,
			&role.Name,
			&role.DisplayName,
			&role.Description,
			&role.IsSystemRole,
			&role.IsActive,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetByNameGlobal retrieves a system role by name (no organization scope)
func (r *RoleRepository) GetByNameGlobal(ctx context.Context, name string) (*domain.Role, error) {
	query := `
        SELECT id, organization_id, name, description, is_system_role, is_active, 
               created_at, updated_at, created_by, updated_by
        FROM roles
        WHERE name = $1 
          AND is_system_role = true
          AND is_active = true
        LIMIT 1
    `

	var role domain.Role
	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID,
		&role.OrganizationID,
		&role.Name,
		&role.Description,
		&role.IsSystemRole,
		&role.IsActive,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

// GetUserRoles retrieves all roles for a user
func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*domain.Role, error) {
	query := `
        SELECT r.id, r.organization_id, r.name, r.description, r.is_system_role, 
               r.is_active, r.created_at, r.updated_at, r.created_by, r.updated_by
        FROM roles r
        INNER JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1 AND r.is_active = true
        ORDER BY r.name
    `

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(
			&role.ID,
			&role.OrganizationID,
			&role.Name,
			&role.Description,
			&role.IsSystemRole,
			&role.IsActive,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

// ChangeUserRole updates a user's role (replaces existing role).
func (r *RoleRepository) ChangeUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Remove existing roles
	_, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to remove existing roles: %w", err)
	}

	// Assign new role
	_, err = tx.Exec(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign new role: %w", err)
	}

	return tx.Commit(ctx)
}

// RemoveUserRole removes a specific role from a user.
func (r *RoleRepository) RemoveUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	result, err := r.db.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove user role: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("role not assigned to user")
	}
	return nil
}

// AssignToUser assigns a role to a user.
func (r *RoleRepository) AssignToUser(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}
	return nil
}
