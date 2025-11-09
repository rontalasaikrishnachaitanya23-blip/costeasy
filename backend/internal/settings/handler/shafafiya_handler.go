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

func NewShafafiyaHandler(service service.ShafafiyaServiceInterface) *ShafafiyaHandler {
	return &ShafafiyaHandler{service: service}
}

// Request DTOs ─────────────────────────────

type CreateShafafiyaSettingsRequest struct {
	Username             string `json:"username" binding:"required"`
	Password             string `json:"password" binding:"required"`
	ProviderCode         string `json:"provider_code" binding:"required"`
	DefaultCurrencyCode  string `json:"default_currency_code"`
	DefaultLanguage      string `json:"default_language"`
	IncludeSensitiveData bool   `json:"include_sensitive_data"`
	CostingMethod        string `json:"costing_method"`
	AllocationMethod     string `json:"allocation_method"`
}

type UpdateCredentialsRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	ProviderCode string `json:"provider_code" binding:"required"`
}

type UpdateCostingRequest struct {
	CostingMethod    string `json:"costing_method" binding:"required"`
	AllocationMethod string `json:"allocation_method" binding:"required"`
}

type UpdateSubmissionRequest struct {
	Language         string `json:"language" binding:"required"`
	Currency         string `json:"currency" binding:"required"`
	IncludeSensitive bool   `json:"include_sensitive"`
}

// ─────────────────────────────────────────────
// @Tags Settings / Shafafiya
// @Summary Create new Shafafiya settings for organization
// @Router /api/v1/settings/shafafiya/org/{org_id} [post]
func (h *ShafafiyaHandler) CreateShafafiyaSettings(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// ─────────────────────────────────────────────
// @Summary Get Shafafiya settings for an organization
// @Router /api/v1/settings/shafafiya/org/{org_id} [get]
func (h *ShafafiyaHandler) GetShafafiyaSettings(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	settings, err := h.service.GetShafafiyaSettings(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// ─────────────────────────────────────────────
// @Summary Update Shafafiya settings
// @Router /api/v1/settings/shafafiya/org/{org_id} [put]
func (h *ShafafiyaHandler) UpdateShafafiyaSettings(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// ─────────────────────────────────────────────
// @Summary Update credentials only
// @Router /api/v1/settings/shafafiya/org/{org_id}/credentials [put]
func (h *ShafafiyaHandler) UpdateShafafiyaCredentials(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Credentials updated successfully"})
}

// ─────────────────────────────────────────────
// @Summary Update costing/allocation settings
// @Router /api/v1/settings/shafafiya/org/{org_id}/costing [put]
func (h *ShafafiyaHandler) UpdateShafafiyaCosting(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Costing updated successfully"})
}

// ─────────────────────────────────────────────
// @Summary Update submission preferences
// @Router /api/v1/settings/shafafiya/org/{org_id}/submission [put]
func (h *ShafafiyaHandler) UpdateShafafiyaSubmission(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Submission updated successfully"})
}

// ─────────────────────────────────────────────
// @Summary Delete Shafafiya settings
// @Router /api/v1/settings/shafafiya/org/{org_id} [delete]
func (h *ShafafiyaHandler) DeleteShafafiyaSettings(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.DeleteShafafiyaSettings(c.Request.Context(), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shafafiya settings deleted successfully"})
}

// ─────────────────────────────────────────────
// @Summary Validate configuration
// @Router /api/v1/settings/shafafiya/org/{org_id}/validate [get]
func (h *ShafafiyaHandler) ValidateShafafiyaConfiguration(c *gin.Context) {
	orgID, err := uuid.Parse(c.Param("org_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.ValidateShafafiyaConfiguration(c.Request.Context(), orgID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"isValid": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"isValid": true, "message": "Configuration is valid"})
}

// ─────────────────────────────────────────────
// @Summary List all Shafafiya configurations (admin)
// @Router /api/v1/settings/shafafiya/list [get]
func (h *ShafafiyaHandler) ListShafafiyaSettings(c *gin.Context) {
	settings, err := h.service.ListAllShafafiyaSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings, "count": len(settings)})
}

// ─────────────────────────────────────────────
// @Summary List failed submissions
// @Router /api/v1/settings/shafafiya/failed-submissions [get]
func (h *ShafafiyaHandler) ListFailedSubmissions(c *gin.Context) {
	records, err := h.service.ListFailedSubmissions(c.Request.Context(), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"records": records, "count": len(records)})
}
