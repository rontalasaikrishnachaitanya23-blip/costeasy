package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

type OrganizationRepositoryInterface interface {
	Create(ctx context.Context, org *domain.Organization) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	GetByCode(ctx context.Context, code string) (*domain.Organization, error)
	GetByEstablishmentID(ctx context.Context, establishmentID string) (*domain.Organization, error) // ADD
	Update(ctx context.Context, org *domain.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*domain.Organization, error)
	Count(ctx context.Context) (int, error)
	CountByType(ctx context.Context, orgType domain.OrganizationType) (int, error) // ADD
	CountByEmirate(ctx context.Context, emirate domain.UAEmirate) (int, error)     // ADD
	ActivateOrganization(ctx context.Context, id uuid.UUID) error
	DeactivateOrganization(ctx context.Context, id uuid.UUID) error
	UpdateMFASettings(ctx context.Context, orgID uuid.UUID, mfaEnabled, mfaEnforced bool, mfaMethod domain.MFAMethod, allowedMethods []domain.MFAMethod) error
}
