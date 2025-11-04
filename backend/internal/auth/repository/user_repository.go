// backend/internal/auth/repository/user_repository.go
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository implements UserRepositoryInterface using pgxpool
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id, organization_id, username, email, password_hash,
			first_name, last_name, phone,
			mfa_enabled, mfa_secret, mfa_backup_codes,
			is_active, is_verified, password_changed_at,
			allow_remote_access, remote_access_reason,
			remote_access_approved_by, remote_access_approved_at,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11,
			$12, $13, $14,
			$15, $16,
			$17, $18,
			$19, $20
		)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.OrganizationID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.MFAEnabled,
		user.MFASecret,
		user.MFABackupCodes,
		user.IsActive,
		user.IsVerified,
		user.PasswordChangedAt,
		user.AllowRemoteAccess,
		user.RemoteAccessReason,
		user.RemoteAccessApprovedBy,
		user.RemoteAccessApprovedAt,
		user.CreatedBy,
		user.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT 
			id, organization_id, username, email, password_hash,
			first_name, last_name, phone,
			mfa_enabled, mfa_secret, mfa_backup_codes,
			is_active, is_verified, email_verified_at,
			last_login_at, last_login_ip,
			password_changed_at, failed_login_attempts, locked_until,
			allow_remote_access, remote_access_reason, 
			remote_access_approved_by, remote_access_approved_at,
			created_at, updated_at, created_by, updated_by
		FROM users
		WHERE id = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.MFABackupCodes,
		&user.IsActive,
		&user.IsVerified,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.LastLoginIP,
		&user.PasswordChangedAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.AllowRemoteAccess,
		&user.RemoteAccessReason,
		&user.RemoteAccessApprovedBy,
		&user.RemoteAccessApprovedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT 
			id, organization_id, username, email, password_hash,
			first_name, last_name, phone,
			mfa_enabled, mfa_secret, mfa_backup_codes,
			is_active, is_verified, email_verified_at,
			last_login_at, last_login_ip,
			password_changed_at, failed_login_attempts, locked_until,
			allow_remote_access, remote_access_reason, 
			remote_access_approved_by, remote_access_approved_at,
			created_at, updated_at, created_by, updated_by
		FROM users
		WHERE username = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.MFABackupCodes,
		&user.IsActive,
		&user.IsVerified,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.LastLoginIP,
		&user.PasswordChangedAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.AllowRemoteAccess,
		&user.RemoteAccessReason,
		&user.RemoteAccessApprovedBy,
		&user.RemoteAccessApprovedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT 
			id, organization_id, username, email, password_hash,
			first_name, last_name, phone,
			mfa_enabled, mfa_secret, mfa_backup_codes,
			is_active, is_verified, email_verified_at,
			last_login_at, last_login_ip,
			password_changed_at, failed_login_attempts, locked_until,
			allow_remote_access, remote_access_reason, 
			remote_access_approved_by, remote_access_approved_at,
			created_at, updated_at, created_by, updated_by
		FROM users
		WHERE email = $1
	`

	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.MFAEnabled,
		&user.MFASecret,
		&user.MFABackupCodes,
		&user.IsActive,
		&user.IsVerified,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.LastLoginIP,
		&user.PasswordChangedAt,
		&user.FailedLoginAttempts,
		&user.LockedUntil,
		&user.AllowRemoteAccess,
		&user.RemoteAccessReason,
		&user.RemoteAccessApprovedBy,
		&user.RemoteAccessApprovedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users SET
			username = $2,
			email = $3,
			first_name = $4,
			last_name = $5,
			phone = $6,
			is_active = $7,
			is_verified = $8,
			allow_remote_access = $9,
			remote_access_reason = $10,
			remote_access_approved_by = $11,
			remote_access_approved_at = $12,
			updated_by = $13,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
		user.IsVerified,
		user.AllowRemoteAccess,
		user.RemoteAccessReason,
		user.RemoteAccessApprovedBy,
		user.RemoteAccessApprovedAt,
		user.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete soft deletes a user (sets is_active = false)
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// HardDelete permanently deletes a user
func (r *UserRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT 
			id, organization_id, username, email, password_hash,
			first_name, last_name, phone,
			mfa_enabled, is_active, is_verified,
			last_login_at, allow_remote_access,
			created_at, updated_at
		FROM users
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.OrganizationID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.MFAEnabled,
			&user.IsActive,
			&user.IsVerified,
			&user.LastLoginAt,
			&user.AllowRemoteAccess,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

// Count returns total number of users in organization
func (r *UserRepository) Count(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE organization_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// GetWithRoles retrieves a user with their roles
func (r *UserRepository) GetWithRoles(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// Get user
	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get roles
	roles, err := r.GetRoles(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Roles = roles
	return user, nil
}

// UpdateLastLogin updates user's last login timestamp and IP
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, ipAddress string) error {
	query := `
		UPDATE users 
		SET last_login_at = CURRENT_TIMESTAMP, last_login_ip = $2
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, ipAddress)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates user's password hash
func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `
		UPDATE users 
		SET password_hash = $2, password_changed_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// IncrementFailedAttempts increments failed login attempts
func (r *UserRepository) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users 
		SET failed_login_attempts = failed_login_attempts + 1,
		    locked_until = CASE 
		        WHEN failed_login_attempts + 1 >= 5 
		        THEN CURRENT_TIMESTAMP + INTERVAL '30 minutes'
		        ELSE locked_until
		    END
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment failed attempts: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// ResetFailedAttempts resets failed login attempts to 0
func (r *UserRepository) ResetFailedAttempts(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users 
		SET failed_login_attempts = 0, locked_until = NULL
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to reset failed attempts: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// LockAccount locks user account until specified time
func (r *UserRepository) LockAccount(ctx context.Context, id uuid.UUID, lockedUntil *time.Time) error {
	query := `UPDATE users SET locked_until = $2 WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, lockedUntil)
	if err != nil {
		return fmt.Errorf("failed to lock account: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// UnlockAccount unlocks user account
func (r *UserRepository) UnlockAccount(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET locked_until = NULL, failed_login_attempts = 0 WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to unlock account: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// EnableMFA enables MFA for user
func (r *UserRepository) EnableMFA(ctx context.Context, id uuid.UUID, secret string, backupCodes []string) error {
	query := `
		UPDATE users 
		SET mfa_enabled = true, mfa_secret = $2, mfa_backup_codes = $3
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id, secret, backupCodes)
	if err != nil {
		return fmt.Errorf("failed to enable MFA: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// DisableMFA disables MFA for user
func (r *UserRepository) DisableMFA(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users 
		SET mfa_enabled = false, mfa_secret = NULL, mfa_backup_codes = NULL
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to disable MFA: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// AssignRole assigns a role to user
func (r *UserRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID, assignedBy *uuid.UUID) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.Exec(ctx, query, userID, roleID, assignedBy)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RemoveRole removes a role from user
func (r *UserRepository) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	result, err := r.db.Exec(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("role assignment not found")
	}

	return nil
}

// GetRoles gets all roles for a user
func (r *UserRepository) GetRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.organization_id, r.name, r.display_name, 
		       r.description, r.is_system_role, r.is_active,
		       r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.is_active = true
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
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
