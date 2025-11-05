package repository

import (
	"context"
	"errors"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (
            id, email, password_hash, username, first_name, last_name, phone,
            is_active, is_verified, organization_id,
            mfa_enabled, mfa_method, totp_secret, backup_codes,
            created_at, updated_at, created_by
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
        )
    `

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
		user.IsVerified,
		user.OrganizationID,
		user.MFAEnabled,
		user.MFAMethod,
		user.TOTPSecret,
		pq.Array(user.BackupCodes),
		user.CreatedAt,
		user.UpdatedAt,
		user.CreatedBy,
	)

	return err
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
        SELECT id, email, password_hash,username, first_name, last_name, phone,
               is_active, is_verified, last_login_at, organization_id,
               mfa_enabled, mfa_method, totp_secret, backup_codes, mfa_verified_at,
               created_at, updated_at, created_by, updated_by
        FROM users
        WHERE id = $1
    `

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.IsVerified,
		&user.LastLoginAt,
		&user.OrganizationID,
		&user.MFAEnabled,
		&user.MFAMethod,
		&user.TOTPSecret,
		pq.Array(&user.BackupCodes),
		&user.MFAVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	// If you don't have username column, use email instead or add username column
	return r.GetByEmail(ctx, username)
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
        SELECT id, email, password_hash,username, first_name, last_name, phone,
               is_active, is_verified, last_login_at, organization_id,
               mfa_enabled, mfa_method, totp_secret, backup_codes, mfa_verified_at,
               created_at, updated_at, created_by, updated_by
        FROM users
        WHERE email = $1
    `

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.IsActive,
		&user.IsVerified,
		&user.LastLoginAt,
		&user.OrganizationID,
		&user.MFAEnabled,
		&user.MFAMethod,
		&user.TOTPSecret,
		pq.Array(&user.BackupCodes),
		&user.MFAVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
        UPDATE users SET
            email = $2,
            password_hash = $3,
			username=$4,
            first_name = $5,
            last_name = $6,
            phone = $7,
            is_active = $8,
            is_verified = $9,
            last_login_at = $10,
            mfa_enabled = $11,
            mfa_method = $12,
            totp_secret = $13,
            backup_codes = $14,
            mfa_verified_at = $15,
            updated_at = $16,
            updated_by = $17
        WHERE id = $1
    `

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.IsActive,
		user.IsVerified,
		user.LastLoginAt,
		user.MFAEnabled,
		user.MFAMethod,
		user.TOTPSecret,
		pq.Array(user.BackupCodes),
		user.MFAVerifiedAt,
		user.UpdatedAt,
		user.UpdatedBy,
	)

	return err
}

// Delete soft deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// HardDelete permanently deletes a user
func (r *UserRepository) HardDelete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// List retrieves users with pagination
func (r *UserRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.User, error) {
	query := `
        SELECT id, email, password_hash, username,first_name, last_name, phone,
               is_active, is_verified, last_login_at, organization_id,
               mfa_enabled, mfa_method, totp_secret, backup_codes, mfa_verified_at,
               created_at, updated_at, created_by, updated_by
        FROM users
        WHERE organization_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.Username,
			&user.FirstName, &user.LastName, &user.Phone,
			&user.IsActive, &user.IsVerified, &user.LastLoginAt, &user.OrganizationID,
			&user.MFAEnabled, &user.MFAMethod, &user.TOTPSecret,
			pq.Array(&user.BackupCodes), &user.MFAVerifiedAt,
			&user.CreatedAt, &user.UpdatedAt, &user.CreatedBy, &user.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// Count returns total number of users in organization
func (r *UserRepository) Count(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE organization_id = $1`
	var count int
	err := r.db.QueryRow(ctx, query, orgID).Scan(&count)
	return count, err
}

// GetWithRoles retrieves a user with their roles
func (r *UserRepository) GetWithRoles(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// For now, just get the user. Roles can be fetched separately
	return r.GetByID(ctx, id)
}

