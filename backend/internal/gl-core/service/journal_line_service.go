// backend/internal/gl-core/service/journal_line_service.go
package service

import (
	"context"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/repository"
	"github.com/google/uuid"
)

type JournalLineService struct {
	lineRepo    repository.JournalLineRepositoryInterface
	accountRepo repository.GLAccountRepositoryInterface
}

// NewJournalLineService creates a new journal line service
func NewJournalLineService(
	lineRepo repository.JournalLineRepositoryInterface,
	accountRepo repository.GLAccountRepositoryInterface,
) *JournalLineService {
	return &JournalLineService{
		lineRepo:    lineRepo,
		accountRepo: accountRepo,
	}
}

// GetLinesByEntry retrieves all lines for a journal entry
func (s *JournalLineService) GetLinesByEntry(ctx context.Context, entryID uuid.UUID) ([]domain.JournalLine, error) {
	lines, err := s.lineRepo.GetLinesByEntryID(ctx, entryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lines for entry: %w", err)
	}

	return lines, nil
}

// GetLinesByAccount retrieves all posted lines for an account
func (s *JournalLineService) GetLinesByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.JournalLine, error) {
	if limit <= 0 {
		limit = 100
	}

	// Verify account exists
	_, err := s.accountRepo.GetGLAccountByID(ctx, accountID, false)
	if err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}

	lines, err := s.lineRepo.GetLinesByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get lines for account: %w", err)
	}

	return lines, nil
}

// GetLinesByReference retrieves lines by reference code
func (s *JournalLineService) GetLinesByReference(ctx context.Context, orgID uuid.UUID, reference string) ([]domain.JournalLine, error) {
	if reference == "" {
		return nil, fmt.Errorf("reference is required")
	}

	lines, err := s.lineRepo.GetLinesByReference(ctx, orgID, reference)
	if err != nil {
		return nil, fmt.Errorf("failed to get lines by reference: %w", err)
	}

	return lines, nil
}

// GetAccountBalance calculates running balance for an account
func (s *JournalLineService) GetAccountBalance(ctx context.Context, accountID uuid.UUID) (float64, error) {
	// Get account to determine type
	account, err := s.accountRepo.GetGLAccountByID(ctx, accountID, false)
	if err != nil {
		return 0, fmt.Errorf("account not found: %w", err)
	}

	// Get all posted lines for this account
	lines, err := s.lineRepo.GetLinesByAccountID(ctx, accountID, 10000, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to get account lines: %w", err)
	}

	// Calculate balance based on account type
	balance := 0.0

	for _, line := range lines {
		switch account.Type {
		case domain.AccountTypeAsset, domain.AccountTypeExpense:
			// Debit increases, credit decreases
			balance += line.Debit - line.Credit

		case domain.AccountTypeLiability, domain.AccountTypeEquity, domain.AccountTypeRevenue:
			// Credit increases, debit decreases
			balance += line.Credit - line.Debit
		}
	}

	return balance, nil
}
