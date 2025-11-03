// backend/internal/gl-core/domain/journal_entry.go
package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// JournalEntry represents a complete journal entry with header and lines
type JournalEntry struct {
	ID              uuid.UUID     `json:"id"`
	OrganizationID  uuid.UUID     `json:"organization_id"`
	EntryNumber     string        `json:"entry_number"`     // Auto-generated: JE-20251031-0001
	TransactionDate time.Time     `json:"transaction_date"` // When transaction occurred
	PostingDate     *time.Time    `json:"posting_date"`     // When entry was posted (nil if not posted)
	Reference       string        `json:"reference"`        // External reference (invoice, receipt, etc.)
	Description     string        `json:"description"`      // Entry description
	Status          EntryStatus   `json:"status"`
	Lines           []JournalLine `json:"lines"`        // Entry lines (debits/credits)
	TotalDebit      float64       `json:"total_debit"`  // Calculated total debits
	TotalCredit     float64       `json:"total_credit"` // Calculated total credits
	CreatedBy       uuid.UUID     `json:"created_by"`   // User who created
	PostedBy        *uuid.UUID    `json:"posted_by"`    // User who posted (nil if not posted)
	ReversedBy      *uuid.UUID    `json:"reversed_by"`  // User who reversed (nil if not reversed)
	ReversalOf      *uuid.UUID    `json:"reversal_of"`  // Original entry ID if this is a reversal
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// Validate performs domain validation on JournalEntry
func (je *JournalEntry) Validate() error {
	// Entry number validation
	if je.EntryNumber == "" {
		return NewGLError("entry number is required", ErrJournalEntryNumberRequired)
	}

	// Organization validation
	if je.OrganizationID == uuid.Nil {
		return NewGLError("organization ID is required", ErrJournalOrgRequired)
	}

	// Transaction date validation
	if je.TransactionDate.IsZero() {
		return NewGLError("transaction date is required", ErrJournalTransactionDateRequired)
	}

	// Description validation
	if je.Description == "" {
		return NewGLError("description is required", ErrJournalDescriptionRequired)
	}

	if len(je.Description) > 500 {
		return NewGLError("description cannot exceed 500 characters", ErrJournalDescriptionTooLong)
	}

	// Lines validation
	if len(je.Lines) == 0 {
		return NewGLError("journal entry must have at least one line", ErrJournalNoLines)
	}

	if len(je.Lines) < 2 {
		return NewGLError("journal entry must have at least two lines (debit and credit)", ErrJournalInsufficientLines)
	}

	// Validate each line
	for i, line := range je.Lines {
		if err := line.Validate(); err != nil {
			return NewGLErrorf(ErrJournalLineInvalid, "line %d: %v", i+1, err)
		}
	}

	// Balanced entry validation
	if !je.IsBalanced() {
		return NewGLError("journal entry is not balanced (debits must equal credits)", ErrJournalNotBalanced)
	}

	return nil
}

// CalculateTotals calculates total debits and credits from lines
func (je *JournalEntry) CalculateTotals() {
	je.TotalDebit = 0
	je.TotalCredit = 0

	for _, line := range je.Lines {
		je.TotalDebit += line.Debit
		je.TotalCredit += line.Credit
	}
}

// IsBalanced checks if debits equal credits (within precision tolerance)
func (je *JournalEntry) IsBalanced() bool {
	je.CalculateTotals()

	// Use small tolerance for floating-point comparison (0.01 for 2 decimal places)
	tolerance := 0.01
	difference := je.TotalDebit - je.TotalCredit

	return difference >= -tolerance && difference <= tolerance
}

// CanPost checks if entry can be posted
func (je *JournalEntry) CanPost() bool {
	return je.Status == EntryStatusDraft && je.IsBalanced()
}

// CanVoid checks if entry can be voided
func (je *JournalEntry) CanVoid() bool {
	return je.Status == EntryStatusPosted
}

// CanEdit checks if entry can be edited
func (je *JournalEntry) CanEdit() bool {
	return je.Status == EntryStatusDraft
}

// CanReverse checks if entry can be reversed
func (je *JournalEntry) CanReverse() bool {
	return je.Status == EntryStatusPosted
}

// Post marks entry as posted
func (je *JournalEntry) Post(postedBy uuid.UUID) error {
	if !je.CanPost() {
		return NewGLError("entry cannot be posted", ErrJournalCannotPost)
	}

	now := time.Now()
	je.Status = EntryStatusPosted
	je.PostingDate = &now
	je.PostedBy = &postedBy
	je.UpdatedAt = now

	return nil
}

// Void marks entry as voided
func (je *JournalEntry) Void() error {
	if !je.CanVoid() {
		return NewGLError("entry cannot be voided", ErrJournalCannotVoid)
	}

	je.Status = EntryStatusVoid
	je.UpdatedAt = time.Now()

	return nil
}

// CreateReversal creates a reversal entry
func (je *JournalEntry) CreateReversal(reversedBy uuid.UUID, newEntryNumber string) (*JournalEntry, error) {
	if !je.CanReverse() {
		return nil, NewGLError("entry cannot be reversed", ErrJournalCannotReverse)
	}

	reversal := &JournalEntry{
		ID:              uuid.New(),
		OrganizationID:  je.OrganizationID,
		EntryNumber:     newEntryNumber,
		TransactionDate: time.Now(),
		Reference:       "REV-" + je.Reference,
		Description:     "Reversal of " + je.EntryNumber + ": " + je.Description,
		Status:          EntryStatusDraft,
		Lines:           make([]JournalLine, len(je.Lines)),
		CreatedBy:       reversedBy,
		ReversalOf:      &je.ID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Reverse all lines (swap debits and credits)
	for i, line := range je.Lines {
		reversal.Lines[i] = JournalLine{
			ID:          uuid.New(),
			AccountID:   line.AccountID,
			Description: "Reversal: " + line.Description,
			Debit:       line.Credit, // Swap
			Credit:      line.Debit,  // Swap
			LineNumber:  line.LineNumber,
		}
	}

	reversal.CalculateTotals()

	return reversal, nil
}

// IsReversal checks if this entry is a reversal
func (je *JournalEntry) IsReversal() bool {
	return je.ReversalOf != nil
}

// AddLine adds a new line to the entry
func (je *JournalEntry) AddLine(line JournalLine) error {
	if !je.CanEdit() {
		return NewGLError("cannot add line to posted entry", ErrJournalCannotEdit)
	}

	if err := line.Validate(); err != nil {
		return err
	}

	line.LineNumber = len(je.Lines) + 1
	je.Lines = append(je.Lines, line)
	je.CalculateTotals()

	return nil
}

// RemoveLine removes a line from the entry
func (je *JournalEntry) RemoveLine(lineNumber int) error {
	if !je.CanEdit() {
		return NewGLError("cannot remove line from posted entry", ErrJournalCannotEdit)
	}

	if lineNumber < 1 || lineNumber > len(je.Lines) {
		return NewGLError("invalid line number", ErrJournalLineNotFound)
	}

	// Remove line and renumber
	je.Lines = append(je.Lines[:lineNumber-1], je.Lines[lineNumber:]...)

	for i := range je.Lines {
		je.Lines[i].LineNumber = i + 1
	}

	je.CalculateTotals()

	return nil
}

// GetLineCount returns the number of lines
func (je *JournalEntry) GetLineCount() int {
	return len(je.Lines)
}

// GenerateEntryNumber generates a unique entry number (format: JE-YYYYMMDD-####)
func GenerateEntryNumber(date time.Time, sequence int) string {
	return fmt.Sprintf("JE-%s-%04d", date.Format("20060102"), sequence)
}
