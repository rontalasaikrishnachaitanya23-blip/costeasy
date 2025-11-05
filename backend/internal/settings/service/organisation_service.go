package service

import (
	"context"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/chaitu35/costeasy/backend/internal/settings/handler/dto"
	"github.com/chaitu35/costeasy/backend/internal/settings/repository"
	"github.com/google/uuid"
)

type OrganizationService struct {
	repo repository.OrganizationRepositoryInterface
}

func NewOrganizationService(repo repository.OrganizationRepositoryInterface) OrganizationServiceInterface {
	return &OrganizationService{
		repo: repo,
	}
}

// CreateOrganization creates a new organization
func (s *OrganizationService) CreateOrganization(ctx context.Context, req *dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error) {
	// Validate organization name uniqueness
	if req.Code != "" {
		existing, _ := s.repo.GetByCode(ctx, req.Code)
		if existing != nil {
			return nil, domain.NewDomainError("organization with this code already exists", domain.ErrOrgAlreadyExists)
		}
	}

	// For healthcare organizations, validate establishment ID uniqueness
	if req.Type == string(domain.OrganizationTypeHealthcare) && req.EstablishmentID != "" {
		existingEst, _ := s.repo.GetByEstablishmentID(ctx, req.EstablishmentID)
		if existingEst != nil {
			return nil, domain.NewDomainError("organization with this establishment ID already exists", domain.ErrOrgEstablishmentIDExists)
		}
	}

	// Create organization domain object
	org := &domain.Organization{
		ID:              uuid.New(),
		Name:            req.Name,
		Code:            domain.StringPtr(req.Code),
		DisplayName:     domain.StringPtr(req.DisplayName),
		Type:            domain.OrganizationType(req.Type),
		Country:         req.Country,
		Emirate:         domain.EmiratePtr(req.Emirate),
		Area:            domain.StringPtr(req.Area),
		Address:         domain.StringPtr(req.Address),
		City:            domain.StringPtr(req.City),
		State:           domain.StringPtr(req.State),
		PostalCode:      domain.StringPtr(req.PostalCode),
		Phone:           domain.StringPtr(req.Phone),
		Email:           domain.StringPtr(req.Email),
		Website:         domain.StringPtr(req.Website),
		Currency:        req.Currency,
		TaxID:           domain.StringPtr(req.TaxID),
		LicenseNumber:   domain.StringPtr(req.LicenseNumber),
		EstablishmentID: domain.StringPtr(req.EstablishmentID),
		Description:     domain.StringPtr(req.Description),
		IsActive:        true,
		MFAEnabled:      false,
		MFAEnforced:     false,
		MFAMethod:       domain.MFAMethodNone,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Validate domain rules
	if err := org.Validate(); err != nil {
		return nil, err
	}

	// Create in database
	if err := s.repo.Create(ctx, org); err != nil {
		return nil, err
	}

	return dto.ToOrganizationResponse(org), nil
}

// GetOrganizationByID retrieves organization by ID
func (s *OrganizationService) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*dto.OrganizationResponse, error) {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewDomainError("organization not found", domain.ErrOrgNotFound)
	}

	return dto.ToOrganizationResponse(org), nil
}

// GetOrganizationByEstablishmentID retrieves organization by establishment ID
func (s *OrganizationService) GetOrganizationByEstablishmentID(ctx context.Context, establishmentID string) (*dto.OrganizationResponse, error) {
	org, err := s.repo.GetByEstablishmentID(ctx, establishmentID)
	if err != nil {
		return nil, domain.NewDomainError("organization not found", domain.ErrOrgNotFound)
	}

	return dto.ToOrganizationResponse(org), nil
}

// ListOrganizations retrieves all organizations with pagination
func (s *OrganizationService) ListOrganizations(ctx context.Context, limit, offset int) ([]*dto.OrganizationResponse, error) {
	orgs, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return dto.ToOrganizationListResponse(orgs), nil
}

