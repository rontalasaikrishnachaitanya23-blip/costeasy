package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/chaitu35/costeasy/backend/internal/settings/handler/dto"
	"github.com/chaitu35/costeasy/backend/internal/settings/service"
)

type OrganizationHandler struct {
	service service.OrganizationServiceInterface
}

func NewOrganizationHandler(service service.OrganizationServiceInterface) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
	}
}

// CreateOrganization creates a new organization
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var req dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.service.CreateOrganization(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// GetOrganization retrieves organization by ID
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	org, err := h.service.GetOrganizationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// GetOrganizationByEstablishmentID retrieves organization by establishment ID
func (h *OrganizationHandler) GetOrganizationByEstablishmentID(c *gin.Context) {
	establishmentID := c.Param("establishmentId")

	org, err := h.service.GetOrganizationByEstablishmentID(c.Request.Context(), establishmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// ListOrganizations retrieves all organizations with pagination
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	limit := 50
	offset := 0

	// Parse query parameters if provided
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	orgs, err := h.service.ListOrganizations(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orgs)
}

// UpdateOrganization updates an existing organization
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	var req dto.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.service.UpdateOrganization(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// ActivateOrganization activates an organization
func (h *OrganizationHandler) ActivateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	if err := h.service.ActivateOrganization(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organization activated successfully"})
}

// DeactivateOrganization deactivates an organization
func (h *OrganizationHandler) DeactivateOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization ID"})
		return
	}

	if err := h.service.DeactivateOrganization(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "organization deactivated successfully"})
}

// GetOrganizationStats retrieves organization statistics
func (h *OrganizationHandler) GetOrganizationStats(c *gin.Context) {
	stats, err := h.service.GetOrganizationStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// OrganizationStatsResponse represents organization statistics
type OrganizationStatsResponse struct {
	TotalOrganizations      int            `json:"total_organizations"`
	ActiveOrganizations     int            `json:"active_organizations"`
	InactiveOrganizations   int            `json:"inactive_organizations"`
	OrganizationsByType     map[string]int `json:"organizations_by_type"`
	OrganizationsByEmirate  map[string]int `json:"organizations_by_emirate"`
	HealthcareOrganizations int            `json:"healthcare_organizations"`
}
