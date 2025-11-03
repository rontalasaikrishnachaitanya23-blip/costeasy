// backend/internal/gl-core/repository/journal_entry_repository_interface.go
package repository

import (
    "context"

    "github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
    "github.com/google/uuid"
)

// JournalEntryRepositoryInterface defines the data access layer for journal entries
type JournalEntryRepositoryInterface interface {
    // Create creates a new journal entry with its lines
    Create(ctx context.Context, entry *domain.JournalEntry) error

    // Update updates an existing journal entry
    Update(ctx context.Context, entry *domain.JournalEntry) error

    // Delete soft deletes a journal entry
    Delete(ctx context.Context, entryID uuid.UUID) error

    // GetByID retrieves a journal entry by ID with all its lines
    GetByID(ctx context.Context, entryID uuid.UUID) (*domain.JournalEntry, error)

    // GetByEntryNumber retrieves a journal entry by entry number
    GetByEntryNumber(ctx context.Context, orgID uuid.UUID, entryNumber string) (*domain.JournalEntry, error)

    // ListByOrganization lists journal entries for an organization with pagination
    ListByOrganization(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.JournalEntry, error)

    // ListByStatus lists journal entries by status
    ListByStatus(ctx context.Context, orgID uuid.UUID, status domain.EntryStatus, limit, offset int) ([]*domain.JournalEntry, error)

    // ListByDateRange lists journal entries within a date range
    ListByDateRange(ctx context.Context, orgID uuid.UUID, startDate, endDate string, limit, offset int) ([]*domain.JournalEntry, error)

    // GetNextEntryNumber generates the next entry number for a given date
    GetNextEntryNumber(ctx context.Context, orgID uuid.UUID, date string) (int, error)

    // CountByOrganization counts total entries for an organization
    CountByOrganization(ctx context.Context, orgID uuid.UUID) (int, error)
}
