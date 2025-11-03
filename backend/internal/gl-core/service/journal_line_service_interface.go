// backend/internal/gl-core/service/journal_line_service_interface.go
package service

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
)

// JournalLineServiceInterface defines business logic for querying journal lines
type JournalLineServiceInterface interface {
	// GetLinesByEntry retrieves all lines for a journal entry
	GetLinesByEntry(ctx context.Context, entryID uuid.UUID) ([]domain.JournalLine, error)

	// GetLinesByAccount retrieves all posted lines for an account (for account statement)
	GetLinesByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.JournalLine, error)

	// GetLinesByReference retrieves lines by reference code
	GetLinesByReference(ctx context.Context, orgID uuid.UUID, reference string) ([]domain.JournalLine, error)

	// GetAccountBalance calculates running balance for an account
	GetAccountBalance(ctx context.Context, accountID uuid.UUID) (float64, error)
}