// UpdateOrganization updates an existing organization
func (s *OrganizationService) UpdateOrganization(ctx context.Context, id uuid.UUID, req *dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error) {
	// Get existing organization
	existingOrg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.NewDomainError("organization not found", domain.ErrOrgNotFound)
	}

	// Check code uniqueness if changed
	if req.Code != "" && domain.StringValue(existingOrg.Code) != req.Code {
		existing, _ := s.repo.GetByCode(ctx, req.Code)
		if existing != nil && existing.ID != id {
			return nil, domain.NewDomainError("organization with this code already exists", domain.ErrOrgAlreadyExists)
		}
	}

	// Update fields using helper functions
	existingOrg.Name = req.Name
	existingOrg.Code = domain.StringPtr(req.Code)
	existingOrg.DisplayName = domain.StringPtr(req.DisplayName)
	existingOrg.Type = domain.OrganizationType(req.Type)
	existingOrg.Country = req.Country
	existingOrg.Emirate = domain.EmiratePtr(req.Emirate)
	existingOrg.Area = domain.StringPtr(req.Area)
	existingOrg.Address = domain.StringPtr(req.Address)
	existingOrg.City = domain.StringPtr(req.City)
	existingOrg.State = domain.StringPtr(req.State)
	existingOrg.PostalCode = domain.StringPtr(req.PostalCode)
	existingOrg.Phone = domain.StringPtr(req.Phone)
	existingOrg.Email = domain.StringPtr(req.Email)
	existingOrg.Website = domain.StringPtr(req.Website)
	existingOrg.Currency = req.Currency
	existingOrg.TaxID = domain.StringPtr(req.TaxID)
	existingOrg.LicenseNumber = domain.StringPtr(req.LicenseNumber)
	existingOrg.EstablishmentID = domain.StringPtr(req.EstablishmentID)
	existingOrg.Description = domain.StringPtr(req.Description)
	existingOrg.UpdatedAt = time.Now()

	// Validate
	if err := existingOrg.Validate(); err != nil {
		return nil, err
	}

	// Update in database
	if err := s.repo.Update(ctx, existingOrg); err != nil {
		return nil, err
	}

	return dto.ToOrganizationResponse(existingOrg), nil
}

// ActivateOrganization activates an organization
func (s *OrganizationService) ActivateOrganization(ctx context.Context, id uuid.UUID) error {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.NewDomainError("organization not found", domain.ErrOrgNotFound)
	}

	if org.IsActive {
		return domain.NewDomainError("organization is already active", domain.ErrOrgAlreadyActive)
	}

	return s.repo.ActivateOrganization(ctx, id)
}

// DeactivateOrganization deactivates an organization
func (s *OrganizationService) DeactivateOrganization(ctx context.Context, id uuid.UUID) error {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.NewDomainError("organization not found", domain.ErrOrgNotFound)
	}

	if !org.IsActive {
		return domain.NewDomainError("organization is already inactive", domain.ErrOrgAlreadyInactive)
	}

	return s.repo.DeactivateOrganization(ctx, id)
}

// GetOrganizationStats returns statistics about organizations
func (s *OrganizationService) GetOrganizationStats(ctx context.Context) (*dto.OrganizationStatsResponse, error) {
	totalCount, err := s.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	healthcareCount, _ := s.repo.CountByType(ctx, domain.OrganizationTypeHealthcare)

	// Get all organizations for emirate counting
	orgs, _ := s.repo.List(ctx, 1000, 0)

	typeMap := make(map[string]int)
	emirateMap := make(map[string]int)

	for _, org := range orgs {
		// Count by type
		typeMap[string(org.Type)]++

		// Count by emirate
		if org.Emirate != nil {
			emirateMap[string(*org.Emirate)]++
		}
	}

	return &dto.OrganizationStatsResponse{
		TotalOrganizations:      totalCount,
		ActiveOrganizations:     totalCount,
		HealthcareOrganizations: healthcareCount,
		OrganizationsByType:     typeMap,
		OrganizationsByEmirate:  emirateMap,
	}, nil
}

// GetOrganizationsByType retrieves organizations by type
func (s *OrganizationService) GetOrganizationsByType(ctx context.Context, orgType domain.OrganizationType) ([]*dto.OrganizationResponse, error) {
	orgs, err := s.repo.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	var filtered []*domain.Organization
	for _, org := range orgs {
		if org.Type == orgType {
			filtered = append(filtered, org)
		}
	}

	return dto.ToOrganizationListResponse(filtered), nil
}

// GetOrganizationsByEmirate retrieves organizations by emirate
func (s *OrganizationService) GetOrganizationsByEmirate(ctx context.Context, emirate domain.UAEmirate) ([]*dto.OrganizationResponse, error) {
	orgs, err := s.repo.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	var filtered []*domain.Organization
	for _, org := range orgs {
		if org.Emirate != nil && *org.Emirate == emirate {
			filtered = append(filtered, org)
		}
	}

	return dto.ToOrganizationListResponse(filtered), nil
}
