// backend/settings/internal/repository/organization_repository_interface.go
package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

// OrganizationRepositoryInterface defines all repository operations for Organization
type OrganizationRepositoryInterface interface {
	// CreateOrganization creates a new organization
	// Returns error if organization with same name/establishment_id already exists
	CreateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error)

	// GetOrganizationByID retrieves an organization by ID
	// Returns error if organization not found
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (domain.Organization, error)

	// GetOrganizationByEstablishmentID retrieves organization by establishment ID (Shafafiya provider code)
	// Used to find existing organizations quickly
	GetOrganizationByEstablishmentID(ctx context.Context, establishmentID string) (domain.Organization, error)

	// ListOrganizations retrieves organizations with optional filters
	// Can be filtered by type, emirate, country, active status
	// Supports pagination
	ListOrganizations(ctx context.Context, filters *OrganizationFilters) ([]domain.Organization, error)

	// UpdateOrganization updates an organization
	// Only updates: name, emirate, area, currency, description
	// Cannot update: type, country, license_number, establishment_id
	UpdateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error)

	// DeactivateOrganization soft-deactivates an organization
	// Sets is_active to false
	DeactivateOrganization(ctx context.Context, id uuid.UUID) error

	// ActivateOrganization reactivates a deactivated organization
	ActivateOrganization(ctx context.Context, id uuid.UUID) error

	// ExistsWithName checks if organization with given name exists (case-insensitive)
	// excludeID allows checking for duplicates during updates
	ExistsWithName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error)

	// ExistsWithEstablishmentID checks if organization with given establishment ID exists
	// excludeID allows checking for duplicates during updates
	ExistsWithEstablishmentID(ctx context.Context, establishmentID string, excludeID *uuid.UUID) (bool, error)

	// CountByType counts organizations by type
	CountByType(ctx context.Context, orgType domain.OrganizationType) (int, error)

	// CountByEmirate counts healthcare organizations by emirate
	CountByEmirate(ctx context.Context, emirate domain.UAEmirate) (int, error)
}
