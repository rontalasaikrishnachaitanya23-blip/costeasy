// backend/internal/auth/repository/permission_repository_interface.go
package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
)

// PermissionRepositoryInterface defines the contract for permission data operations
type PermissionRepositoryInterface interface {
	// Create creates a new permission
	Create(ctx context.Context, permission *domain.Permission) error

	// GetByID retrieves a permission by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Permission, error)

	// GetByKey retrieves a permission by module:resource:action key
	GetByKey(ctx context.Context, module, resource, action string) (*domain.Permission, error)

	// List retrieves all permissions
	List(ctx context.Context) ([]*domain.Permission, error)

	// ListByModule retrieves permissions by module
	ListByModule(ctx context.Context, module string) ([]*domain.Permission, error)

	// Delete deletes a permission
	Delete(ctx context.Context, id uuid.UUID) error

	// GetUserPermissions gets all permissions for a user (via roles)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]domain.Permission, error)

	// CheckUserPermission checks if user has specific permission
	CheckUserPermission(ctx context.Context, userID uuid.UUID, module, resource, action string) (bool, error)

	// BulkCreate creates multiple permissions
	BulkCreate(ctx context.Context, permissions []*domain.Permission) error
}
