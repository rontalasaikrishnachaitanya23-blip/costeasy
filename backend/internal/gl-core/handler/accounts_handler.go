// backend/internal/gl-core/handler/account_handler.go
package handler

import (
	"net/http"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler/dto"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler/mapper"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/repository"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	service service.AccountServiceInterface
}

func NewAccountHandler(service service.AccountServiceInterface) *AccountHandler {
	return &AccountHandler{service: service}
}

// CreateAccount handles POST /accounts
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	account := domain.GLAccount{
		Code:       req.Code,
		Name:       req.Name,
		Type:       req.Type,
		ParentCode: req.ParentCode,
	}

	created, err := h.service.CreateAccount(c.Request.Context(), account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create account",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, mapper.ToAccountResponse(created))
}

// GetAccount handles GET /accounts/:id
func (h *AccountHandler) GetAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	includeInactive := c.DefaultQuery("include_inactive", "false") == "true"

	account, err := h.service.GetAccountByID(c.Request.Context(), id, !includeInactive)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Account not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, mapper.ToAccountResponse(account))
}

// GetAccountByCode handles GET /accounts/code/:code
func (h *AccountHandler) GetAccountByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Account code is required",
			Message: "Code parameter cannot be empty",
		})
		return
	}

	includeInactive := c.DefaultQuery("include_inactive", "false") == "true"

	account, err := h.service.GetAccountByCode(c.Request.Context(), code, !includeInactive)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Account not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, mapper.ToAccountResponse(account))
}

// ListAccounts handles GET /accounts
func (h *AccountHandler) ListAccounts(c *gin.Context) {
	includeInactive := c.DefaultQuery("include_inactive", "false") == "true"

	accounts, err := h.service.ListAccounts(c.Request.Context(), !includeInactive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to list accounts",
			Message: err.Error(),
		})
		return
	}

	response := make([]dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = mapper.ToAccountResponse(account)
	}

	c.JSON(http.StatusOK, response)
}

// UpdateAccount handles PUT /accounts/:id
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// Get existing account to preserve immutable fields
	existing, err := h.service.GetAccountByID(c.Request.Context(), id, false)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "Account not found",
			Message: err.Error(),
		})
		return
	}

	account := domain.GLAccount{
		ID:         id,
		Code:       existing.Code, // Code cannot be changed
		Name:       req.Name,
		Type:       existing.Type, // Type cannot be changed
		ParentCode: req.ParentCode,
	}

	updated, err := h.service.UpdateAccount(c.Request.Context(), account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to update account",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, mapper.ToAccountResponse(updated))
}

// DeactivateAccount handles POST /accounts/:id/deactivate
func (h *AccountHandler) DeactivateAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	err = h.service.DeactivateAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to deactivate account",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Account deactivated successfully",
	})
}

// ActivateAccount handles POST /accounts/:id/activate
func (h *AccountHandler) ActivateAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	err = h.service.ActivateAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Failed to activate account",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Account activated successfully",
	})
}

// SoftDeleteAccount handles DELETE /accounts/:id
func (h *AccountHandler) SoftDeleteAccount(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid account ID",
			Message: err.Error(),
		})
		return
	}

	err = h.service.SoftDeleteAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to delete account",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "Account deleted successfully",
	})
}

// SearchAccounts handles GET /accounts/search
func (h *AccountHandler) SearchAccounts(c *gin.Context) {
	params := repository.GLAccountSearchParams{
		Name: c.Query("name"),
	}

	// Parse account type if provided
	if typeStr := c.Query("type"); typeStr != "" {
		accountType := domain.AccountType(typeStr)
		params.Type = &accountType
	}

	// Parse is_active if provided
	if activeStr := c.Query("is_active"); activeStr != "" {
		isActive := activeStr == "true"
		params.IsActive = &isActive
	}

	accounts, err := h.service.SearchAccounts(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to search accounts",
			Message: err.Error(),
		})
		return
	}

	response := make([]dto.AccountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = mapper.ToAccountResponse(account)
	}

	c.JSON(http.StatusOK, response)
}
