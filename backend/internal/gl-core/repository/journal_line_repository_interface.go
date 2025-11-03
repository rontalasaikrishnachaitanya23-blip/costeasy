// backend/internal/gl-core/repository/journal_line_repository_interface.go
package repository

import (
    "context"

    "github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
    "github.com/google/uuid"
)

// JournalLineRepositoryInterface defines data access for journal lines
type JournalLineRepositoryInterface interface {
    // GetLinesByEntryID retrieves all lines for a journal entry
    GetLinesByEntryID(ctx context.Context, entryID uuid.UUID) ([]domain.JournalLine, error)

    // GetLinesByAccountID retrieves all lines for a specific account
    GetLinesByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.JournalLine, error)

    // GetLinesByReference retrieves lines by reference code
    GetLinesByReference(ctx context.Context, orgID uuid.UUID, reference string) ([]domain.JournalLine, error)

    // CountLinesByAccount counts lines for an account
    CountLinesByAccount(ctx context.Context, accountID uuid.UUID) (int, error)
}
