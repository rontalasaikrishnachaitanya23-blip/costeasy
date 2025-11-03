// backend/settings/internal/service/organization_service_interface.go
package service

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/chaitu35/costeasy/backend/internal/settings/repository"
	"github.com/google/uuid"
)

// OrganizationServiceInterface defines business operations for organizations
type OrganizationServiceInterface interface {
	// CreateOrganization creates a new organization with validation
	CreateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error)

	// GetOrganizationByID retrieves an organization by ID
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (domain.Organization, error)

	// GetOrganizationByEstablishmentID retrieves organization by establishment ID
	GetOrganizationByEstablishmentID(ctx context.Context, establishmentID string) (domain.Organization, error)

	// ListOrganizations retrieves organizations with filters
	ListOrganizations(ctx context.Context, filters *repository.OrganizationFilters) ([]domain.Organization, error)

	// UpdateOrganization updates an organization
	UpdateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error)

	// DeactivateOrganization deactivates an organization
	DeactivateOrganization(ctx context.Context, id uuid.UUID) error

	// ActivateOrganization reactivates an organization
	ActivateOrganization(ctx context.Context, id uuid.UUID) error

	// GetOrganizationStats returns statistics about organizations
	GetOrganizationStats(ctx context.Context) (*OrganizationStats, error)
}

// OrganizationStats contains organization statistics
type OrganizationStats struct {
	TotalOrganizations    int                             `json:"total_organizations"`
	ActiveOrganizations   int                             `json:"active_organizations"`
	InactiveOrganizations int                             `json:"inactive_organizations"`
	ByType                map[domain.OrganizationType]int `json:"by_type"`
	ByEmirate             map[domain.UAEmirate]int        `json:"by_emirate"`
}
