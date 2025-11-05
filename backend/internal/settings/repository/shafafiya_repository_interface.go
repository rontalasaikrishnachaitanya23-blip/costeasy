// backend/settings/internal/repository/shafafiya_repository_interface.go
package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

// ShafafiyaSettingsRepositoryInterface defines all repository operations for ShafafiyaOrgSettings
type ShafafiyaSettingsRepositoryInterface interface {
	CreateShafafiyaSettings(ctx context.Context, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error)
	GetShafafiyaSettings(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error)
	GetShafafiyaSettingsWithPassword(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error)
	UpdateShafafiyaSettings(ctx context.Context, orgID uuid.UUID, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error) // ADD THIS
	UpdateShafafiyaCredentials(ctx context.Context, orgID uuid.UUID, username, encryptedPassword, providerCode string) error
	UpdateShafafiyaCosting(ctx context.Context, orgID uuid.UUID, costingMethod, allocationMethod string) error
	UpdateShafafiyaSubmission(ctx context.Context, orgID uuid.UUID, language, currency string, includeSensitive bool) error
	UpdateSubmissionStatus(ctx context.Context, orgID uuid.UUID, status, errorMsg string) error
	DeleteShafafiyaSettings(ctx context.Context, orgID uuid.UUID) error
	ExistsForOrganization(ctx context.Context, orgID uuid.UUID) (bool, error)
	ListBySubmissionStatus(ctx context.Context, status string, limit int) ([]domain.ShafafiyaOrgSettings, error)
	List(ctx context.Context, limit, offset int) ([]domain.ShafafiyaOrgSettings, error) // ADD THIS
}
