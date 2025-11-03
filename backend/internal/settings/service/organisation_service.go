// backend/settings/internal/service/organization_service.go
package service

import (
	"context"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/chaitu35/costeasy/backend/internal/settings/repository"
	"github.com/google/uuid"
)

// OrganizationService implements OrganizationServiceInterface
type OrganizationService struct {
	repo repository.OrganizationRepositoryInterface
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(repo repository.OrganizationRepositoryInterface) *OrganizationService {
	return &OrganizationService{repo: repo}
}

// CreateOrganization creates a new organization with validation
func (s *OrganizationService) CreateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error) {
	// Domain validation
	if err := org.Validate(); err != nil {
		return domain.Organization{}, err
	}

	// Business rule: Check for duplicate name
	exists, err := s.repo.ExistsWithName(ctx, org.Name, nil)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to check organization name: %w", err)
	}
	if exists {
		return domain.Organization{}, domain.NewDomainError(
			fmt.Sprintf("organization with name '%s' already exists", org.Name),
			domain.ErrOrgAlreadyExists,
		)
	}

	// Business rule: Check for duplicate establishment ID (if provided)
	if org.EstablishmentID != "" {
		exists, err := s.repo.ExistsWithEstablishmentID(ctx, org.EstablishmentID, nil)
		if err != nil {
			return domain.Organization{}, fmt.Errorf("failed to check establishment ID: %w", err)
		}
		if exists {
			return domain.Organization{}, domain.NewDomainError(
				fmt.Sprintf("organization with establishment ID '%s' already exists", org.EstablishmentID),
				domain.ErrOrgAlreadyExists,
			)
		}
	}

	// Create organization
	created, err := s.repo.CreateOrganization(ctx, org)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to create organization: %w", err)
	}

	return created, nil
}

// GetOrganizationByID retrieves an organization by ID
func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id uuid.UUID) (domain.Organization, error) {
	if id == uuid.Nil {
		return domain.Organization{}, domain.NewDomainError("organization ID is required", domain.ErrOrgNotFound)
	}

	org, err := s.repo.GetOrganizationByID(ctx, id)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("organization not found: %w", err)
	}

	return org, nil
}

// GetOrganizationByEstablishmentID retrieves organization by establishment ID
func (s *OrganizationService) GetOrganizationByEstablishmentID(ctx context.Context, establishmentID string) (domain.Organization, error) {
	if establishmentID == "" {
		return domain.Organization{}, domain.NewDomainError("establishment ID is required", domain.ErrOrgNotFound)
	}

	org, err := s.repo.GetOrganizationByEstablishmentID(ctx, establishmentID)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("organization not found: %w", err)
	}

	return org, nil
}

// ListOrganizations retrieves organizations with filters
func (s *OrganizationService) ListOrganizations(ctx context.Context, filters *repository.OrganizationFilters) ([]domain.Organization, error) {
	if filters == nil {
		filters = repository.DefaultFilters()
	}

	orgs, err := s.repo.ListOrganizations(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return orgs, nil
}

// UpdateOrganization updates an organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, org domain.Organization) (domain.Organization, error) {
	if org.ID == uuid.Nil {
		return domain.Organization{}, domain.NewDomainError("organization ID is required", domain.ErrOrgNotFound)
	}

	// Verify organization exists
	existing, err := s.repo.GetOrganizationByID(ctx, org.ID)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("organization not found: %w", err)
	}

	// Preserve immutable fields
	org.Type = existing.Type
	org.Country = existing.Country
	org.TaxID = existing.TaxID
	org.LicenseNumber = existing.LicenseNumber
	org.EstablishmentID = existing.EstablishmentID
	org.CreatedAt = existing.CreatedAt

	// Domain validation
	if err := org.Validate(); err != nil {
		return domain.Organization{}, err
	}

	// Business rule: Check for duplicate name (excluding current org)
	exists, err := s.repo.ExistsWithName(ctx, org.Name, &org.ID)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to check organization name: %w", err)
	}
	if exists {
		return domain.Organization{}, domain.NewDomainError(
			fmt.Sprintf("organization with name '%s' already exists", org.Name),
			domain.ErrOrgAlreadyExists,
		)
	}

	// Update organization
	updated, err := s.repo.UpdateOrganization(ctx, org)
	if err != nil {
		return domain.Organization{}, fmt.Errorf("failed to update organization: %w", err)
	}

	return updated, nil
}

