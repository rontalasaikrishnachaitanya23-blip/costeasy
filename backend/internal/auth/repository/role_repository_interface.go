package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
)

type RoleRepositoryInterface interface {
	Create(ctx context.Context, role *domain.Role) error
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, roleID uuid.UUID) error
	GetByID(ctx context.Context, roleID uuid.UUID) (*domain.Role, error)
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Role, error)

	// User role management
	AssignToUser(ctx context.Context, userID, roleID uuid.UUID) error
	ChangeUserRole(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveUserRole(ctx context.Context, userID, roleID uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*domain.Role, error)

	// Permission management for roles
	AssignPermission(ctx context.Context, roleID, permID uuid.UUID, grantedBy *uuid.UUID) error
	RemovePermission(ctx context.Context, roleID, permID uuid.UUID) error
}
