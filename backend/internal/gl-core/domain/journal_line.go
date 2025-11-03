// backend/internal/gl-core/domain/journal_line.go
package domain

import (
	"github.com/google/uuid"
)

// JournalLine represents a single line in a journal entry
type JournalLine struct {
	ID          uuid.UUID `json:"id"`
	AccountID   uuid.UUID `json:"account_id"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Debit       float64   `json:"debit"`
	Credit      float64   `json:"credit"`
	LineNumber  int       `json:"line_number"`
}

// Validate performs domain validation on JournalLine
func (jl *JournalLine) Validate() error {
	if jl.AccountID == uuid.Nil {
		return NewGLError("account ID is required", ErrJournalLineAccountRequired)
	}

	if jl.Debit == 0 && jl.Credit == 0 {
		return NewGLError("line must have either debit or credit amount", ErrJournalLineNoAmount)
	}

	if jl.Debit != 0 && jl.Credit != 0 {
		return NewGLError("line cannot have both debit and credit amounts", ErrJournalLineBothAmounts)
	}

	if jl.Debit < 0 {
		return NewGLError("debit amount cannot be negative", ErrJournalLineNegativeDebit)
	}

	if jl.Credit < 0 {
		return NewGLError("credit amount cannot be negative", ErrJournalLineNegativeCredit)
	}

	if jl.Description == "" {
		return NewGLError("line description is required", ErrJournalLineDescriptionRequired)
	}

	if len(jl.Description) > 255 {
		return NewGLError("line description cannot exceed 255 characters", ErrJournalLineDescriptionTooLong)
	}

	if len(jl.Reference) > 100 {
		return NewGLError("line reference cannot exceed 100 characters", ErrJournalLineReferenceTooLong)
	}

	return nil
}

// IsDebit checks if this is a debit line
func (jl *JournalLine) IsDebit() bool {
	return jl.Debit > 0
}

// IsCredit checks if this is a credit line
func (jl *JournalLine) IsCredit() bool {
	return jl.Credit > 0
}

// GetAmount returns the non-zero amount
func (jl *JournalLine) GetAmount() float64 {
	if jl.IsDebit() {
		return jl.Debit
	}
	return jl.Credit
}

// GetAmountType returns "debit" or "credit"
func (jl *JournalLine) GetAmountType() string {
	if jl.IsDebit() {
		return "debit"
	}
	return "credit"
}

// Reverse creates a reversed copy
func (jl *JournalLine) Reverse() JournalLine {
	return JournalLine{
		ID:          uuid.New(),
		AccountID:   jl.AccountID,
		Reference:   jl.Reference,
		Description: "Reversal: " + jl.Description,
		Debit:       jl.Credit,
		Credit:      jl.Debit,
		LineNumber:  jl.LineNumber,
	}
}
