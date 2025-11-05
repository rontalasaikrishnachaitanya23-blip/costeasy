// internal/auth/service/auth_service_interface.go
package service

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/auth/handler/dto"
	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	// Authentication
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error

	// User Management
	GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error
	ListUsers(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*dto.UserResponse, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// Role Management
	AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
	ChangeRole(ctx context.Context, userID, roleID uuid.UUID) error
	RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error

	// Role CRUD
	CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*dto.RoleResponse, error)
	UpdateRole(ctx context.Context, roleID uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error)
	GetRole(ctx context.Context, roleID uuid.UUID) (*dto.RoleResponse, error)
	ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*dto.RoleResponse, error)
	DeleteRole(ctx context.Context, roleID uuid.UUID) error

	// Permission Management
	AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	RemovePermissionsFromRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*dto.PermissionResponse, error)

	// Permission Check
	CheckPermission(ctx context.Context, userID uuid.UUID, module, resource, action string) (bool, error)
}
