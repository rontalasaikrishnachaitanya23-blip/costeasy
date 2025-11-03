// backend/internal/gl-core/domain/posting_validation.go
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PostingValidationResult represents the result of posting validation
type PostingValidationResult struct {
	IsValid  bool
	Errors   []string
	Warnings []string
}

// AddError adds an error to the validation result
func (pvr *PostingValidationResult) AddError(err string) {
	pvr.IsValid = false
	pvr.Errors = append(pvr.Errors, err)
}

// AddWarning adds a warning to the validation result
func (pvr *PostingValidationResult) AddWarning(warning string) {
	pvr.Warnings = append(pvr.Warnings, warning)
}

// HasErrors checks if there are any errors
func (pvr *PostingValidationResult) HasErrors() bool {
	return len(pvr.Errors) > 0
}

// HasWarnings checks if there are any warnings
func (pvr *PostingValidationResult) HasWarnings() bool {
	return len(pvr.Warnings) > 0
}

// ValidateForPosting performs comprehensive validation before posting
func ValidateForPosting(entry *JournalEntry, accounts map[uuid.UUID]*GLAccount) *PostingValidationResult {
	result := &PostingValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Basic domain validation
	if err := entry.Validate(); err != nil {
		result.AddError(err.Error())
		return result
	}

	// Status validation
	if entry.Status != EntryStatusDraft {
		result.AddError(fmt.Sprintf("entry must be in DRAFT status to be posted (current: %s)", entry.Status))
	}

	// Balance validation
	if !entry.IsBalanced() {
		result.AddError(fmt.Sprintf("entry is not balanced: debits %.2f != credits %.2f", entry.TotalDebit, entry.TotalCredit))
	}

	// Account existence and status validation
	for i, line := range entry.Lines {
		account, exists := accounts[line.AccountID]
		if !exists {
			result.AddError(fmt.Sprintf("line %d: account %s does not exist", i+1, line.AccountID))
			continue
		}

		if !account.IsActive {
			result.AddError(fmt.Sprintf("line %d: account %s (%s) is inactive", i+1, account.Code, account.Name))
		}

		// Warn if posting to control account
		if account.IsParentAccount() {
			result.AddWarning(fmt.Sprintf("line %d: posting to control account %s (%s)", i+1, account.Code, account.Name))
		}
	}

	// Transaction date validation
	if entry.TransactionDate.After(time.Now().Add(24 * time.Hour)) {
		result.AddError("transaction date cannot be in the future")
	}

	// Check for backdated entries (more than 30 days old)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	if entry.TransactionDate.Before(thirtyDaysAgo) {
		result.AddWarning(fmt.Sprintf("transaction date is more than 30 days old (%s)", entry.TransactionDate.Format("2006-01-02")))
	}

	return result
}

// ValidateBalance validates that debits equal credits
func ValidateBalance(totalDebit, totalCredit float64) error {
	tolerance := 0.01
	difference := totalDebit - totalCredit

	if difference < -tolerance || difference > tolerance {
		return NewGLErrorf(ErrJournalNotBalanced, "entry is not balanced: debits %.2f != credits %.2f (difference: %.2f)", totalDebit, totalCredit, difference)
	}

	return nil
}
