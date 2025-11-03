// backend/settings/internal/service/shafafiya_service_interface.go
package service

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

// ShafafiyaServiceInterface defines business operations for Shafafiya settings
type ShafafiyaServiceInterface interface {
	// CreateShafafiyaSettings creates new Shafafiya settings for an organization
	// Validates that organization is healthcare and in Abu Dhabi
	CreateShafafiyaSettings(ctx context.Context, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error)

	// GetShafafiyaSettings retrieves Shafafiya settings by organization ID
	// Does NOT return decrypted password
	GetShafafiyaSettings(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error)

	// UpdateShafafiyaCredentials updates credentials (username, password, provider code)
	// Password will be encrypted before storage
	UpdateShafafiyaCredentials(ctx context.Context, orgID uuid.UUID, username, password, providerCode string) error

	// UpdateShafafiyaCosting updates costing configuration
	UpdateShafafiyaCosting(ctx context.Context, orgID uuid.UUID, costingMethod, allocationMethod string) error

	// UpdateShafafiyaSubmission updates submission configuration
	UpdateShafafiyaSubmission(ctx context.Context, orgID uuid.UUID, language, currency string, includeSensitive bool) error

	// UpdateSubmissionStatus updates submission tracking information
	UpdateSubmissionStatus(ctx context.Context, orgID uuid.UUID, status, errorMsg string) error

	// DeleteShafafiyaSettings deletes settings for an organization
	DeleteShafafiyaSettings(ctx context.Context, orgID uuid.UUID) error

	// ValidateShafafiyaConfiguration validates that settings are complete and correct
	ValidateShafafiyaConfiguration(ctx context.Context, orgID uuid.UUID) error

	// GetDecryptedCredentials returns decrypted credentials for SOAP client (internal use)
	// NEVER expose this to API - only for internal submission service
	GetDecryptedCredentials(ctx context.Context, orgID uuid.UUID) (*ShafafiyaCredentials, error)

	// ListFailedSubmissions retrieves organizations with failed submissions
	ListFailedSubmissions(ctx context.Context, limit int) ([]domain.ShafafiyaOrgSettings, error)
}

// ShafafiyaCredentials holds decrypted credentials for API calls
type ShafafiyaCredentials struct {
	Username     string
	Password     string
	ProviderCode string
	Endpoint     string
	WSDLUrl      string
	Environment  string
}

// ShafafiyaSOAPConfig holds complete SOAP configuration
type ShafafiyaSOAPConfig struct {
	Credentials             ShafafiyaCredentials
	SOAPNamespace           string
	SOAPEnvelopeNamespace   string
	SOAPVersion             string
	DataDictionaryNamespace string
	Timeout                 int
	MaxRetries              int
	RetryDelay              int
}