// UpdateLastLogin updates user's last login timestamp and IP
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, ipAddress string) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdatePassword updates user's password hash
func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, passwordHash)
	return err
}

// IncrementFailedAttempts increments failed login attempts
func (r *UserRepository) IncrementFailedAttempts(ctx context.Context, id uuid.UUID) error {
	// If you don't have failed_attempts column, you can skip this or add the column
	// query := `UPDATE users SET failed_attempts = failed_attempts + 1 WHERE id = $1`
	// _, err := r.db.Exec(ctx, query, id)
	// return err
	return nil // Placeholder
}

// ResetFailedAttempts resets failed login attempts to 0
func (r *UserRepository) ResetFailedAttempts(ctx context.Context, id uuid.UUID) error {
	// If you don't have failed_attempts column, you can skip this
	// query := `UPDATE users SET failed_attempts = 0 WHERE id = $1`
	// _, err := r.db.Exec(ctx, query, id)
	// return err
	return nil // Placeholder
}

// LockAccount locks user account until specified time
func (r *UserRepository) LockAccount(ctx context.Context, id uuid.UUID, lockedUntil *time.Time) error {
	// If you don't have locked_until column, you can skip this
	// query := `UPDATE users SET locked_until = $2 WHERE id = $1`
	// _, err := r.db.Exec(ctx, query, id, lockedUntil)
	// return err
	return nil // Placeholder
}

// UnlockAccount unlocks user account
func (r *UserRepository) UnlockAccount(ctx context.Context, id uuid.UUID) error {
	// If you don't have locked_until column, you can skip this
	// query := `UPDATE users SET locked_until = NULL WHERE id = $1`
	// _, err := r.db.Exec(ctx, query, id)
	// return err
	return nil // Placeholder
}

// EnableMFA enables MFA for user
func (r *UserRepository) EnableMFA(ctx context.Context, id uuid.UUID, secret string, backupCodes []string) error {
	query := `
        UPDATE users SET
            mfa_enabled = true,
            totp_secret = $2,
            backup_codes = $3,
            mfa_verified_at = NOW(),
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.Exec(ctx, query, id, secret, pq.Array(backupCodes))
	return err
}

// DisableMFA disables MFA for user
func (r *UserRepository) DisableMFA(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE users SET
            mfa_enabled = false,
            mfa_method = 'none',
            totp_secret = NULL,
            backup_codes = NULL,
            mfa_verified_at = NULL,
            updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// AssignRole assigns a role to a user
func (r *UserRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `
        INSERT INTO user_roles (user_id, role_id, assigned_at)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id, role_id) DO NOTHING
    `
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

// RemoveRole removes a role from user
func (r *UserRepository) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

// GetRoles gets all roles for a user
func (r *UserRepository) GetRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	query := `
        SELECT r.id, r.organization_id, r.name, r.description, r.is_system_role,
               r.is_active, r.created_at, r.updated_at, r.created_by, r.updated_by
        FROM roles r
        INNER JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1 AND r.is_active = true
    `

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID, &role.OrganizationID, &role.Name, &role.Description,
			&role.IsSystemRole, &role.IsActive, &role.CreatedAt, &role.UpdatedAt,
			// &role.CreatedBy, &role.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetUserRole gets the current role for a user (returns nil if no role)
func (r *UserRepository) GetUserRole(ctx context.Context, userID uuid.UUID) (*domain.Role, error) {
	query := `
        SELECT r.id, r.organization_id, r.name, r.display_name, 
               r.description, r.is_system_role, r.is_active,
               r.created_at, r.updated_at
        FROM roles r
        INNER JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1
        LIMIT 1
    `

	var role domain.Role
	err := r.db.QueryRow(ctx, query, userID).Scan(
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
			return nil, nil // No role assigned
		}
		return nil, err
	}

	return &role, nil
}
