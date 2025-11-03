// backend/internal/gl-core/service/journal_entry_service_interface.go
package service

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
)

// JournalEntryServiceInterface defines business logic for journal entries
type JournalEntryServiceInterface interface {
	// CreateEntry creates a new journal entry in DRAFT status
	CreateEntry(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error)

	// UpdateEntry updates an existing draft entry
	UpdateEntry(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error)

	// GetEntry retrieves a journal entry by ID
	GetEntry(ctx context.Context, entryID uuid.UUID) (*domain.JournalEntry, error)

	// ListEntries lists journal entries for an organization
	ListEntries(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.JournalEntry, error)

	// ListByStatus lists entries by status
	ListByStatus(ctx context.Context, orgID uuid.UUID, status domain.EntryStatus, limit, offset int) ([]*domain.JournalEntry, error)

	// PostEntry posts a draft entry to the ledger
	PostEntry(ctx context.Context, entryID uuid.UUID, postedBy uuid.UUID) error

	// VoidEntry voids a posted entry
	VoidEntry(ctx context.Context, entryID uuid.UUID) error

	// ReverseEntry creates a reversal entry for a posted entry
	ReverseEntry(ctx context.Context, entryID uuid.UUID, reversedBy uuid.UUID) (*domain.JournalEntry, error)

	// DeleteEntry soft deletes a draft entry
	DeleteEntry(ctx context.Context, entryID uuid.UUID) error

	// ValidateEntry validates an entry before posting
	ValidateEntry(ctx context.Context, entryID uuid.UUID) (*domain.PostingValidationResult, error)

	// GenerateEntryNumber generates the next entry number
	GenerateEntryNumber(ctx context.Context, orgID uuid.UUID, date string) (string, error)
}
