// backend/internal/gl-core/handler/dto/journal_entry_dto.go
package dto

// CreateJournalEntryRequest represents the request body for creating a journal entry
type CreateJournalEntryRequest struct {
    OrganizationID  string               `json:"organization_id" binding:"required"`
    TransactionDate string               `json:"transaction_date" binding:"required"` // YYYY-MM-DD
    Reference       string               `json:"reference"`
    Description     string               `json:"description" binding:"required"`
    Lines           []JournalLineRequest `json:"lines" binding:"required,min=2"`
}

// JournalLineRequest represents a journal line in the request
type JournalLineRequest struct {
    AccountID   string  `json:"account_id" binding:"required"`
    Reference   string  `json:"reference"`
    Description string  `json:"description" binding:"required"`
    Debit       float64 `json:"debit"`
    Credit      float64 `json:"credit"`
}

// UpdateJournalEntryRequest represents the request body for updating a journal entry
type UpdateJournalEntryRequest struct {
    TransactionDate string               `json:"transaction_date" binding:"required"` // YYYY-MM-DD
    Reference       string               `json:"reference"`
    Description     string               `json:"description" binding:"required"`
    Lines           []JournalLineRequest `json:"lines" binding:"required,min=2"`
}

// JournalEntryResponse represents the response for a journal entry
type JournalEntryResponse struct {
    ID              string                `json:"id"`
    OrganizationID  string                `json:"organization_id"`
    EntryNumber     string                `json:"entry_number"`
    TransactionDate string                `json:"transaction_date"`
    PostingDate     *string               `json:"posting_date"`
    Reference       string                `json:"reference"`
    Description     string                `json:"description"`
    Status          string                `json:"status"`
    TotalDebit      float64               `json:"total_debit"`
    TotalCredit     float64               `json:"total_credit"`
    Lines           []JournalLineResponse `json:"lines"`
    CreatedAt       string                `json:"created_at"`
    UpdatedAt       string                `json:"updated_at"`
}

// JournalLineResponse represents a journal line in the response
type JournalLineResponse struct {
    ID          string  `json:"id"`
    AccountID   string  `json:"account_id"`
    LineNumber  int     `json:"line_number"`
    Reference   string  `json:"reference"`
    Description string  `json:"description"`
    Debit       float64 `json:"debit"`
    Credit      float64 `json:"credit"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
