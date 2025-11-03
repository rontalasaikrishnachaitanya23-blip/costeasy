// backend/internal/gl-core/service/journal_entry_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/repository"
	"github.com/google/uuid"
)

type JournalEntryService struct {
	repo        repository.JournalEntryRepositoryInterface
	accountRepo repository.GLAccountRepositoryInterface
}

// NewJournalEntryService creates a new journal entry service
func NewJournalEntryService(
	repo repository.JournalEntryRepositoryInterface,
	accountRepo repository.GLAccountRepositoryInterface,
) *JournalEntryService {
	return &JournalEntryService{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

// CreateEntry creates a new journal entry in DRAFT status
func (s *JournalEntryService) CreateEntry(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	// Set defaults
	entry.ID = uuid.New()
	entry.Status = domain.EntryStatusDraft
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()

	// Generate entry number if not provided
	if entry.EntryNumber == "" {
		entryNumber, err := s.GenerateEntryNumber(ctx, entry.OrganizationID, entry.TransactionDate.Format("20060102"))
		if err != nil {
			return nil, fmt.Errorf("failed to generate entry number: %w", err)
		}
		entry.EntryNumber = entryNumber
	}

	// Generate line IDs
	for i := range entry.Lines {
		entry.Lines[i].ID = uuid.New()
		entry.Lines[i].LineNumber = i + 1
	}

	// Domain validation
	if err := entry.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Calculate totals
	entry.CalculateTotals()

	// Save to repository
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to create journal entry: %w", err)
	}

	return entry, nil
}

// UpdateEntry updates an existing draft entry
func (s *JournalEntryService) UpdateEntry(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	// Get existing entry
	existing, err := s.repo.GetByID(ctx, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %w", err)
	}

	// Check if can edit
	if !existing.CanEdit() {
		return nil, fmt.Errorf("entry cannot be edited (status: %s)", existing.Status)
	}

	// Update timestamps
	entry.UpdatedAt = time.Now()
	entry.CreatedAt = existing.CreatedAt
	entry.CreatedBy = existing.CreatedBy

	// Domain validation
	if err := entry.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Calculate totals
	entry.CalculateTotals()

	// Save to repository
	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to update journal entry: %w", err)
	}

	return entry, nil
}

// GetEntry retrieves a journal entry by ID
func (s *JournalEntryService) GetEntry(ctx context.Context, entryID uuid.UUID) (*domain.JournalEntry, error) {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %w", err)
	}

	return entry, nil
}

// ListEntries lists journal entries for an organization
func (s *JournalEntryService) ListEntries(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.JournalEntry, error) {
	if limit <= 0 {
		limit = 50
	}

	entries, err := s.repo.ListByOrganization(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}

	return entries, nil
}

// ListByStatus lists entries by status
func (s *JournalEntryService) ListByStatus(ctx context.Context, orgID uuid.UUID, status domain.EntryStatus, limit, offset int) ([]*domain.JournalEntry, error) {
	if limit <= 0 {
		limit = 50
	}

	if !status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	entries, err := s.repo.ListByStatus(ctx, orgID, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list entries: %w", err)
	}

	return entries, nil
}

// PostEntry posts a draft entry to the ledger
func (s *JournalEntryService) PostEntry(ctx context.Context, entryID uuid.UUID, postedBy uuid.UUID) error {
	// Get entry
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return fmt.Errorf("entry not found: %w", err)
	}

	// Validate can post
	if !entry.CanPost() {
		return fmt.Errorf("entry cannot be posted (status: %s)", entry.Status)
	}

	// Validate entry for posting
	validationResult, err := s.ValidateEntry(ctx, entryID)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !validationResult.IsValid {
		return fmt.Errorf("entry validation failed: %v", validationResult.Errors)
	}

	// Post the entry (domain logic)
	if err := entry.Post(postedBy); err != nil {
		return fmt.Errorf("failed to post entry: %w", err)
	}

	// Save updated entry
	if err := s.repo.Update(ctx, entry); err != nil {
		return fmt.Errorf("failed to save posted entry: %w", err)
	}

	// TODO: Update account balances here
	// This could be done via a separate AccountBalanceService
	// For each line:
	//   - If debit: increase asset/expense, decrease liability/equity/revenue
	//   - If credit: increase liability/equity/revenue, decrease asset/expense

	return nil
}

