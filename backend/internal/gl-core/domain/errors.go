// backend/internal/gl-core/domain/errors_gl.go
package domain

import "fmt"

// GLError represents a GL domain error
type GLError struct {
    Message string
    Code    string
}

// Error implements the error interface
func (e *GLError) Error() string {
    return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

// NewGLError creates a new GL error
func NewGLError(message, code string) *GLError {
    return &GLError{
        Message: message,
        Code:    code,
    }
}

// NewGLErrorf creates a new GL error with formatted message
func NewGLErrorf(code, format string, args ...interface{}) *GLError {
    return &GLError{
        Message: fmt.Sprintf(format, args...),
        Code:    code,
    }
}

// GL Error Codes
const (
    // Journal Entry errors
    ErrJournalEntryNumberRequired     = "JOURNAL_ENTRY_NUMBER_REQUIRED"
    ErrJournalOrgRequired             = "JOURNAL_ORG_REQUIRED"
    ErrJournalTransactionDateRequired = "JOURNAL_TRANSACTION_DATE_REQUIRED"
    ErrJournalDescriptionRequired     = "JOURNAL_DESCRIPTION_REQUIRED"
    ErrJournalDescriptionTooLong      = "JOURNAL_DESCRIPTION_TOO_LONG"
    ErrJournalNoLines                 = "JOURNAL_NO_LINES"
    ErrJournalInsufficientLines       = "JOURNAL_INSUFFICIENT_LINES"
    ErrJournalLineInvalid             = "JOURNAL_LINE_INVALID"
    ErrJournalNotBalanced             = "JOURNAL_NOT_BALANCED"
    ErrJournalCannotPost              = "JOURNAL_CANNOT_POST"
    ErrJournalCannotVoid              = "JOURNAL_CANNOT_VOID"
    ErrJournalCannotReverse           = "JOURNAL_CANNOT_REVERSE"
    ErrJournalCannotEdit              = "JOURNAL_CANNOT_EDIT"
    ErrJournalLineNotFound            = "JOURNAL_LINE_NOT_FOUND"

    // Journal Line errors
    ErrJournalLineAccountRequired     = "JOURNAL_LINE_ACCOUNT_REQUIRED"
    ErrJournalLineNoAmount            = "JOURNAL_LINE_NO_AMOUNT"
    ErrJournalLineBothAmounts         = "JOURNAL_LINE_BOTH_AMOUNTS"
    ErrJournalLineNegativeDebit       = "JOURNAL_LINE_NEGATIVE_DEBIT"
    ErrJournalLineNegativeCredit      = "JOURNAL_LINE_NEGATIVE_CREDIT"
    ErrJournalLineDescriptionRequired = "JOURNAL_LINE_DESCRIPTION_REQUIRED"
    ErrJournalLineDescriptionTooLong  = "JOURNAL_LINE_DESCRIPTION_TOO_LONG"
    ErrJournalLineReferenceTooLong    = "JOURNAL_LINE_REFERENCE_TOO_LONG"
)
