// backend/internal/gl-core/repository/journal_line_repository.go
package repository

import (
	"context"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JournalLineRepository struct {
	pool *pgxpool.Pool
}

// NewJournalLineRepository creates a new journal line repository
func NewJournalLineRepository(pool *pgxpool.Pool) *JournalLineRepository {
	return &JournalLineRepository{pool: pool}
}

// GetLinesByEntryID retrieves all lines for a journal entry
func (r *JournalLineRepository) GetLinesByEntryID(ctx context.Context, entryID uuid.UUID) ([]domain.JournalLine, error) {
	query := `
        SELECT id, account_id, line_number, reference, description, debit, credit
        FROM journal_lines
        WHERE journal_entry_id = $1
        ORDER BY line_number
    `

	rows, err := r.pool.Query(ctx, query, entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal lines: %w", err)
	}
	defer rows.Close()

	var lines []domain.JournalLine
	for rows.Next() {
		var line domain.JournalLine
		err := rows.Scan(
			&line.ID,
			&line.AccountID,
			&line.LineNumber,
			&line.Reference,
			&line.Description,
			&line.Debit,
			&line.Credit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan journal line: %w", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// GetLinesByAccountID retrieves all lines for a specific account
func (r *JournalLineRepository) GetLinesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.JournalLine, error) {
	query := `
        SELECT jl.id, jl.account_id, jl.line_number, jl.reference, jl.description, jl.debit, jl.credit
        FROM journal_lines jl
        INNER JOIN journal_entries je ON jl.journal_entry_id = je.id
        WHERE jl.account_id = $1 AND je.status = 'POSTED'
        ORDER BY je.transaction_date DESC, jl.line_number
        LIMIT $2 OFFSET $3
    `

	rows, err := r.pool.Query(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal lines by account: %w", err)
	}
	defer rows.Close()

	var lines []domain.JournalLine
	for rows.Next() {
		var line domain.JournalLine
		err := rows.Scan(
			&line.ID,
			&line.AccountID,
			&line.LineNumber,
			&line.Reference,
			&line.Description,
			&line.Debit,
			&line.Credit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan journal line: %w", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// GetLinesByReference retrieves lines by reference code
func (r *JournalLineRepository) GetLinesByReference(ctx context.Context, orgID uuid.UUID, reference string) ([]domain.JournalLine, error) {
	query := `
        SELECT jl.id, jl.account_id, jl.line_number, jl.reference, jl.description, jl.debit, jl.credit
        FROM journal_lines jl
        INNER JOIN journal_entries je ON jl.journal_entry_id = je.id
        WHERE je.organization_id = $1 AND jl.reference = $2 AND je.status = 'POSTED'
        ORDER BY je.transaction_date DESC, jl.line_number
    `

	rows, err := r.pool.Query(ctx, query, orgID, reference)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal lines by reference: %w", err)
	}
	defer rows.Close()

	var lines []domain.JournalLine
	for rows.Next() {
		var line domain.JournalLine
		err := rows.Scan(
			&line.ID,
			&line.AccountID,
			&line.LineNumber,
			&line.Reference,
			&line.Description,
			&line.Debit,
			&line.Credit,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan journal line: %w", err)
		}
		lines = append(lines, line)
	}

	return lines, nil
}

// CountLinesByAccount counts lines for an account
func (r *JournalLineRepository) CountLinesByAccount(ctx context.Context, accountID uuid.UUID) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM journal_lines jl
        INNER JOIN journal_entries je ON jl.journal_entry_id = je.id
        WHERE jl.account_id = $1 AND je.status = 'POSTED'
    `

	var count int
	err := r.pool.QueryRow(ctx, query, accountID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count journal lines: %w", err)
	}

	return count, nil
}
