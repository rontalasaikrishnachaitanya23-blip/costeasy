// backend/internal/gl-core/repository/journal_entry_repository.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JournalEntryRepository struct {
	pool *pgxpool.Pool
}

// NewJournalEntryRepository creates a new journal entry repository
func NewJournalEntryRepository(pool *pgxpool.Pool) *JournalEntryRepository {
	return &JournalEntryRepository{pool: pool}
}

// Create creates a new journal entry with its lines in a transaction
func (r *JournalEntryRepository) Create(ctx context.Context, entry *domain.JournalEntry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert journal entry header
	entryQuery := `
        INSERT INTO journal_entries (
            id, organization_id, entry_number, transaction_date, posting_date,
            reference, description, status, total_debit, total_credit,
            created_by, posted_by, reversed_by, reversal_of,
            created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
    `

	_, err = tx.Exec(ctx, entryQuery,
		entry.ID,
		entry.OrganizationID,
		entry.EntryNumber,
		entry.TransactionDate,
		entry.PostingDate,
		entry.Reference,
		entry.Description,
		entry.Status,
		entry.TotalDebit,
		entry.TotalCredit,
		entry.CreatedBy,
		entry.PostedBy,
		entry.ReversedBy,
		entry.ReversalOf,
		entry.CreatedAt,
		entry.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to insert journal entry: %w", err)
	}

	// Insert journal lines
	lineQuery := `
        INSERT INTO journal_lines (
            id, journal_entry_id, account_id, line_number,
            reference, description, debit, credit
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	for _, line := range entry.Lines {
		_, err = tx.Exec(ctx, lineQuery,
			line.ID,
			entry.ID,
			line.AccountID,
			line.LineNumber,
			line.Reference,
			line.Description,
			line.Debit,
			line.Credit,
		)

		if err != nil {
			return fmt.Errorf("failed to insert journal line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Update updates an existing journal entry
func (r *JournalEntryRepository) Update(ctx context.Context, entry *domain.JournalEntry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update journal entry header
	entryQuery := `
        UPDATE journal_entries
        SET transaction_date = $2,
            posting_date = $3,
            reference = $4,
            description = $5,
            status = $6,
            total_debit = $7,
            total_credit = $8,
            posted_by = $9,
            reversed_by = $10,
            updated_at = $11
        WHERE id = $1
    `

	entry.UpdatedAt = time.Now()

	_, err = tx.Exec(ctx, entryQuery,
		entry.ID,
		entry.TransactionDate,
		entry.PostingDate,
		entry.Reference,
		entry.Description,
		entry.Status,
		entry.TotalDebit,
		entry.TotalCredit,
		entry.PostedBy,
		entry.ReversedBy,
		entry.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update journal entry: %w", err)
	}

	// Delete existing lines and re-insert (simpler than update logic)
	_, err = tx.Exec(ctx, "DELETE FROM journal_lines WHERE journal_entry_id = $1", entry.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing lines: %w", err)
	}

	// Re-insert lines
	lineQuery := `
        INSERT INTO journal_lines (
            id, journal_entry_id, account_id, line_number,
            reference, description, debit, credit
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	for _, line := range entry.Lines {
		_, err = tx.Exec(ctx, lineQuery,
			line.ID,
			entry.ID,
			line.AccountID,
			line.LineNumber,
			line.Reference,
			line.Description,
			line.Debit,
			line.Credit,
		)

		if err != nil {
			return fmt.Errorf("failed to insert journal line: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Delete soft deletes a journal entry (marks as void)
func (r *JournalEntryRepository) Delete(ctx context.Context, entryID uuid.UUID) error {
	query := `
        UPDATE journal_entries
        SET status = 'VOID', updated_at = NOW()
        WHERE id = $1
    `

	_, err := r.pool.Exec(ctx, query, entryID)
	if err != nil {
		return fmt.Errorf("failed to delete journal entry: %w", err)
	}

	return nil
}

// GetByID retrieves a journal entry by ID with all its lines
func (r *JournalEntryRepository) GetByID(ctx context.Context, entryID uuid.UUID) (*domain.JournalEntry, error) {
	// Get entry header
	entryQuery := `
        SELECT id, organization_id, entry_number, transaction_date, posting_date,
               reference, description, status, total_debit, total_credit,
               created_by, posted_by, reversed_by, reversal_of,
               created_at, updated_at
        FROM journal_entries
        WHERE id = $1
    `

	entry := &domain.JournalEntry{}
	var postingDate *time.Time
	var postedBy, reversedBy, reversalOf *uuid.UUID // ✅ FIXED: Changed from *time.Time to *uuid.UUID

	err := r.pool.QueryRow(ctx, entryQuery, entryID).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.EntryNumber,
		&entry.TransactionDate,
		&postingDate,
		&entry.Reference,
		&entry.Description,
		&entry.Status,
		&entry.TotalDebit,
		&entry.TotalCredit,
		&entry.CreatedBy,
		&postedBy,
		&reversedBy,
		&reversalOf,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}

	// Assign nullable fields
	entry.PostingDate = postingDate
	entry.PostedBy = postedBy     // ✅ FIXED: Direct assignment (already *uuid.UUID)
	entry.ReversedBy = reversedBy // ✅ FIXED: Direct assignment (already *uuid.UUID)
	entry.ReversalOf = reversalOf // ✅ FIXED: Direct assignment (already *uuid.UUID)

	// Get entry lines
	linesQuery := `
        SELECT id, account_id, line_number, reference, description, debit, credit
        FROM journal_lines
        WHERE journal_entry_id = $1
        ORDER BY line_number
    `

	rows, err := r.pool.Query(ctx, linesQuery, entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get journal lines: %w", err)
	}
	defer rows.Close()

	entry.Lines = []domain.JournalLine{}
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
		entry.Lines = append(entry.Lines, line)
	}

	return entry, nil
}

// GetByEntryNumber retrieves a journal entry by entry number
func (r *JournalEntryRepository) GetByEntryNumber(ctx context.Context, orgID uuid.UUID, entryNumber string) (*domain.JournalEntry, error) {
	query := `
        SELECT id
        FROM journal_entries
        WHERE organization_id = $1 AND entry_number = $2
    `

	var entryID uuid.UUID
	err := r.pool.QueryRow(ctx, query, orgID, entryNumber).Scan(&entryID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("journal entry not found")
		}
		return nil, fmt.Errorf("failed to get journal entry: %w", err)
	}

	return r.GetByID(ctx, entryID)
}

// ListByOrganization lists journal entries for an organization
func (r *JournalEntryRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.JournalEntry, error) {
	query := `
        SELECT id
        FROM journal_entries
        WHERE organization_id = $1
        ORDER BY transaction_date DESC, entry_number DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.pool.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list journal entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.JournalEntry
	for rows.Next() {
		var entryID uuid.UUID
		if err := rows.Scan(&entryID); err != nil {
			return nil, fmt.Errorf("failed to scan entry ID: %w", err)
		}

		entry, err := r.GetByID(ctx, entryID)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// ListByStatus lists journal entries by status
func (r *JournalEntryRepository) ListByStatus(ctx context.Context, orgID uuid.UUID, status domain.EntryStatus, limit, offset int) ([]*domain.JournalEntry, error) {
	query := `
        SELECT id
        FROM journal_entries
        WHERE organization_id = $1 AND status = $2
        ORDER BY transaction_date DESC
        LIMIT $3 OFFSET $4
    `

	rows, err := r.pool.Query(ctx, query, orgID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list journal entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.JournalEntry
	for rows.Next() {
		var entryID uuid.UUID
		if err := rows.Scan(&entryID); err != nil {
			return nil, fmt.Errorf("failed to scan entry ID: %w", err)
		}

		entry, err := r.GetByID(ctx, entryID)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// ListByDateRange lists journal entries within a date range
func (r *JournalEntryRepository) ListByDateRange(ctx context.Context, orgID uuid.UUID, startDate, endDate string, limit, offset int) ([]*domain.JournalEntry, error) {
	query := `
        SELECT id
        FROM journal_entries
        WHERE organization_id = $1
          AND transaction_date >= $2
          AND transaction_date <= $3
        ORDER BY transaction_date DESC
        LIMIT $4 OFFSET $5
    `

	rows, err := r.pool.Query(ctx, query, orgID, startDate, endDate, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list journal entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.JournalEntry
	for rows.Next() {
		var entryID uuid.UUID
		if err := rows.Scan(&entryID); err != nil {
			return nil, fmt.Errorf("failed to scan entry ID: %w", err)
		}

		entry, err := r.GetByID(ctx, entryID)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

// GetNextEntryNumber generates the next entry number for a given date
func (r *JournalEntryRepository) GetNextEntryNumber(ctx context.Context, orgID uuid.UUID, date string) (int, error) {
	query := `
        SELECT COUNT(*) + 1
        FROM journal_entries
        WHERE organization_id = $1
          AND entry_number LIKE $2
    `

	pattern := fmt.Sprintf("JE-%s-%%", date)

	var sequence int
	err := r.pool.QueryRow(ctx, query, orgID, pattern).Scan(&sequence)
	if err != nil {
		return 0, fmt.Errorf("failed to get next entry number: %w", err)
	}

	return sequence, nil
}

// CountByOrganization counts total entries for an organization
func (r *JournalEntryRepository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int, error) {
	query := `
        SELECT COUNT(*)
        FROM journal_entries
        WHERE organization_id = $1
    `

	var count int
	err := r.pool.QueryRow(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count journal entries: %w", err)
	}

	return count, nil
}
