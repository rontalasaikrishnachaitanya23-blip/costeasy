// backend/settings/internal/handler/shafafiya_handler.go
package handler

import (
	"net/http"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/chaitu35/costeasy/backend/internal/settings/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShafafiyaHandler handles HTTP requests for Shafafiya settings
type ShafafiyaHandler struct {
	service service.ShafafiyaServiceInterface
}

// NewShafafiyaHandler creates a new Shafafiya handler
func NewShafafiyaHandler(service service.ShafafiyaServiceInterface) *ShafafiyaHandler {
	return &ShafafiyaHandler{service: service}
}

// CreateShafafiyaSettingsRequest represents the request body for creating Shafafiya settings
type CreateShafafiyaSettingsRequest struct {
	OrganizationID       uuid.UUID `json:"organization_id" binding:"required"`
	Username             string    `json:"username" binding:"required"`
	Password             string    `json:"password" binding:"required"`
	ProviderCode         string    `json:"provider_code" binding:"required"`
	DefaultCurrencyCode  string    `json:"default_currency_code"`
	DefaultLanguage      string    `json:"default_language"`
	IncludeSensitiveData bool      `json:"include_sensitive_data"`
	CostingMethod        string    `json:"costing_method"`
	AllocationMethod     string    `json:"allocation_method"`
}

// UpdateCredentialsRequest represents the request body for updating credentials
type UpdateCredentialsRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	ProviderCode string `json:"provider_code" binding:"required"`
}

// UpdateCostingRequest represents the request body for updating costing configuration
type UpdateCostingRequest struct {
	CostingMethod    string `json:"costing_method" binding:"required"`
	AllocationMethod string `json:"allocation_method" binding:"required"`
}

// UpdateSubmissionRequest represents the request body for updating submission configuration
type UpdateSubmissionRequest struct {
	Language         string `json:"language" binding:"required"`
	Currency         string `json:"currency" binding:"required"`
	IncludeSensitive bool   `json:"include_sensitive"`
}

// CreateShafafiyaSettings handles POST /api/v1/organizations/:org_id/shafafiya
func (h *ShafafiyaHandler) CreateShafafiyaSettings(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req CreateShafafiyaSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings := domain.ShafafiyaOrgSettings{
		OrganizationID:       orgID,
		Username:             req.Username,
		Password:             req.Password,
		ProviderCode:         req.ProviderCode,
		DefaultCurrencyCode:  req.DefaultCurrencyCode,
		DefaultLanguage:      req.DefaultLanguage,
		IncludeSensitiveData: req.IncludeSensitiveData,
		CostingMethod:        req.CostingMethod,
		AllocationMethod:     req.AllocationMethod,
	}

	created, err := h.service.CreateShafafiyaSettings(c.Request.Context(), settings)
	if err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Shafafiya settings"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetShafafiyaSettings handles GET /api/v1/organizations/:org_id/shafafiya
func (h *ShafafiyaHandler) GetShafafiyaSettings(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	settings, err := h.service.GetShafafiyaSettings(c.Request.Context(), orgID)
	if err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Shafafiya settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateShafafiyaCredentials handles PUT /api/v1/organizations/:org_id/shafafiya/credentials
func (h *ShafafiyaHandler) UpdateShafafiyaCredentials(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateCredentialsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateShafafiyaCredentials(c.Request.Context(), orgID, req.Username, req.Password, req.ProviderCode); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credentials updated successfully"})
}

// UpdateShafafiyaCosting handles PUT /api/v1/organizations/:org_id/shafafiya/costing
func (h *ShafafiyaHandler) UpdateShafafiyaCosting(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateCostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateShafafiyaCosting(c.Request.Context(), orgID, req.CostingMethod, req.AllocationMethod); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update costing configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Costing configuration updated successfully"})
}

// UpdateShafafiyaSubmission handles PUT /api/v1/organizations/:org_id/shafafiya/submission
func (h *ShafafiyaHandler) UpdateShafafiyaSubmission(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateShafafiyaSubmission(c.Request.Context(), orgID, req.Language, req.Currency, req.IncludeSensitive); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update submission configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission configuration updated successfully"})
}

// DeleteShafafiyaSettings handles DELETE /api/v1/organizations/:org_id/shafafiya
func (h *ShafafiyaHandler) DeleteShafafiyaSettings(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.DeleteShafafiyaSettings(c.Request.Context(), orgID); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete Shafafiya settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shafafiya settings deleted successfully"})
}

// ValidateShafafiyaConfiguration handles GET /api/v1/organizations/:org_id/shafafiya/validate
func (h *ShafafiyaHandler) ValidateShafafiyaConfiguration(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.ValidateShafafiyaConfiguration(c.Request.Context(), orgID); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   domainErr.GetMessage(),
				"code":    domainErr.GetCode(),
				"isValid": false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shafafiya configuration is valid",
		"isValid": true,
	})
}

// ListFailedSubmissions handles GET /api/v1/shafafiya/failed-submissions
func (h *ShafafiyaHandler) ListFailedSubmissions(c *gin.Context) {
	limit := 100 // Default limit

	settings, err := h.service.ListFailedSubmissions(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list failed submissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"submissions": settings,
		"total":       len(settings),
	})
}

// UpdateShafafiyaSettings handles PUT /api/v1/shafafiya/organizations/:org_id
func (h *ShafafiyaHandler) UpdateShafafiyaSettings(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req CreateShafafiyaSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings := domain.ShafafiyaOrgSettings{
		OrganizationID:       orgID,
		Username:             req.Username,
		Password:             req.Password,
		ProviderCode:         req.ProviderCode,
		DefaultCurrencyCode:  req.DefaultCurrencyCode,
		DefaultLanguage:      req.DefaultLanguage,
		IncludeSensitiveData: req.IncludeSensitiveData,
		CostingMethod:        req.CostingMethod,
		AllocationMethod:     req.AllocationMethod,
	}

	updated, err := h.service.UpdateShafafiyaSettings(c.Request.Context(), orgID, settings)
	if err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Shafafiya settings"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// ValidateShafafiyaConfig handles POST /api/v1/shafafiya/validate/:org_id
func (h *ShafafiyaHandler) ValidateShafafiyaConfig(c *gin.Context) {
	orgIDParam := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.ValidateShafafiyaConfiguration(c.Request.Context(), orgID); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   domainErr.GetMessage(),
				"code":    domainErr.GetCode(),
				"isValid": false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shafafiya configuration is valid",
		"isValid": true,
	})
}

// ListShafafiyaSettings handles GET /api/v1/shafafiya/
func (h *ShafafiyaHandler) ListShafafiyaSettings(c *gin.Context) {
	// Get all Shafafiya configurations (admin only)
	settings, err := h.service.ListAllShafafiyaSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list Shafafiya settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"settings": settings,
		"total":    len(settings),
	})
}
