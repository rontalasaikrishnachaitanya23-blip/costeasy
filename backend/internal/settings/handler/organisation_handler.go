// backend/settings/internal/handler/organization_handler.go
package handler

import (
	"net/http"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/chaitu35/costeasy/backend/internal/settings/repository"
	"github.com/chaitu35/costeasy/backend/internal/settings/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OrganizationHandler handles HTTP requests for organizations
type OrganizationHandler struct {
	service service.OrganizationServiceInterface
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(service service.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

// CreateOrganizationRequest represents the request body for creating an organization
type CreateOrganizationRequest struct {
	Name            string                  `json:"name" binding:"required"`
	Type            domain.OrganizationType `json:"type" binding:"required"`
	Country         string                  `json:"country" binding:"required"`
	Emirate         domain.UAEmirate        `json:"emirate"`
	Area            string                  `json:"area"`
	Currency        string                  `json:"currency" binding:"required"`
	TaxID           string                  `json:"tax_id"`
	LicenseNumber   string                  `json:"license_number"`
	EstablishmentID string                  `json:"establishment_id"`
	Description     string                  `json:"description"`
}

// UpdateOrganizationRequest represents the request body for updating an organization
type UpdateOrganizationRequest struct {
	Name        string           `json:"name" binding:"required"`
	Emirate     domain.UAEmirate `json:"emirate"`
	Area        string           `json:"area"`
	Currency    string           `json:"currency" binding:"required"`
	Description string           `json:"description"`
}

// CreateOrganization handles POST /api/v1/organizations
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := domain.Organization{
		Name:            req.Name,
		Type:            req.Type,
		Country:         req.Country,
		Emirate:         req.Emirate,
		Area:            req.Area,
		Currency:        req.Currency,
		TaxID:           req.TaxID,
		LicenseNumber:   req.LicenseNumber,
		EstablishmentID: req.EstablishmentID,
		Description:     req.Description,
	}

	created, err := h.service.CreateOrganization(c.Request.Context(), org)
	if err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetOrganization handles GET /api/v1/organizations/:id
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := h.service.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusNotFound, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// ListOrganizations handles GET /api/v1/organizations
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	filters := repository.DefaultFilters()

	// Parse query parameters
	if orgType := c.Query("type"); orgType != "" {
		t := domain.OrganizationType(orgType)
		filters.Type = &t
	}

	if emirate := c.Query("emirate"); emirate != "" {
		e := domain.UAEmirate(emirate)
		filters.Emirate = &e
	}

	if country := c.Query("country"); country != "" {
		filters.Country = &country
	}

	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filters.IsActive = &active
	}

	orgs, err := h.service.ListOrganizations(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list organizations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"organizations": orgs,
		"total":         len(orgs),
	})
}

// UpdateOrganization handles PUT /api/v1/organizations/:id
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := domain.Organization{
		ID:          id,
		Name:        req.Name,
		Emirate:     req.Emirate,
		Area:        req.Area,
		Currency:    req.Currency,
		Description: req.Description,
	}

	updated, err := h.service.UpdateOrganization(c.Request.Context(), org)
	if err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeactivateOrganization handles DELETE /api/v1/organizations/:id
func (h *OrganizationHandler) DeactivateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.DeactivateOrganization(c.Request.Context(), id); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deactivated successfully"})
}

// ActivateOrganization handles POST /api/v1/organizations/:id/activate
func (h *OrganizationHandler) ActivateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := h.service.ActivateOrganization(c.Request.Context(), id); err != nil {
		if domainErr, ok := domain.AsDomainError(err); ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": domainErr.GetMessage(),
				"code":  domainErr.GetCode(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate organization"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization activated successfully"})
}

// GetOrganizationStats handles GET /api/v1/organizations/stats
func (h *OrganizationHandler) GetOrganizationStats(c *gin.Context) {
	stats, err := h.service.GetOrganizationStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get organization statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
