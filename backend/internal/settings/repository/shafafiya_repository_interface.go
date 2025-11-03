// backend/settings/internal/repository/shafafiya_repository_interface.go
package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

// ShafafiyaSettingsRepositoryInterface defines all repository operations for ShafafiyaOrgSettings
type ShafafiyaSettingsRepositoryInterface interface {
	// CreateShafafiyaSettings creates new Shafafiya settings for an organization
	// Returns error if settings already exist for the organization
	CreateShafafiyaSettings(ctx context.Context, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error)

	// GetShafafiyaSettings retrieves Shafafiya settings by organization ID
	// Does NOT return encrypted password
	GetShafafiyaSettings(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error)

	// GetShafafiyaSettingsWithPassword retrieves Shafafiya settings WITH encrypted password
	// For internal use only (e.g., XML submission)
	// NEVER return this to API client
	GetShafafiyaSettingsWithPassword(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error)

	// UpdateShafafiyaCredentials updates only the credentials (username, password, provider_code)
	UpdateShafafiyaCredentials(ctx context.Context, orgID uuid.UUID, username, encryptedPassword, providerCode string) error

	// UpdateShafafiyaCosting updates costing configuration
	UpdateShafafiyaCosting(ctx context.Context, orgID uuid.UUID, costingMethod, allocationMethod string) error

	// UpdateShafafiyaSubmission updates submission configuration (language, currency, sensitive data flag)
	UpdateShafafiyaSubmission(ctx context.Context, orgID uuid.UUID, language, currency string, includeSensitive bool) error

	// UpdateSubmissionStatus updates submission tracking (last submission time, status, error)
	UpdateSubmissionStatus(ctx context.Context, orgID uuid.UUID, status, errorMsg string) error

	// DeleteShafafiyaSettings deletes Shafafiya settings for an organization
	DeleteShafafiyaSettings(ctx context.Context, orgID uuid.UUID) error

	// ExistsForOrganization checks if Shafafiya settings exist for an organization
	ExistsForOrganization(ctx context.Context, orgID uuid.UUID) (bool, error)

	// ListBySubmissionStatus retrieves all settings with specific submission status
	// Useful for tracking failed submissions
	ListBySubmissionStatus(ctx context.Context, status string, limit int) ([]domain.ShafafiyaOrgSettings, error)
}
