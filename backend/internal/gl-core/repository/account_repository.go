package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	//"github.com/chaitu35/costeasy/backend/database"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GLAccountSearchParams struct {
	Name        string
	Type        *domain.AccountType
	IsActive    *bool
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

var _ GLAccountRepositoryInterface = (*GLAccountRepository)(nil)

type GLAccountRepository struct {
	pool *pgxpool.Pool
}

func NewGLAccountRepository(pool *pgxpool.Pool) *GLAccountRepository {
	return &GLAccountRepository{
		pool: pool,
	}
}

func (r *GLAccountRepository) CreateGLAccount(ctx context.Context, GLAccount domain.GLAccount) (domain.GLAccount, error) {
	prefix, err := domain.PrefixForType(GLAccount.Type)
	if err != nil {
		return domain.GLAccount{}, err
	}
	if !strings.HasPrefix(GLAccount.Code, prefix) {
		GLAccount.Code = prefix + GLAccount.Code
	}

	query := `
        INSERT INTO gl_accounts (code, name, type, parent_code, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `
	var newGLAccount domain.GLAccount
	err = r.pool.QueryRow(ctx, query,
		GLAccount.Code,
		GLAccount.Name,
		GLAccount.Type,
		GLAccount.ParentCode,
		true,
	).Scan(&newGLAccount.ID, &newGLAccount.CreateAt, &newGLAccount.UpdateAt)
	if err != nil {
		return domain.GLAccount{}, err
	}

	newGLAccount.Code = GLAccount.Code
	newGLAccount.Name = GLAccount.Name
	newGLAccount.Type = GLAccount.Type
	newGLAccount.ParentCode = GLAccount.ParentCode
	newGLAccount.IsActive = true

	return newGLAccount, nil
}

func (r *GLAccountRepository) GetGLAccountByID(ctx context.Context, id uuid.UUID, includeInactive bool) (domain.GLAccount, error) {
	query := `
        SELECT id, code, name, type, parent_code, is_active, created_at, updated_at 
        FROM gl_accounts WHERE id = $1
    `
	if !includeInactive {
		query += " AND is_active = TRUE"
	}

	var GLAccount domain.GLAccount
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&GLAccount.ID,
		&GLAccount.Code,
		&GLAccount.Name,
		&GLAccount.Type,
		&GLAccount.ParentCode,
		&GLAccount.IsActive,
		&GLAccount.CreateAt,
		&GLAccount.UpdateAt,
	)
	if err != nil {
		return domain.GLAccount{}, err
	}
	return GLAccount, nil
}

func (r *GLAccountRepository) GetGLAccountByCode(ctx context.Context, code string, includeInactive bool) (domain.GLAccount, error) {
	query := `
		SELECT id, code, name, type, parent_code, is_active, created_at, updated_at
		FROM gl_accounts WHERE code = $1
	`
	if !includeInactive {
		query += " AND is_active = TRUE"
	}

	var GLAccount domain.GLAccount
	err := r.pool.QueryRow(ctx, query, code).Scan(

		&GLAccount.ID,
		&GLAccount.Code,
		&GLAccount.Name,
		&GLAccount.Type,
		&GLAccount.ParentCode,
		&GLAccount.IsActive,
		&GLAccount.CreateAt,
		&GLAccount.UpdateAt,
	)
	if err != nil {
		return domain.GLAccount{}, err
	}
	return GLAccount, nil
}

func (r *GLAccountRepository) DeactivateGLAccount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE gl_accounts SET is_active = FALSE, updated_at = NOW()	 WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *GLAccountRepository) ActivateGLAccount(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE gl_accounts SET is_active = TRUE, updated_at = NOW() WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *GLAccountRepository) UpdateGLAccount(ctx context.Context, GLAccount domain.GLAccount) (domain.GLAccount, error) {
	query := `
		UPDATE gl_accounts 
		SET name = $1, parent_code = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING code, type, is_active, created_at, updated_at
	`
	var updatedGLAccount domain.GLAccount
	err := r.pool.QueryRow(ctx, query,
		GLAccount.Name,
		GLAccount.ParentCode,
		GLAccount.ID,
	).Scan(
		&updatedGLAccount.Code,
		&updatedGLAccount.Type,
		&updatedGLAccount.IsActive,
		&updatedGLAccount.CreateAt,
		&updatedGLAccount.UpdateAt,
	)
	if err != nil {
		return domain.GLAccount{}, err
	}

	updatedGLAccount.ID = GLAccount.ID
	updatedGLAccount.Name = GLAccount.Name
	updatedGLAccount.ParentCode = GLAccount.ParentCode

	return updatedGLAccount, nil
}

