// backend/settings/internal/service/shafafiya_service_interface.go
package service

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

// ShafafiyaServiceInterface defines business operations for Shafafiya settings
type ShafafiyaServiceInterface interface {
	CreateShafafiyaSettings(ctx context.Context, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error)
	GetShafafiyaSettings(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error)
	UpdateShafafiyaSettings(ctx context.Context, orgID uuid.UUID, settings domain.ShafafiyaOrgSettings) (*domain.ShafafiyaOrgSettings, error) // ADD THIS
	UpdateShafafiyaCredentials(ctx context.Context, orgID uuid.UUID, username, password, providerCode string) error
	UpdateShafafiyaCosting(ctx context.Context, orgID uuid.UUID, costingMethod, allocationMethod string) error
	UpdateShafafiyaSubmission(ctx context.Context, orgID uuid.UUID, language, currency string, includeSensitive bool) error
	UpdateSubmissionStatus(ctx context.Context, orgID uuid.UUID, status, errorMsg string) error
	DeleteShafafiyaSettings(ctx context.Context, orgID uuid.UUID) error
	ValidateShafafiyaConfiguration(ctx context.Context, orgID uuid.UUID) error
	GetDecryptedCredentials(ctx context.Context, orgID uuid.UUID) (*ShafafiyaCredentials, error)
	ListFailedSubmissions(ctx context.Context, limit int) ([]domain.ShafafiyaOrgSettings, error)
	ListAllShafafiyaSettings(ctx context.Context) ([]domain.ShafafiyaOrgSettings, error) // ADD THIS
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
