// backend/internal/auth/repository/permission_repository.go
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

// PermissionRepository implements PermissionRepositoryInterface using pgxpool
type PermissionRepository struct {
	db *pgxpool.Pool
}

// NewPermissionRepository creates a new PermissionRepository
func NewPermissionRepository(db *pgxpool.Pool) PermissionRepositoryInterface {
	return &PermissionRepository{db: db}
}

// Create creates a new permission
func (r *PermissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	query := `
		INSERT INTO permissions (
			id, module, resource, action,
			display_name, description
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		permission.ID,
		permission.Module,
		permission.Resource,
		permission.Action,
		permission.DisplayName,
		permission.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

// GetByID retrieves a permission by ID
func (r *PermissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	query := `
		SELECT id, module, resource, action,
		       display_name, description, created_at
		FROM permissions
		WHERE id = $1
	`

	permission := &domain.Permission{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&permission.ID,
		&permission.Module,
		&permission.Resource,
		&permission.Action,
		&permission.DisplayName,
		&permission.Description,
		&permission.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return permission, nil
}

// GetByKey retrieves a permission by module:resource:action key
func (r *PermissionRepository) GetByKey(ctx context.Context, module, resource, action string) (*domain.Permission, error) {
	query := `
		SELECT id, module, resource, action,
		       display_name, description, created_at
		FROM permissions
		WHERE module = $1 AND resource = $2 AND action = $3
	`

	permission := &domain.Permission{}
	err := r.db.QueryRow(ctx, query, module, resource, action).Scan(
		&permission.ID,
		&permission.Module,
		&permission.Resource,
		&permission.Action,
		&permission.DisplayName,
		&permission.Description,
		&permission.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPermissionNotFound
		}
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	return permission, nil
}

// List retrieves all permissions
func (r *PermissionRepository) List(ctx context.Context) ([]*domain.Permission, error) {
	query := `
		SELECT id, module, resource, action,
		       display_name, description, created_at
		FROM permissions
		ORDER BY module, resource, action
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		permission := &domain.Permission{}
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

// ListByModule retrieves permissions by module
func (r *PermissionRepository) ListByModule(ctx context.Context, module string) ([]*domain.Permission, error) {
	query := `
		SELECT id, module, resource, action,
		       display_name, description, created_at
		FROM permissions
		WHERE module = $1
		ORDER BY resource, action
	`

	rows, err := r.db.Query(ctx, query, module)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions by module: %w", err)
	}
	defer rows.Close()

	var permissions []*domain.Permission
	for rows.Next() {
		permission := &domain.Permission{}
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

// Delete deletes a permission
func (r *PermissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM permissions WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrPermissionNotFound
	}

	return nil
}

// GetUserPermissions gets all permissions for a user (via roles)
func (r *PermissionRepository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.module, p.resource, p.action,
		       p.display_name, p.description, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		INNER JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND r.is_active = true
		ORDER BY p.module, p.resource, p.action
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
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

// CheckUserPermission checks if user has specific permission
func (r *PermissionRepository) CheckUserPermission(ctx context.Context, userID uuid.UUID, module, resource, action string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			INNER JOIN user_roles ur ON rp.role_id = ur.role_id
			INNER JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 
			  AND p.module = $2 
			  AND p.resource = $3 
			  AND p.action = $4
			  AND r.is_active = true
		)
	`

	var hasPermission bool
	err := r.db.QueryRow(ctx, query, userID, module, resource, action).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("failed to check user permission: %w", err)
	}

	return hasPermission, nil
}

// BulkCreate creates multiple permissions
func (r *PermissionRepository) BulkCreate(ctx context.Context, permissions []*domain.Permission) error {
	if len(permissions) == 0 {
		return nil
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO permissions (
			id, module, resource, action,
			display_name, description
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (module, resource, action) DO NOTHING
	`

	for _, permission := range permissions {
		_, err := tx.Exec(ctx, query,
			permission.ID,
			permission.Module,
			permission.Resource,
			permission.Action,
			permission.DisplayName,
			permission.Description,
		)
		if err != nil {
			return fmt.Errorf("failed to create permission: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
