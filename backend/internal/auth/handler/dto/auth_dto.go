package dto

import (
	"time"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
)

// ============================================================================
// Authentication DTOs
// ============================================================================

// LoginRequest for user login
type LoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	IPAddress string `json:"-"`
	UserAgent string `json:"-"`
}

// RegisterRequest for user registration
type RegisterRequest struct {
	Email     string     `json:"email" binding:"required,email"`
	Password  string     `json:"password" binding:"required,min=8"`
	FirstName string     `json:"first_name" binding:"required"`
	LastName  string     `json:"last_name" binding:"required"`
	Username  string     `json:"username" binding:"required"`
	Phone     string     `json:"phone"`
	RoleID    *uuid.UUID `json:"role_id"` // Optional - if not provided, assign default role
}

// RefreshTokenRequest for token refresh
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest for password change
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ForgotPasswordRequest for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest for password reset
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// TokenResponse returned after successful authentication
type TokenResponse struct {
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	TokenType    string               `json:"token_type"`
	ExpiresIn    int                  `json:"expires_in"`
	User         *UserResponse        `json:"user"`
	Roles        []string             `json:"roles"`
	Permissions  []PermissionResponse `json:"permissions"`
}

// ============================================================================
// User Management DTOs
// ============================================================================

// CreateUserRequest for admin creating users
type CreateUserRequest struct {
	Username       string    `json:"username" binding:"required,min=3,max=50"`
	Email          string    `json:"email" binding:"required,email"`
	Password       string    `json:"password" binding:"required,min=8"`
	FirstName      string    `json:"first_name" binding:"required"`
	LastName       string    `json:"last_name" binding:"required"`
	Phone          string    `json:"phone"`
	OrganizationID uuid.UUID `json:"organization_id" binding:"required"`
	RoleID         uuid.UUID `json:"role_id" binding:"required"` // Required for admin user creation
}

// UpdateProfileRequest for updating user profile
type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// UserResponse for user data
type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Phone       string     `json:"phone,omitempty"`
	IsActive    bool       `json:"is_active"`
	IsVerified  bool       `json:"is_verified"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	Token       string     `json:"token,omitempty"`
}

// ToUserResponse converts domain User to UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	firstName := ""
	if user.FirstName != nil {
		firstName = *user.FirstName
	}

	lastName := ""
	if user.LastName != nil {
		lastName = *user.LastName
	}

	phone := ""
	if user.Phone != nil {
		phone = *user.Phone
	}

	return &UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		FirstName:   firstName,
		LastName:    lastName,
		Phone:       phone,
		IsActive:    user.IsActive,
		IsVerified:  user.IsVerified,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
	}
}

// ============================================================================
// Role Management DTOs
// ============================================================================

// CreateRoleRequest for creating roles
type CreateRoleRequest struct {
	Name           string      `json:"name" binding:"required,min=2,max=50"`
	DisplayName    string      `json:"display_name" binding:"required"`
	Description    string      `json:"description"`
	OrganizationID *uuid.UUID  `json:"organization_id"` // nil for system roles
	PermissionIDs  []uuid.UUID `json:"permission_ids"`  // Optional - can be assigned later
}

// UpdateRoleRequest for updating roles
type UpdateRoleRequest struct {
	Name          string      `json:"name"`
	DisplayName   string      `json:"display_name"`
	Description   string      `json:"description"`
	PermissionIDs []uuid.UUID `json:"permission_ids"` // Update permissions (replaces all)
}

// RoleResponse for role data
type RoleResponse struct {
	ID             uuid.UUID            `json:"id"`
	Name           string               `json:"name"`
	DisplayName    string               `json:"display_name"`
	Description    string               `json:"description"`
	OrganizationID *uuid.UUID           `json:"organization_id,omitempty"`
	IsSystemRole   bool                 `json:"is_system_role"`
	IsActive       bool                 `json:"is_active"`
	Permissions    []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt      string               `json:"created_at"`
	UpdatedAt      string               `json:"updated_at"`
}

// ToRoleResponse converts domain Role to RoleResponse
func ToRoleResponse(role *domain.Role, permissions []PermissionResponse) *RoleResponse {
	desc := ""
	if role.Description != nil {
		desc = *role.Description
	}

	var orgID *uuid.UUID
	if role.OrganizationID != uuid.Nil {
		orgID = &role.OrganizationID
	}
	return &RoleResponse{
		ID:             role.ID,
		Name:           role.Name,
		DisplayName:    role.DisplayName,
		Description:    desc,
		OrganizationID: orgID,
		IsSystemRole:   role.IsSystemRole,
		IsActive:       role.IsActive,
		Permissions:    permissions,
		CreatedAt:      role.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      role.UpdatedAt.Format(time.RFC3339),
	}
}

// AssignRoleRequest for assigning role to user
type AssignRoleRequest struct {
	RoleID uuid.UUID `json:"role_id" binding:"required"`
}

// ============================================================================
// Permission Management DTOs
// ============================================================================

// PermissionResponse for permission data
type PermissionResponse struct {
	Module    string `json:"module"`
	Page      string `json:"page"`
	CanView   bool   `json:"can_view"`
	CanAdd    bool   `json:"can_add"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
	CanPrint  bool   `json:"can_print"`
	CanExport bool   `json:"can_export"`
}

// AssignPermissionsRequest for bulk permission assignment
type AssignPermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required"`
}

// CreatePermissionRequest for creating permissions
type CreatePermissionRequest struct {
	Module   string `json:"module" binding:"required"`
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// ============================================================================
// MFA DTOs
// ============================================================================

// VerifyMFARequest for MFA verification
type VerifyMFARequest struct {
	Code string `json:"code" binding:"required"`
}

// EnableMFARequest for enabling MFA
type EnableMFARequest struct {
	Method string `json:"method" binding:"required,oneof=sms email totp"`
}

// EnableMFAResponse returned when MFA is enabled
type EnableMFAResponse struct {
	Secret      string   `json:"secret,omitempty"`       // For TOTP
	QRCode      string   `json:"qr_code,omitempty"`      // For TOTP
	BackupCodes []string `json:"backup_codes,omitempty"` // Backup codes
}

// ============================================================================
// Pagination & List DTOs
// ============================================================================

// PaginationRequest for paginated requests
type PaginationRequest struct {
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
	SortBy   string `json:"sort_by" form:"sort_by"`
	SortDesc bool   `json:"sort_desc" form:"sort_desc"`
}

// PaginatedResponse for paginated responses
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// ListUsersResponse for user list
type ListUsersResponse struct {
	Users      []*UserResponse `json:"users"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ListRolesResponse for role list
type ListRolesResponse struct {
	Roles      []*RoleResponse `json:"roles"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ============================================================================
// Error Response DTOs
// ============================================================================

// ErrorResponse for error messages
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Details map[string]string `json:"details,omitempty"`
}

// SuccessResponse for success messages
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
