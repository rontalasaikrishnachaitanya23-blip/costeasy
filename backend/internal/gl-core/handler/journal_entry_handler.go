// backend/internal/gl-core/handler/journal_entry_handler.go
package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler/dto"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler/mapper"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type JournalEntryHandler struct {
	service service.JournalEntryServiceInterface
}

// NewJournalEntryHandler creates a new journal entry handler
func NewJournalEntryHandler(service service.JournalEntryServiceInterface) *JournalEntryHandler {
	return &JournalEntryHandler{service: service}
}

// CreateJournalEntry creates a new journal entry
func (h *JournalEntryHandler) CreateJournalEntry(c *gin.Context) {
	var req dto.CreateJournalEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Parse organization ID
	orgID, err := uuid.Parse(req.OrganizationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid organization ID",
			Message: err.Error(),
		})
		return
	}

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid transaction date format",
			Message: "Use YYYY-MM-DD format",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := getUserIDFromContext(c)

	// Convert request to domain model
	entry := &domain.JournalEntry{
		OrganizationID:  orgID,
		TransactionDate: transactionDate,
		Reference:       req.Reference,
		Description:     req.Description,
		CreatedBy:       userID,
		Lines:           make([]domain.JournalLine, len(req.Lines)),
	}

	for i, line := range req.Lines {
		accountID, err := uuid.Parse(line.AccountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid account ID",
				Message: err.Error(),
			})
			return
		}

		entry.Lines[i] = domain.JournalLine{
			AccountID:   accountID,
			Reference:   line.Reference,
			Description: line.Description,
			Debit:       line.Debit,
			Credit:      line.Credit,
		}
	}

	// Create entry
	createdEntry, err := h.service.CreateEntry(c.Request.Context(), entry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create journal entry",
			Message: err.Error(),
		})
		return
	}

	// Convert to response
	response := mapper.ToJournalEntryResponse(createdEntry)
	c.JSON(http.StatusCreated, response)
}

// UpdateJournalEntry updates an existing journal entry
func (h *JournalEntryHandler) UpdateJournalEntry(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid entry ID",
			Message: err.Error(),
		})
		return
	}

	var req dto.UpdateJournalEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid transaction date format",
			Message: "Use YYYY-MM-DD format",
		})
		return
	}

	// Get existing entry
	existingEntry, err := h.service.GetEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Journal entry not found",
			Message: err.Error(),
		})
		return
	}

	// Update entry fields
	existingEntry.TransactionDate = transactionDate
	existingEntry.Reference = req.Reference
	existingEntry.Description = req.Description
	existingEntry.Lines = make([]domain.JournalLine, len(req.Lines))

	for i, line := range req.Lines {
		accountID, err := uuid.Parse(line.AccountID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid account ID",
				Message: err.Error(),
			})
			return
		}

		existingEntry.Lines[i] = domain.JournalLine{
			AccountID:   accountID,
			Reference:   line.Reference,
			Description: line.Description,
			Debit:       line.Debit,
			Credit:      line.Credit,
		}
	}

	// Update entry
	updatedEntry, err := h.service.UpdateEntry(c.Request.Context(), existingEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to update journal entry",
			Message: err.Error(),
		})
		return
	}

	// Convert to response
	response := mapper.ToJournalEntryResponse(updatedEntry)
	c.JSON(http.StatusOK, response)
}

// GetJournalEntry retrieves a journal entry by ID
func (h *JournalEntryHandler) GetJournalEntry(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid entry ID",
			Message: err.Error(),
		})
		return
	}

	entry, err := h.service.GetEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Journal entry not found",
			Message: err.Error(),
		})
		return
	}

	response := mapper.ToJournalEntryResponse(entry)
	c.JSON(http.StatusOK, response)
}

// ListJournalEntries lists journal entries with pagination
func (h *JournalEntryHandler) ListJournalEntries(c *gin.Context) {
	orgID, err := uuid.Parse(c.Query("organization_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid organization ID",
			Message: err.Error(),
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.Query("status")

	var entries []*domain.JournalEntry

	if status != "" {
		entryStatus := domain.EntryStatus(status)
		entries, err = h.service.ListByStatus(c.Request.Context(), orgID, entryStatus, limit, offset)
	} else {
		entries, err = h.service.ListEntries(c.Request.Context(), orgID, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to list journal entries",
			Message: err.Error(),
		})
		return
	}

	responses := make([]dto.JournalEntryResponse, len(entries))
	for i, entry := range entries {
		responses[i] = mapper.ToJournalEntryResponse(entry)
	}

	c.JSON(http.StatusOK, responses)
}

// PostJournalEntry posts a draft entry
func (h *JournalEntryHandler) PostJournalEntry(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid entry ID",
			Message: err.Error(),
		})
		return
	}

	userID := getUserIDFromContext(c)

	err = h.service.PostEntry(c.Request.Context(), entryID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to post journal entry",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Journal entry posted successfully",
	})
}

// VoidJournalEntry voids a posted entry
func (h *JournalEntryHandler) VoidJournalEntry(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid entry ID",
			Message: err.Error(),
		})
		return
	}

	err = h.service.VoidEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to void journal entry",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Journal entry voided successfully",
	})
}

// ReverseJournalEntry creates a reversal entry
func (h *JournalEntryHandler) ReverseJournalEntry(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid entry ID",
			Message: err.Error(),
		})
		return
	}

	userID := getUserIDFromContext(c)

	reversalEntry, err := h.service.ReverseEntry(c.Request.Context(), entryID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to reverse journal entry",
			Message: err.Error(),
		})
		return
	}

	response := mapper.ToJournalEntryResponse(reversalEntry)
	c.JSON(http.StatusCreated, response)
}

// DeleteJournalEntry deletes a draft entry
func (h *JournalEntryHandler) DeleteJournalEntry(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid entry ID",
			Message: err.Error(),
		})
		return
	}

	err = h.service.DeleteEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to delete journal entry",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Journal entry deleted successfully",
	})
}

// Helper function
func getUserIDFromContext(c *gin.Context) uuid.UUID {
	// TODO: Extract user ID from JWT token or session
	// For now, return a dummy UUID
	return uuid.New()
}
