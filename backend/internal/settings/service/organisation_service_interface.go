package service

import (
    "context"

    "github.com/chaitu35/costeasy/backend/internal/settings/domain"
    "github.com/chaitu35/costeasy/backend/internal/settings/handler/dto"
    "github.com/google/uuid"
)

type OrganizationServiceInterface interface {
    CreateOrganization(ctx context.Context, req *dto.CreateOrganizationRequest) (*dto.OrganizationResponse, error)
    GetOrganizationByID(ctx context.Context, id uuid.UUID) (*dto.OrganizationResponse, error)
    GetOrganizationByEstablishmentID(ctx context.Context, establishmentID string) (*dto.OrganizationResponse, error)
    ListOrganizations(ctx context.Context, limit, offset int) ([]*dto.OrganizationResponse, error)
    UpdateOrganization(ctx context.Context, id uuid.UUID, req *dto.UpdateOrganizationRequest) (*dto.OrganizationResponse, error)
    ActivateOrganization(ctx context.Context, id uuid.UUID) error
    DeactivateOrganization(ctx context.Context, id uuid.UUID) error
    GetOrganizationStats(ctx context.Context) (*dto.OrganizationStatsResponse, error)
    GetOrganizationsByType(ctx context.Context, orgType domain.OrganizationType) ([]*dto.OrganizationResponse, error)
    GetOrganizationsByEmirate(ctx context.Context, emirate domain.UAEmirate) ([]*dto.OrganizationResponse, error)
}