// DeactivateOrganization deactivates an organization
func (s *OrganizationService) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.NewDomainError("organization ID is required", domain.ErrOrgNotFound)
	}

	// Verify organization exists
	_, err := s.repo.GetOrganizationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	// TODO: Add business rule checks
	// - Can't deactivate if there are active accounts
	// - Can't deactivate if there are pending transactions
	// - etc.

	if err := s.repo.DeactivateOrganization(ctx, id); err != nil {
		return fmt.Errorf("failed to deactivate organization: %w", err)
	}

	return nil
}

// ActivateOrganization reactivates an organization
func (s *OrganizationService) ActivateOrganization(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return domain.NewDomainError("organization ID is required", domain.ErrOrgNotFound)
	}

	// Verify organization exists
	_, err := s.repo.GetOrganizationByID(ctx, id)
	if err != nil {
		return fmt.Errorf("organization not found: %w", err)
	}

	if err := s.repo.ActivateOrganization(ctx, id); err != nil {
		return fmt.Errorf("failed to activate organization: %w", err)
	}

	return nil
}

// GetOrganizationStats returns statistics about organizations
func (s *OrganizationService) GetOrganizationStats(ctx context.Context) (*OrganizationStats, error) {
	stats := &OrganizationStats{
		ByType:    make(map[domain.OrganizationType]int),
		ByEmirate: make(map[domain.UAEmirate]int),
	}

	// Get all active organizations
	activeFilters := repository.DefaultFilters().WithActive(true)
	activeOrgs, err := s.repo.ListOrganizations(ctx, activeFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get active organizations: %w", err)
	}
	stats.ActiveOrganizations = len(activeOrgs)

	// Get all inactive organizations
	inactiveFilters := repository.DefaultFilters().WithActive(false)
	inactiveOrgs, err := s.repo.ListOrganizations(ctx, inactiveFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive organizations: %w", err)
	}
	stats.InactiveOrganizations = len(inactiveOrgs)

	stats.TotalOrganizations = stats.ActiveOrganizations + stats.InactiveOrganizations

	// Count by type
	orgTypes := []domain.OrganizationType{
		domain.OrganizationTypeHealthcare,
		domain.OrganizationTypeRetail,
		domain.OrganizationTypeManufacturing,
		domain.OrganizationTypeFinance,
		domain.OrganizationTypeEducation,
		domain.OrganizationTypeHospitality,
		domain.OrganizationTypeLogistics,
		domain.OrganizationTypeRealEstate,
		domain.OrganizationTypeService,
		domain.OrganizationTypeOther,
	}

	for _, orgType := range orgTypes {
		count, err := s.repo.CountByType(ctx, orgType)
		if err != nil {
			return nil, fmt.Errorf("failed to count by type: %w", err)
		}
		if count > 0 {
			stats.ByType[orgType] = count
		}
	}

	// Count by emirate (healthcare only)
	emirates := []domain.UAEmirate{
		domain.EmirateAbuDhabi,
		domain.EmirateDubai,
		domain.EmirateSharjah,
		domain.EmirateRasAlKhaimah,
		domain.EmirateUmAlQuwain,
		domain.EmirateFujairah,
		domain.EmirateAjman,
	}

	for _, emirate := range emirates {
		count, err := s.repo.CountByEmirate(ctx, emirate)
		if err != nil {
			return nil, fmt.Errorf("failed to count by emirate: %w", err)
		}
		if count > 0 {
			stats.ByEmirate[emirate] = count
		}
	}

	return stats, nil
}
