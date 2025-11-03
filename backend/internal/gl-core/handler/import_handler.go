package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/service"
)

// ImportHandler handles Excel import requests
type ImportHandler struct {
	importService *service.ImportService
}

// NewImportHandler creates a new ImportHandler
func NewImportHandler(importService *service.ImportService) *ImportHandler {
	return &ImportHandler{
		importService: importService,
	}
}

// ImportChartOfAccountsRequest represents the request for importing chart of accounts
type ImportChartOfAccountsRequest struct {
	LegacySystem   string `form:"legacy_system"`
	LegacyVersion  string `form:"legacy_version"`
	SkipDuplicates bool   `form:"skip_duplicates"`
	UpdateExisting bool   `form:"update_existing"`
	ValidateOnly   bool   `form:"validate_only"`
	MigrationNotes string `form:"migration_notes"`
}

// ImportChartOfAccounts handles chart of accounts import
// @Summary Import chart of accounts from Excel
// @Description Upload Excel file to import chart of accounts
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file"
// @Param legacy_system formData string false "Legacy system (tally, quickbooks, excel)"
// @Param skip_duplicates formData bool false "Skip duplicate accounts"
// @Param update_existing formData bool false "Update existing accounts"
// @Param validate_only formData bool false "Only validate, don't import"
// @Success 200 {object} service.ImportResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gl/import/accounts [post]
func (h *ImportHandler) ImportChartOfAccounts(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 32 MB max
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to parse form",
			Message: err.Error(),
		})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "File is required",
			Message: err.Error(),
		})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid file format",
			Message: "Only .xlsx and .xls files are supported",
		})
		return
	}

	// Save uploaded file temporarily
	tempFile := filepath.Join("/tmp", fmt.Sprintf("import_%s_%s", uuid.New().String(), header.Filename))
	if err := c.SaveUploadedFile(header, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to save file",
			Message: err.Error(),
		})
		return
	}

	// Parse form parameters
	var req ImportChartOfAccountsRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (assuming auth middleware sets this)
	userID := uuid.MustParse(c.GetString("user_id"))
	if userID == uuid.Nil {
		userID = uuid.New() // Fallback for testing
	}

	// Prepare import options
	options := service.ImportOptions{
		LegacySystem:   req.LegacySystem,
		LegacyVersion:  req.LegacyVersion,
		SkipDuplicates: req.SkipDuplicates,
		UpdateExisting: req.UpdateExisting,
		ValidateOnly:   req.ValidateOnly,
		MigrationNotes: req.MigrationNotes,
	}

	// Perform import
	result, err := h.importService.ImportChartOfAccounts(
		c.Request.Context(),
		tempFile,
		header.Filename,
		header.Size,
		userID,
		options,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Import failed",
			Message: err.Error(),
		})
		return
	}

	// Return result
	c.JSON(http.StatusOK, result)
}

// ImportJournalEntriesRequest represents the request for importing journal entries
type ImportJournalEntriesRequest struct {
	LegacySystem   string `form:"legacy_system"`
	LegacyVersion  string `form:"legacy_version"`
	ValidateOnly   bool   `form:"validate_only"`
	MigrationNotes string `form:"migration_notes"`
}

// ImportJournalEntries handles journal entries import
// @Summary Import journal entries from Excel
// @Description Upload Excel file to import journal entries
// @Tags Import
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file"
// @Param legacy_system formData string false "Legacy system"
// @Param validate_only formData bool false "Only validate, don't import"
// @Success 200 {object} service.ImportResult
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/gl/import/journal-entries [post]
func (h *ImportHandler) ImportJournalEntries(c *gin.Context) {
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to parse form",
			Message: err.Error(),
		})
		return
	}

	// Get file from form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "File is required",
			Message: err.Error(),
		})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := filepath.Ext(header.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid file format",
			Message: "Only .xlsx and .xls files are supported",
		})
		return
	}

	// Save uploaded file temporarily
	tempFile := filepath.Join("/tmp", fmt.Sprintf("import_%s_%s", uuid.New().String(), header.Filename))
	if err := c.SaveUploadedFile(header, tempFile); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to save file",
			Message: err.Error(),
		})
		return
	}

	// Parse form parameters
	var req ImportJournalEntriesRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request parameters",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context
	userID := uuid.MustParse(c.GetString("user_id"))
	if userID == uuid.Nil {
		userID = uuid.New() // Fallback for testing
	}

	// Prepare import options
	options := service.ImportOptions{
		LegacySystem:   req.LegacySystem,
		LegacyVersion:  req.LegacyVersion,
		ValidateOnly:   req.ValidateOnly,
		MigrationNotes: req.MigrationNotes,
	}

	// Perform import
	result, err := h.importService.ImportJournalEntries(
		c.Request.Context(),
		tempFile,
		header.Filename,
		header.Size,
		userID,
		options,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Import failed",
			Message: err.Error(),
		})
		return
	}

	// Return result
	c.JSON(http.StatusOK, result)
}

// DownloadTemplate downloads an import template
// @Summary Download import template
// @Description Download Excel template for data import
// @Tags Import
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param type path string true "Template type (chart_of_accounts, journal_entries, tally, quickbooks)"
// @Success 200 {file} binary
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/gl/import/template/{type} [get]
func (h *ImportHandler) DownloadTemplate(c *gin.Context) {
	templateType := c.Param("type")

	if templateType == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Template type is required",
			Message: "Valid types: chart_of_accounts, journal_entries, tally, quickbooks",
		})
		return
	}

	// Get template file path
	filePath, err := h.importService.DownloadTemplateFile(templateType)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Template not found",
			Message: err.Error(),
		})
		return
	}

	// Send file
	c.FileAttachment(filePath, filepath.Base(filePath))
}

// ErrorResponse represents an error response
