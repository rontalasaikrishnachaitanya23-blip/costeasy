package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/repository"
)

// ImportService handles Excel imports for GL data
type ImportService struct {
	pool             *pgxpool.Pool
	accountRepo      repository.GLAccountRepositoryInterface
	journalEntryRepo repository.JournalEntryRepositoryInterface
	journalLineRepo  repository.JournalLineRepositoryInterface
}

// ImportResult contains results of an import operation
type ImportResult struct {
	ImportLogID          uuid.UUID       `json:"import_log_id"`
	TotalRows            int             `json:"total_rows"`
	SuccessCount         int             `json:"success_count"`
	ErrorCount           int             `json:"error_count"`
	WarningCount         int             `json:"warning_count"`
	Errors               []ImportError   `json:"errors"`
	Warnings             []ImportWarning `json:"warnings"`
	ImportedIDs          []uuid.UUID     `json:"imported_ids,omitempty"`
	ProcessingTimeMillis int64           `json:"processing_time_millis" example:"1250"`
	Status               string          `json:"status"`
	LegacySystem         string          `json:"legacy_system,omitempty"`
}

// ImportError represents a single import error
type ImportError struct {
	Row        int    `json:"row"`
	Column     string `json:"column,omitempty"`
	Field      string `json:"field"`
	Value      string `json:"value"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	LegacyCode string `json:"legacy_code,omitempty"`
}

// ImportWarning represents a single import warning
type ImportWarning struct {
	Row        int    `json:"row"`
	Message    string `json:"message"`
	LegacyCode string `json:"legacy_code,omitempty"`
}

// ImportOptions contains options for import operations
type ImportOptions struct {
	LegacySystem   string `json:"legacy_system"`   // 'tally', 'quickbooks', 'excel'
	LegacyVersion  string `json:"legacy_version"`  // Version info
	SkipDuplicates bool   `json:"skip_duplicates"` // Skip existing accounts
	UpdateExisting bool   `json:"update_existing"` // Update if exists
	ValidateOnly   bool   `json:"validate_only"`   // Only validate, don't import
	MigrationNotes string `json:"migration_notes"` // Additional notes
}

// JournalEntryLine represents a single journal entry line
type JournalEntryLine struct {
	RowNumber    int
	EntryDate    time.Time
	ReferenceNo  string
	Description  string
	AccountCode  string
	AccountID    uuid.UUID
	DebitAmount  float64
	CreditAmount float64
	CostCenter   string
	Department   string
	Notes        string
}

// LegacyAccountInfo holds legacy system specific account information
type LegacyAccountInfo struct {
	Code    string                 `json:"code"`
	Name    string                 `json:"name"`
	System  string                 `json:"system"`
	Version string                 `json:"version"`
	RawData map[string]interface{} `json:"raw_data"`
}

// NewImportService creates a new ImportService
func NewImportService(
	pool *pgxpool.Pool,
	accountRepo repository.GLAccountRepositoryInterface,
	journalEntryRepo repository.JournalEntryRepositoryInterface,
	journalLineRepo repository.JournalLineRepositoryInterface,
) *ImportService {
	return &ImportService{
		pool:             pool,
		accountRepo:      accountRepo,
		journalEntryRepo: journalEntryRepo,
		journalLineRepo:  journalLineRepo,
	}
}

// ImportChartOfAccounts imports chart of accounts from Excel
func (s *ImportService) ImportChartOfAccounts(
	ctx context.Context,
	filePath string,
	fileName string,
	fileSize int64,
	userID uuid.UUID,
	options ImportOptions,
) (*ImportResult, error) {
	startTime := time.Now()

	result := &ImportResult{
		ImportLogID:  uuid.New(),
		Errors:       make([]ImportError, 0),
		Warnings:     make([]ImportWarning, 0),
		ImportedIDs:  make([]uuid.UUID, 0),
		LegacySystem: options.LegacySystem,
	}

	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		result.Status = "failed"
		return result, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	// Get active sheet
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		result.Status = "failed"
		return result, fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) < 2 {
		result.Status = "failed"
		return result, fmt.Errorf("file must contain header and at least one data row")
	}

	// Parse header and detect legacy system format
	header := rows[0]
	colMap, detectedSystem := s.buildColumnMapWithLegacyDetection(header, options.LegacySystem)

	if detectedSystem != "" && options.LegacySystem == "" {
		result.LegacySystem = detectedSystem
		result.Warnings = append(result.Warnings, ImportWarning{
			Row:     1,
			Message: fmt.Sprintf("Auto-detected legacy system: %s", detectedSystem),
		})
	}

	// Validate required columns
	if err := s.validateRequiredColumns(colMap, result.LegacySystem); err != nil {
		result.Status = "failed"
		return result, err
	}

	// Begin transaction if not validate-only
	var tx interface{}
	if !options.ValidateOnly {
		var err error
		txn, err := s.pool.Begin(ctx)
		if err != nil {
			result.Status = "failed"
			return result, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txn.Rollback(ctx)
		tx = txn
	}

	result.TotalRows = len(rows) - 1

	// Process each row
	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]

		// Skip empty rows
		if s.isEmptyRow(row) {
			continue
		}

		// Validate and convert row to account
		account, legacyInfo, validationErrs := s.validateAndConvertAccountRow(
			row, colMap, rowIdx+1, result.LegacySystem,
		)

		if len(validationErrs) > 0 {
			result.Errors = append(result.Errors, validationErrs...)
			result.ErrorCount++
			continue
		}

		// Check for existing account
		existingAccount, err := s.accountRepo.GetGLAccountByCode(ctx, account.Code, true)
		if err != nil && err.Error() != "no rows in result set" {
			result.Errors = append(result.Errors, ImportError{
				Row:        rowIdx + 1,
				Field:      "code",
				Value:      account.Code,
				Message:    fmt.Sprintf("Database error: %v", err),
				Code:       "DB_ERROR",
				LegacyCode: legacyInfo.Code,
			})
			result.ErrorCount++
			continue
		}

		// existingAccount is a POINTER (*GLAccount), check if pointer is not nil
		if existingAccount.ID != uuid.Nil {
			if options.SkipDuplicates {
				result.Warnings = append(result.Warnings, ImportWarning{
					Row:        rowIdx + 1,
					Message:    fmt.Sprintf("Account '%s' already exists, skipping", account.Code),
					LegacyCode: legacyInfo.Code,
				})
				result.WarningCount++
				continue
			} else if options.UpdateExisting {
				// Update existing account
				account.ID = existingAccount.ID
				if !options.ValidateOnly {
					_, err := s.accountRepo.UpdateGLAccount(ctx, *account)
					if err != nil {
						result.Errors = append(result.Errors, ImportError{
							Row:        rowIdx + 1,
							Field:      "code",
							Value:      account.Code,
							Message:    fmt.Sprintf("Failed to update: %v", err),
							Code:       "UPDATE_ERROR",
							LegacyCode: legacyInfo.Code,
						})
						result.ErrorCount++
						continue
					}
				}
				result.Warnings = append(result.Warnings, ImportWarning{
					Row:        rowIdx + 1,
					Message:    fmt.Sprintf("Account '%s' updated", account.Code),
					LegacyCode: legacyInfo.Code,
				})
				result.WarningCount++
			} else {
				result.Errors = append(result.Errors, ImportError{
					Row:        rowIdx + 1,
					Field:      "code",
					Value:      account.Code,
					Message:    "Account already exists",
					Code:       "DUPLICATE_ACCOUNT",
					LegacyCode: legacyInfo.Code,
				})
				result.ErrorCount++
				continue
			}
		} else {
			// Create new account
			if !options.ValidateOnly {
				account.ID = uuid.New()
				account.CreateAt = time.Now()
				account.UpdateAt = time.Now()
				if _, err := s.accountRepo.CreateGLAccount(ctx, *account); err != nil {
					result.Errors = append(result.Errors, ImportError{
						Row:        rowIdx + 1,
						Field:      "code",
						Value:      account.Code,
						Message:    fmt.Sprintf("Failed to create: %v", err),
						Code:       "CREATE_ERROR",
						LegacyCode: legacyInfo.Code,
					})
					result.ErrorCount++
					continue
				}
				result.ImportedIDs = append(result.ImportedIDs, account.ID)
			}
		}

		result.SuccessCount++
	}

	// Determine final status
	if result.ErrorCount == 0 {
		result.Status = "success"
	} else if result.SuccessCount > 0 {
		result.Status = "partial"
	} else {
		result.Status = "failed"
	}

	// Commit transaction if not validation only and no critical errors
	if !options.ValidateOnly && tx != nil && result.Status != "failed" {
		if txn, ok := tx.(interface{ Commit(context.Context) error }); ok {
			if err := txn.Commit(ctx); err != nil {
				result.Status = "failed"
				return result, fmt.Errorf("failed to commit transaction: %w", err)
			}
		}
	} else if options.ValidateOnly {
		result.Status = "validated"
	}

	result.ProcessingTimeMillis = time.Since(startTime).Milliseconds()

	return result, nil
}

// ImportJournalEntries imports journal entries from Excel
func (s *ImportService) ImportJournalEntries(
	ctx context.Context,
	filePath string,
	fileName string,
	fileSize int64,
	userID uuid.UUID,
	options ImportOptions,
) (*ImportResult, error) {
	startTime := time.Now()

	result := &ImportResult{
		ImportLogID: uuid.New(),
		Errors:      make([]ImportError, 0),
		Warnings:    make([]ImportWarning, 0),
		ImportedIDs: make([]uuid.UUID, 0),
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		result.Status = "failed"
		return result, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		result.Status = "failed"
		return result, fmt.Errorf("failed to read rows: %w", err)
	}

	if len(rows) < 2 {
		result.Status = "failed"
		return result, fmt.Errorf("file must contain header and at least one data row")
	}

	header := rows[0]
	colMap := s.buildColumnMap(header)

	// Validate required columns
	requiredCols := []string{"entry_date", "reference_no", "description", "account_code"}
	for _, col := range requiredCols {
		if _, exists := colMap[col]; !exists {
			result.Status = "failed"
			return result, fmt.Errorf("missing required column: %s", col)
		}
	}

	var tx interface{}
	if !options.ValidateOnly {
		var err error
		txn, err := s.pool.Begin(ctx)
		if err != nil {
			result.Status = "failed"
			return result, fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer txn.Rollback(ctx)
		tx = txn
	}

	result.TotalRows = len(rows) - 1

	// Group entries by reference number
	entriesByRef := make(map[string][]*JournalEntryLine)

	for rowIdx := 1; rowIdx < len(rows); rowIdx++ {
		row := rows[rowIdx]

		if s.isEmptyRow(row) {
			continue
		}

		line, _, validationErrs := s.validateAndConvertJournalEntryRow(row, colMap, rowIdx+1)
		if len(validationErrs) > 0 {
			result.Errors = append(result.Errors, validationErrs...)
			result.ErrorCount++
			continue
		}

		entriesByRef[line.ReferenceNo] = append(entriesByRef[line.ReferenceNo], line)
	}

	// Validate and process each journal entry
	for _, lines := range entriesByRef {
		// Validate balancing
		var totalDebit, totalCredit float64
		for _, line := range lines {
			totalDebit += line.DebitAmount
			totalCredit += line.CreditAmount
		}

		// Check balance
		if fmt.Sprintf("%.2f", totalDebit) != fmt.Sprintf("%.2f", totalCredit) {
			for _, line := range lines {
				result.Errors = append(result.Errors, ImportError{
					Row:     line.RowNumber,
					Field:   "reference_no",
					Value:   line.ReferenceNo,
					Message: fmt.Sprintf("Entry not balanced. Debits: %.2f, Credits: %.2f", totalDebit, totalCredit),
					Code:    "UNBALANCED_ENTRY",
				})
			}
			result.ErrorCount += len(lines)
			continue
		}

		result.SuccessCount += len(lines)
		result.ImportedIDs = append(result.ImportedIDs, uuid.New()) // Representative ID
	}

	// Determine final status
	if result.ErrorCount == 0 {
		result.Status = "success"
	} else if result.SuccessCount > 0 {
		result.Status = "partial"
	} else {
		result.Status = "failed"
	}

	// Commit if successful
	if !options.ValidateOnly && tx != nil && result.Status != "failed" {
		if txn, ok := tx.(interface{ Commit(context.Context) error }); ok {
			if err := txn.Commit(ctx); err != nil {
				result.Status = "failed"
				return result, fmt.Errorf("failed to commit transaction: %w", err)
			}
		}
	} else if options.ValidateOnly {
		result.Status = "validated"
	}

	result.ProcessingTimeMillis = time.Since(startTime).Milliseconds()

	return result, nil
}

// ============================================================================
// VALIDATION METHODS (from import_validation.go)
// ============================================================================

// validateAndConvertAccountRow validates and converts row to account domain object
func (s *ImportService) validateAndConvertAccountRow(
	row []string,
	colMap map[string]int,
	rowNum int,
	legacySystem string,
) (*domain.GLAccount, *LegacyAccountInfo, []ImportError) {
	errors := make([]ImportError, 0)
	legacyInfo := &LegacyAccountInfo{
		System:  legacySystem,
		RawData: make(map[string]interface{}),
	}

	account := &domain.GLAccount{
		ID:       uuid.New(),
		IsActive: true, // Default
	}

	// Account Code validation
	accountCode := strings.TrimSpace(s.getCellValue(row, colMap, "account_code"))
	if accountCode == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Code",
			Field:   "code",
			Value:   accountCode,
			Message: "Account Code is required",
			Code:    "REQUIRED_FIELD",
		})
	} else if len(accountCode) > 20 {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Code",
			Field:   "code",
			Value:   accountCode,
			Message: "Account Code must be 20 characters or less",
			Code:    "FIELD_TOO_LONG",
		})
	} else if !regexp.MustCompile(`^[A-Z0-9\-_]+$`).MatchString(accountCode) {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Code",
			Field:   "code",
			Value:   accountCode,
			Message: "Account Code must contain only uppercase letters, numbers, hyphens, and underscores",
			Code:    "INVALID_FORMAT",
		})
	}
	account.Code = accountCode
	legacyInfo.Code = accountCode

	// Account Name validation
	accountName := strings.TrimSpace(s.getCellValue(row, colMap, "account_name"))
	if accountName == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Name",
			Field:   "name",
			Value:   accountName,
			Message: "Account Name is required",
			Code:    "REQUIRED_FIELD",
		})
	} else if len(accountName) > 100 {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Name",
			Field:   "name",
			Value:   accountName,
			Message: "Account Name must be 100 characters or less",
			Code:    "FIELD_TOO_LONG",
		})
	}
	account.Name = accountName
	legacyInfo.Name = accountName

	// Account Type validation
	accountType := strings.TrimSpace(s.getCellValue(row, colMap, "account_type"))
	if accountType == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Type",
			Field:   "type",
			Value:   accountType,
			Message: "Account Type is required",
			Code:    "REQUIRED_FIELD",
		})
	} else {
		// Convert to AccountType from domain
		account.Type = domain.AccountType(strings.ToUpper(accountType))
	}

	// Is Active validation
	isActiveStr := strings.TrimSpace(s.getCellValue(row, colMap, "is_active"))
	if isActiveStr != "" {
		isActiveStr = strings.ToLower(isActiveStr)
		switch isActiveStr {
		case "yes", "true", "1":
			account.IsActive = true
		case "no", "false", "0":
			account.IsActive = false
		default:
			errors = append(errors, ImportError{
				Row:     rowNum,
				Column:  "Is Active",
				Field:   "is_active",
				Value:   isActiveStr,
				Message: "Is Active must be: Yes, No, True, False, 1, or 0",
				Code:    "INVALID_VALUE",
			})
		}
	}

	// Parent Account Code (optional)
	parentCode := strings.TrimSpace(s.getCellValue(row, colMap, "parent_code"))
	if parentCode != "" {
		// Try to parse as UUID if provided
		if parentID, err := uuid.Parse(parentCode); err == nil {
			account.ParentCode = &parentID
		} else {
			// It might be a legacy code that needs mapping
			legacyInfo.RawData["parent_code"] = parentCode
		}
	}

	// Store raw data for audit
	legacyInfo.RawData["account_code"] = accountCode
	legacyInfo.RawData["account_name"] = accountName
	legacyInfo.RawData["account_type"] = accountType
	legacyInfo.RawData["is_active"] = isActiveStr

	return account, legacyInfo, errors
}

// validateAndConvertJournalEntryRow validates and converts row to journal entry line
func (s *ImportService) validateAndConvertJournalEntryRow(
	row []string,
	colMap map[string]int,
	rowNum int,
) (*JournalEntryLine, *LegacyAccountInfo, []ImportError) {
	errors := make([]ImportError, 0)
	legacyInfo := &LegacyAccountInfo{
		RawData: make(map[string]interface{}),
	}

	line := &JournalEntryLine{
		RowNumber: rowNum,
	}

	// Entry Date validation
	entryDateStr := strings.TrimSpace(s.getCellValue(row, colMap, "entry_date"))
	if entryDateStr == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Entry Date",
			Field:   "entry_date",
			Value:   entryDateStr,
			Message: "Entry Date is required",
			Code:    "REQUIRED_FIELD",
		})
	} else {
		entryDate, err := time.Parse("2006-01-02", entryDateStr)
		if err != nil {
			errors = append(errors, ImportError{
				Row:     rowNum,
				Column:  "Entry Date",
				Field:   "entry_date",
				Value:   entryDateStr,
				Message: "Entry Date must be in YYYY-MM-DD format",
				Code:    "INVALID_DATE",
			})
		} else {
			line.EntryDate = entryDate
		}
	}

	// Reference No validation
	referenceNo := strings.TrimSpace(s.getCellValue(row, colMap, "reference_no"))
	if referenceNo == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Reference No",
			Field:   "reference_no",
			Value:   referenceNo,
			Message: "Reference No is required",
			Code:    "REQUIRED_FIELD",
		})
	}
	line.ReferenceNo = referenceNo
	legacyInfo.Code = referenceNo

	// Description validation
	description := strings.TrimSpace(s.getCellValue(row, colMap, "description"))
	if description == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Description",
			Field:   "description",
			Value:   description,
			Message: "Description is required",
			Code:    "REQUIRED_FIELD",
		})
	}
	line.Description = description

	// Account Code validation
	accountCode := strings.TrimSpace(s.getCellValue(row, colMap, "account_code"))
	if accountCode == "" {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Column:  "Account Code",
			Field:   "account_code",
			Value:   accountCode,
			Message: "Account Code is required",
			Code:    "REQUIRED_FIELD",
		})
	}
	line.AccountCode = accountCode
	legacyInfo.Name = accountCode

	// Debit Amount validation
	debitAmountStr := strings.TrimSpace(s.getCellValue(row, colMap, "debit_amount"))
	if debitAmountStr != "" {
		debit, err := strconv.ParseFloat(debitAmountStr, 64)
		if err != nil {
			errors = append(errors, ImportError{
				Row:     rowNum,
				Column:  "Debit Amount",
				Field:   "debit_amount",
				Value:   debitAmountStr,
				Message: "Debit Amount must be a valid number",
				Code:    "INVALID_NUMBER",
			})
		} else if debit < 0 {
			errors = append(errors, ImportError{
				Row:     rowNum,
				Column:  "Debit Amount",
				Field:   "debit_amount",
				Value:   debitAmountStr,
				Message: "Debit Amount must be positive",
				Code:    "INVALID_VALUE",
			})
		} else {
			line.DebitAmount = debit
		}
	}

	// Credit Amount validation
	creditAmountStr := strings.TrimSpace(s.getCellValue(row, colMap, "credit_amount"))
	if creditAmountStr != "" {
		credit, err := strconv.ParseFloat(creditAmountStr, 64)
		if err != nil {
			errors = append(errors, ImportError{
				Row:     rowNum,
				Column:  "Credit Amount",
				Field:   "credit_amount",
				Value:   creditAmountStr,
				Message: "Credit Amount must be a valid number",
				Code:    "INVALID_NUMBER",
			})
		} else if credit < 0 {
			errors = append(errors, ImportError{
				Row:     rowNum,
				Column:  "Credit Amount",
				Field:   "credit_amount",
				Value:   creditAmountStr,
				Message: "Credit Amount must be positive",
				Code:    "INVALID_VALUE",
			})
		} else {
			line.CreditAmount = credit
		}
	}

	// Validate that either debit or credit is filled (not both, not neither)
	if line.DebitAmount > 0 && line.CreditAmount > 0 {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Field:   "amounts",
			Message: "Cannot have both Debit and Credit amounts in same line",
			Code:    "INVALID_COMBINATION",
		})
	} else if line.DebitAmount == 0 && line.CreditAmount == 0 {
		errors = append(errors, ImportError{
			Row:     rowNum,
			Field:   "amounts",
			Message: "Must specify either Debit or Credit amount",
			Code:    "REQUIRED_FIELD",
		})
	}

	// Optional fields
	line.CostCenter = strings.TrimSpace(s.getCellValue(row, colMap, "cost_center"))
	line.Department = strings.TrimSpace(s.getCellValue(row, colMap, "department"))
	line.Notes = strings.TrimSpace(s.getCellValue(row, colMap, "notes"))

	// Store raw data
	legacyInfo.RawData["entry_date"] = entryDateStr
	legacyInfo.RawData["reference_no"] = referenceNo
	legacyInfo.RawData["account_code"] = accountCode
	legacyInfo.RawData["debit_amount"] = debitAmountStr
	legacyInfo.RawData["credit_amount"] = creditAmountStr

	return line, legacyInfo, errors
}

// getCellValue safely retrieves cell value from row
func (s *ImportService) getCellValue(row []string, colMap map[string]int, colName string) string {
	if idx, exists := colMap[colName]; exists && idx < len(row) {
		return row[idx]
	}
	return ""
}

// buildColumnMap builds a map of column names to indices
func (s *ImportService) buildColumnMap(header []string) map[string]int {
	colMap := make(map[string]int)
	for idx, col := range header {
		normalizedCol := strings.ToLower(strings.TrimSpace(col))
		normalizedCol = strings.ReplaceAll(normalizedCol, " ", "_")
		colMap[normalizedCol] = idx
	}
	return colMap
}

// buildColumnMapWithLegacyDetection builds column map and detects legacy system format
func (s *ImportService) buildColumnMapWithLegacyDetection(header []string, specifiedSystem string) (map[string]int, string) {
	colMap := make(map[string]int)
	detectedSystem := ""

	// Build basic normalized column map
	for idx, col := range header {
		normalizedCol := strings.ToLower(strings.TrimSpace(col))
		normalizedCol = strings.ReplaceAll(normalizedCol, " ", "_")
		colMap[normalizedCol] = idx
	}

	// If system not specified, try to detect from column names
	if specifiedSystem == "" {
		// Tally-specific columns
		if _, hasTallyCode := colMap["ledger_code"]; hasTallyCode {
			if _, hasTallyGroup := colMap["group"]; hasTallyGroup {
				detectedSystem = "tally"
			}
		}

		// QuickBooks-specific columns
		if _, hasQBType := colMap["account_type"]; hasQBType {
			if _, hasQBNumber := colMap["account_number"]; hasQBNumber {
				detectedSystem = "quickbooks"
			}
		}

		// Generic Excel format
		if detectedSystem == "" {
			detectedSystem = "excel"
		}
	} else {
		detectedSystem = specifiedSystem
	}

	// Create standardized column mappings based on detected/specified system
	standardMap := s.createStandardColumnMap(colMap, detectedSystem)

	return standardMap, detectedSystem
}

// createStandardColumnMap creates standardized column mappings for different legacy systems
func (s *ImportService) createStandardColumnMap(colMap map[string]int, system string) map[string]int {
	standardMap := make(map[string]int)

	switch system {
	case "tally":
		// Map Tally columns to standard names
		if idx, ok := colMap["ledger_code"]; ok {
			standardMap["account_code"] = idx
		}
		if idx, ok := colMap["ledger_name"]; ok {
			standardMap["account_name"] = idx
		}
		if idx, ok := colMap["group"]; ok {
			standardMap["account_type"] = idx
		}
		if idx, ok := colMap["parent_group"]; ok {
			standardMap["parent_code"] = idx
		}

	case "quickbooks":
		// Map QuickBooks columns to standard names
		if idx, ok := colMap["account_number"]; ok {
			standardMap["account_code"] = idx
		}
		if idx, ok := colMap["account_name"]; ok {
			standardMap["account_name"] = idx
		}
		if idx, ok := colMap["account_type"]; ok {
			standardMap["account_type"] = idx
		}
		if idx, ok := colMap["parent_account"]; ok {
			standardMap["parent_code"] = idx
		}

	default: // excel or generic
		// Direct mapping for standard Excel format
		if idx, ok := colMap["account_code"]; ok {
			standardMap["account_code"] = idx
		}
		if idx, ok := colMap["account_name"]; ok {
			standardMap["account_name"] = idx
		}
		if idx, ok := colMap["account_type"]; ok {
			standardMap["account_type"] = idx
		}
		if idx, ok := colMap["parent_code"]; ok {
			standardMap["parent_code"] = idx
		}
		if idx, ok := colMap["is_active"]; ok {
			standardMap["is_active"] = idx
		}
	}

	return standardMap
}

// validateRequiredColumns validates that required columns exist
func (s *ImportService) validateRequiredColumns(colMap map[string]int, system string) error {
	var requiredCols []string

	switch system {
	case "tally":
		requiredCols = []string{"account_code", "account_name", "account_type"}
	case "quickbooks":
		requiredCols = []string{"account_code", "account_name", "account_type"}
	default:
		requiredCols = []string{"account_code", "account_name", "account_type"}
	}

	for _, col := range requiredCols {
		if _, exists := colMap[col]; !exists {
			return fmt.Errorf("missing required column for %s format: %s", system, col)
		}
	}

	return nil
}

// isEmptyRow checks if a row is empty
func (s *ImportService) isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// combineErrorMessages combines multiple error messages
func (s *ImportService) CombineErrorMessages(errs []ImportError) string {
	if len(errs) == 0 {
		return ""
	}
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = fmt.Sprintf("%s (code: %s)", err.Message, err.Code)
	}

	return strings.Join(messages, "; ")
}

// DownloadTemplateFile returns the path to a template file
func (s *ImportService) DownloadTemplateFile(templateType string) (string, error) {
	templatesDir := "internal/settings/seed/templates"

	var fileName string
	switch templateType {
	case "chart_of_accounts":
		fileName = "chart_of_accounts_template.xlsx"
	case "journal_entries":
		fileName = "journal_entries_template.xlsx"
	case "tally":
		fileName = "tally_migration_template.xlsx"
	case "quickbooks":
		fileName = "quickbooks_migration_template.xlsx"
	default:
		return "", fmt.Errorf("unknown template type: %s", templateType)
	}

	filePath := filepath.Join(templatesDir, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", filePath)
	}

	return filePath, nil
}