// VoidEntry voids a posted entry
func (s *JournalEntryService) VoidEntry(ctx context.Context, entryID uuid.UUID) error {
	// Get entry
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return fmt.Errorf("entry not found: %w", err)
	}

	// Validate can void
	if !entry.CanVoid() {
		return fmt.Errorf("entry cannot be voided (status: %s)", entry.Status)
	}

	// Void the entry (domain logic)
	if err := entry.Void(); err != nil {
		return fmt.Errorf("failed to void entry: %w", err)
	}

	// Save updated entry
	if err := s.repo.Update(ctx, entry); err != nil {
		return fmt.Errorf("failed to save voided entry: %w", err)
	}

	return nil
}

// ReverseEntry creates a reversal entry for a posted entry
func (s *JournalEntryService) ReverseEntry(ctx context.Context, entryID uuid.UUID, reversedBy uuid.UUID) (*domain.JournalEntry, error) {
	// Get original entry
	originalEntry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %w", err)
	}

	// Validate can reverse
	if !originalEntry.CanReverse() {
		return nil, fmt.Errorf("entry cannot be reversed (status: %s)", originalEntry.Status)
	}

	// Generate new entry number for reversal
	reversalDate := time.Now().Format("20060102")
	newEntryNumber, err := s.GenerateEntryNumber(ctx, originalEntry.OrganizationID, reversalDate)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reversal entry number: %w", err)
	}

	// Create reversal entry (domain logic)
	reversalEntry, err := originalEntry.CreateReversal(reversedBy, newEntryNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to create reversal: %w", err)
	}

	// Save reversal entry as DRAFT
	if err := s.repo.Create(ctx, reversalEntry); err != nil {
		return nil, fmt.Errorf("failed to save reversal entry: %w", err)
	}

	// Mark original entry as REVERSED
	originalEntry.Status = domain.EntryStatusReversed
	originalEntry.ReversedBy = &reversedBy
	originalEntry.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, originalEntry); err != nil {
		return nil, fmt.Errorf("failed to update original entry: %w", err)
	}

	return reversalEntry, nil
}

// DeleteEntry soft deletes a draft entry
func (s *JournalEntryService) DeleteEntry(ctx context.Context, entryID uuid.UUID) error {
	// Get entry
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return fmt.Errorf("entry not found: %w", err)
	}

	// Can only delete draft entries
	if entry.Status != domain.EntryStatusDraft {
		return fmt.Errorf("only draft entries can be deleted (status: %s)", entry.Status)
	}

	// Delete entry
	if err := s.repo.Delete(ctx, entryID); err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}

	return nil
}

// ValidateEntry validates an entry before posting
func (s *JournalEntryService) ValidateEntry(ctx context.Context, entryID uuid.UUID) (*domain.PostingValidationResult, error) {
	// Get entry
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return nil, fmt.Errorf("entry not found: %w", err)
	}

	// Get all accounts referenced in the entry
	accounts := make(map[uuid.UUID]*domain.GLAccount)
	for _, line := range entry.Lines {
		// ✅ FIXED: Added third parameter (includeInactive bool)
		account, err := s.accountRepo.GetGLAccountByID(ctx, line.AccountID, false)
		if err != nil {
			return nil, fmt.Errorf("account %s not found: %w", line.AccountID, err)
		}
		// ✅ FIXED: Store pointer to account
		accounts[line.AccountID] = &account
	}

	// Perform validation
	validationResult := domain.ValidateForPosting(entry, accounts)

	return validationResult, nil
}

// GenerateEntryNumber generates the next entry number
func (s *JournalEntryService) GenerateEntryNumber(ctx context.Context, orgID uuid.UUID, date string) (string, error) {
	sequence, err := s.repo.GetNextEntryNumber(ctx, orgID, date)
	if err != nil {
		return "", fmt.Errorf("failed to get next sequence: %w", err)
	}

	return domain.GenerateEntryNumber(time.Now(), sequence), nil
}
