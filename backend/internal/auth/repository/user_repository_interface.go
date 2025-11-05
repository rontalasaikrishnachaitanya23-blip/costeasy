// backend/internal/auth/repository/user_repository_interface.go
package repository

import (
    "context"
    "time"
    
    "github.com/chaitu35/costeasy/backend/internal/auth/domain"
    "github.com/google/uuid"
)

// UserRepositoryInterface defines the contract for user data operations
type UserRepositoryInterface interface {
    // Create creates a new user
    Create(ctx context.Context, user *domain.User) error
    
    // GetByID retrieves a user by ID
    GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
    
    // GetByUsername retrieves a user by username
    GetByUsername(ctx context.Context, username string) (*domain.User, error)
    
    // GetByEmail retrieves a user by email
    GetByEmail(ctx context.Context, email string) (*domain.User, error)
    
    // Update updates an existing user
    Update(ctx context.Context, user *domain.User) error
    
    // Delete soft deletes a user (sets is_active = false)
    Delete(ctx context.Context, id uuid.UUID) error
    
    // HardDelete permanently deletes a user
    HardDelete(ctx context.Context, id uuid.UUID) error
    
    // List retrieves users with pagination
    List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.User, error)
    
    // Count returns total number of users in organization
    Count(ctx context.Context, orgID uuid.UUID) (int, error)
    
    // GetWithRoles retrieves a user with their roles
    GetWithRoles(ctx context.Context, id uuid.UUID) (*domain.User, error)
    
    // UpdateLastLogin updates user's last login timestamp and IP
    UpdateLastLogin(ctx context.Context, id uuid.UUID, ipAddress string) error
    
    // UpdatePassword updates user's password hash
    UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
    
    // IncrementFailedAttempts increments failed login attempts
    IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error
    
    // ResetFailedAttempts resets failed login attempts to 0
    ResetFailedAttempts(ctx context.Context, id uuid.UUID) error
    
    // LockAccount locks user account until specified time
    LockAccount(ctx context.Context, id uuid.UUID, lockedUntil *time.Time) error
    
    // UnlockAccount unlocks user account
    UnlockAccount(ctx context.Context, id uuid.UUID) error
    
    // EnableMFA enables MFA for user
    EnableMFA(ctx context.Context, id uuid.UUID, secret string, backupCodes []string) error
    
    // DisableMFA disables MFA for user
    DisableMFA(ctx context.Context, id uuid.UUID) error
    
    // AssignRole assigns a role to user
    AssignRole(ctx context.Context, userID, roleID uuid.UUID) error
    
    // RemoveRole removes a role from user
    RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error
    
    // GetRoles gets all roles for a user (returns multiple roles)
    GetRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, error)
    
    // GetUserRole gets the primary/current role for a user (returns single role or nil)
    GetUserRole(ctx context.Context, userID uuid.UUID) (*domain.Role, error)  // CHANGED: returns single role, not slice
}