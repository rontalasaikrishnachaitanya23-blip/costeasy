// backend/internal/auth/repository/role_repository_interface.go
package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
)

// RoleRepositoryInterface defines the contract for role data operations
type RoleRepositoryInterface interface {
	// Create creates a new role
	Create(ctx context.Context, role *domain.Role) error
	
	// GetByID retrieves a role by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	
	// GetByName retrieves a role by name within an organization
	GetByName(ctx context.Context, orgID uuid.UUID, name string) (*domain.Role, error)
	
	// Update updates an existing role
	Update(ctx context.Context, role *domain.Role) error
	
	// Delete deletes a role (only if not system role)
	Delete(ctx context.Context, id uuid.UUID) error
	
	// List retrieves roles with pagination
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Role, error)
	
	// Count returns total number of roles in organization
	Count(ctx context.Context, orgID uuid.UUID) (int, error)
	
	// GetWithPermissions retrieves a role with its permissions
	GetWithPermissions(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	
	// AssignPermission assigns a permission to role
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID, grantedBy *uuid.UUID) error
	
	// RemovePermission removes a permission from role
	RemovePermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	
	// GetPermissions gets all permissions for a role
	GetPermissions(ctx context.Context, roleID uuid.UUID) ([]domain.Permission, error)
	
	// GetSystemRoles gets all system roles
	GetSystemRoles(ctx context.Context, orgID uuid.UUID) ([]*domain.Role, error)
}