func (r *GLAccountRepository) ListGLAccounts(ctx context.Context, includeInactive bool) ([]domain.GLAccount, error) {
	query := `
		SELECT id, code, name, type, parent_code, is_active, created_at, updated_at
		FROM gl_accounts
	`
	if !includeInactive {
		query += " WHERE is_active = TRUE"
	}

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var GLAccounts []domain.GLAccount
	for rows.Next() {
		var GLAccount domain.GLAccount
		err := rows.Scan(
			&GLAccount.ID,
			&GLAccount.Code,
			&GLAccount.Name,
			&GLAccount.Type,
			&GLAccount.ParentCode,
			&GLAccount.IsActive,
			&GLAccount.CreateAt,
			&GLAccount.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		GLAccounts = append(GLAccounts, GLAccount)
	}

	return GLAccounts, nil
}

func (r *GLAccountRepository) SearchGLAccounts(ctx context.Context, params GLAccountSearchParams) ([]domain.GLAccount, error) {
	baseQuery := `
        SELECT id, code, name, type, parent_code, is_active, created_at, updated_at
        FROM gl_accounts WHERE 1=1
    `
	args := []interface{}{}
	argIdx := 1

	// Build conditions dynamically
	if params.Name != "" {
		baseQuery += fmt.Sprintf(" AND LOWER(name) LIKE LOWER($%d)", argIdx)
		args = append(args, "%"+params.Name+"%")
		argIdx++
	}

	if params.Type != nil {
		baseQuery += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, *params.Type)
		argIdx++
	}

	if params.IsActive != nil {
		baseQuery += fmt.Sprintf(" AND is_active = $%d", argIdx)
		args = append(args, *params.IsActive)
		argIdx++
	}

	if params.CreatedFrom != nil {
		baseQuery += fmt.Sprintf(" AND created_at >= $%d", argIdx)
		args = append(args, *params.CreatedFrom)
		argIdx++
	}
	if params.CreatedTo != nil {
		baseQuery += fmt.Sprintf(" AND created_at <= $%d", argIdx)
		args = append(args, *params.CreatedTo)
		argIdx++
	}

	rows, err := r.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var GLAccounts []domain.GLAccount
	for rows.Next() {
		var GLAccount domain.GLAccount
		err := rows.Scan(
			&GLAccount.ID,
			&GLAccount.Code,
			&GLAccount.Name,
			&GLAccount.Type,
			&GLAccount.ParentCode,
			&GLAccount.IsActive,
			&GLAccount.CreateAt,
			&GLAccount.UpdateAt,
		)
		if err != nil {
			return nil, err
		}
		GLAccounts = append(GLAccounts, GLAccount)
	}
	return GLAccounts, rows.Err()
}

func (r *GLAccountRepository) SoftDeleteGLAccount(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE gl_accounts SET is_active = FALSE, updated_at = NOW()
        WHERE id = $1
    `
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// In repository/GLAccount_repository.go

func (r *GLAccountRepository) HasActiveTransactions(ctx context.Context, GLAccountID uuid.UUID) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM journal_lines 
        WHERE account_id = $1
    `
	err := r.pool.QueryRow(ctx, query, GLAccountID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check transactions: %w", err)
	}
	return count > 0, nil
}

func (r *GLAccountRepository) HasChildGLAccounts(ctx context.Context, GLAccountID uuid.UUID) (bool, error) {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM gl_accounts 
        WHERE parent_code = (SELECT code FROM gl_accounts WHERE id = $1)
    `
	err := r.pool.QueryRow(ctx, query, GLAccountID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check child GLAccounts: %w", err)
	}
	return count > 0, nil
}
