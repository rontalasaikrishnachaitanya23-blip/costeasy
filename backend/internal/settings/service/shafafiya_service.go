// backend/settings/internal/service/shafafiya_service.go
package service

import (
    "context"
    "fmt"

    "github.com/chaitu35/costeasy/backend/internal/settings/domain"
    "github.com/chaitu35/costeasy/backend/internal/settings/repository"
    "github.com/chaitu35/costeasy/backend/pkg/crypto"
    "github.com/google/uuid"
)



// Config holds application configuration
type Config struct {
    ShafafiyaEnvironment string // "production", "staging", or "test"
}

// ShafafiyaService implements ShafafiyaServiceInterface
type ShafafiyaService struct {
    repo          repository.ShafafiyaSettingsRepositoryInterface
    orgRepo       repository.OrganizationRepositoryInterface
    cryptoService *crypto.CryptoService
    config        *Config
}

// NewShafafiyaService creates a new Shafafiya service
func NewShafafiyaService(
    repo repository.ShafafiyaSettingsRepositoryInterface,
    orgRepo repository.OrganizationRepositoryInterface,
    cryptoService *crypto.CryptoService,
    config *Config,
) *ShafafiyaService {
    return &ShafafiyaService{
        repo:          repo,
        orgRepo:       orgRepo,
        cryptoService: cryptoService,
        config:        config,
    }
}

// CreateShafafiyaSettings creates new Shafafiya settings
func (s *ShafafiyaService) CreateShafafiyaSettings(ctx context.Context, settings domain.ShafafiyaOrgSettings) (domain.ShafafiyaOrgSettings, error) {
    // Domain validation
    if err := settings.Validate(); err != nil {
        return domain.ShafafiyaOrgSettings{}, err
    }

    // Business rule: Verify organization exists
    org, err := s.orgRepo.GetOrganizationByID(ctx, settings.OrganizationID)
    if err != nil {
        return domain.ShafafiyaOrgSettings{}, fmt.Errorf("organization not found: %w", err)
    }

    // Business rule: Must be healthcare organization
    if !org.IsHealthcare() {
        return domain.ShafafiyaOrgSettings{}, domain.NewDomainError(
            "Shafafiya settings can only be created for healthcare organizations",
            domain.ErrShafafiyaNotConfigured,
        )
    }

    // Business rule: Must be in Abu Dhabi
    if !org.IsInAbuDhabi() {
        return domain.ShafafiyaOrgSettings{}, domain.NewDomainError(
            "Shafafiya is only available for Abu Dhabi healthcare facilities",
            domain.ErrShafafiyaOrgNotAbuDhabi,
        )
    }

    // Business rule: Check if settings already exist
    exists, err := s.repo.ExistsForOrganization(ctx, settings.OrganizationID)
    if err != nil {
        return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to check existing settings: %w", err)
    }
    if exists {
        return domain.ShafafiyaOrgSettings{}, domain.NewDomainError(
            "Shafafiya settings already exist for this organization",
            domain.ErrShafafiyaNotConfigured,
        )
    }

    // Encrypt password
    encryptedPassword, err := s.cryptoService.Encrypt(settings.Password)
    if err != nil {
        return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to encrypt password: %w", err)
    }
    settings.Password = encryptedPassword

    // Set defaults if not provided
    if settings.DefaultCurrencyCode == "" {
        settings.DefaultCurrencyCode = "AED"
    }
    if settings.DefaultLanguage == "" {
        settings.DefaultLanguage = "en"
    }
    if settings.CostingMethod == "" {
        settings.CostingMethod = domain.CostingMethodDepartmental
    }
    if settings.AllocationMethod == "" {
        settings.AllocationMethod = domain.AllocationMethodWeighted
    }

    // Create settings
    created, err := s.repo.CreateShafafiyaSettings(ctx, settings)
    if err != nil {
        return domain.ShafafiyaOrgSettings{}, fmt.Errorf("failed to create Shafafiya settings: %w", err)
    }

    return created, nil
}

// GetShafafiyaSettings retrieves settings (without decrypted password)
func (s *ShafafiyaService) GetShafafiyaSettings(ctx context.Context, orgID uuid.UUID) (domain.ShafafiyaOrgSettings, error) {
    if orgID == uuid.Nil {
        return domain.ShafafiyaOrgSettings{}, domain.NewDomainError(
            "organization ID is required",
            domain.ErrShafafiyaNotFound,
        )
    }

    settings, err := s.repo.GetShafafiyaSettings(ctx, orgID)
    if err != nil {
        return domain.ShafafiyaOrgSettings{}, fmt.Errorf("shafafiya settings not found: %w", err)
    }

    return settings, nil
}

// UpdateShafafiyaCredentials updates credentials with encryption
func (s *ShafafiyaService) UpdateShafafiyaCredentials(ctx context.Context, orgID uuid.UUID, username, password, providerCode string) error {
    if orgID == uuid.Nil {
        return domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Validation
    if username == "" {
        return domain.NewDomainError("username is required", domain.ErrShafafiyaUsernameRequired)
    }
    if password == "" {
        return domain.NewDomainError("password is required", domain.ErrShafafiyaPasswordRequired)
    }
    if providerCode == "" {
        return domain.NewDomainError("provider code is required", domain.ErrShafafiyaProviderCodeRequired)
    }

    // Verify settings exist
    _, err := s.repo.GetShafafiyaSettings(ctx, orgID)
    if err != nil {
        return fmt.Errorf("shafafiya settings not found: %w", err)
    }

    // Encrypt password
    encryptedPassword, err := s.cryptoService.Encrypt(password)
    if err != nil {
        return fmt.Errorf("failed to encrypt password: %w", err)
    }

    // Update credentials
    if err := s.repo.UpdateShafafiyaCredentials(ctx, orgID, username, encryptedPassword, providerCode); err != nil {
        return fmt.Errorf("failed to update credentials: %w", err)
    }

    return nil
}

// UpdateShafafiyaCosting updates costing configuration
func (s *ShafafiyaService) UpdateShafafiyaCosting(ctx context.Context, orgID uuid.UUID, costingMethod, allocationMethod string) error {
    if orgID == uuid.Nil {
        return domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Validate costing method
    validCostingMethods := map[string]bool{
        domain.CostingMethodDepartmental:  true,
        domain.CostingMethodActivityBased: true,
        domain.CostingMethodServiceBased:  true,
    }
    if !validCostingMethods[costingMethod] {
        return domain.NewDomainError(
            fmt.Sprintf("invalid costing method: %s", costingMethod),
            domain.ErrShafafiyaInvalidCostingMethod,
        )
    }

    // Validate allocation method
    validAllocationMethods := map[string]bool{
        domain.AllocationMethodWeighted:   true,
        domain.AllocationMethodPercentage: true,
        domain.AllocationMethodFixed:      true,
    }
    if !validAllocationMethods[allocationMethod] {
        return domain.NewDomainError(
            fmt.Sprintf("invalid allocation method: %s", allocationMethod),
            domain.ErrShafafiyaInvalidAllocationMethod,
        )
    }

    // Update costing configuration
    if err := s.repo.UpdateShafafiyaCosting(ctx, orgID, costingMethod, allocationMethod); err != nil {
        return fmt.Errorf("failed to update costing configuration: %w", err)
    }

    return nil
}

// UpdateShafafiyaSubmission updates submission configuration
func (s *ShafafiyaService) UpdateShafafiyaSubmission(ctx context.Context, orgID uuid.UUID, language, currency string, includeSensitive bool) error {
    if orgID == uuid.Nil {
        return domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Validate language
    if language != "en" && language != "ar" {
        return domain.NewDomainError(
            fmt.Sprintf("invalid language: %s, must be 'en' or 'ar'", language),
            domain.ErrShafafiyaInvalidLanguage,
        )
    }

    // Validate currency
    validCurrencies := map[string]bool{
        "AED": true, "USD": true, "EUR": true, "GBP": true,
        "INR": true, "SAR": true, "KWD": true, "QAR": true,
    }
    if !validCurrencies[currency] {
        return domain.NewDomainError(
            fmt.Sprintf("invalid currency: %s", currency),
            domain.ErrShafafiyaInvalidCurrency,
        )
    }

    // Update submission configuration
    if err := s.repo.UpdateShafafiyaSubmission(ctx, orgID, language, currency, includeSensitive); err != nil {
        return fmt.Errorf("failed to update submission configuration: %w", err)
    }

    return nil
}

// UpdateSubmissionStatus updates submission tracking
func (s *ShafafiyaService) UpdateSubmissionStatus(ctx context.Context, orgID uuid.UUID, status, errorMsg string) error {
    if orgID == uuid.Nil {
        return domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Validate status
    validStatuses := map[string]bool{
        domain.SubmissionStatusSuccess: true,
        domain.SubmissionStatusFailure: true,
        domain.SubmissionStatusPending: true,
    }
    if !validStatuses[status] {
        return domain.NewDomainError(
            fmt.Sprintf("invalid submission status: %s", status),
            domain.ErrShafafiyaSubmissionFailed,
        )
    }

    // Update submission status
    if err := s.repo.UpdateSubmissionStatus(ctx, orgID, status, errorMsg); err != nil {
        return fmt.Errorf("failed to update submission status: %w", err)
    }

    return nil
}

// DeleteShafafiyaSettings deletes settings for an organization
func (s *ShafafiyaService) DeleteShafafiyaSettings(ctx context.Context, orgID uuid.UUID) error {
    if orgID == uuid.Nil {
        return domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Verify settings exist
    _, err := s.repo.GetShafafiyaSettings(ctx, orgID)
    if err != nil {
        return fmt.Errorf("shafafiya settings not found: %w", err)
    }

    // Delete settings
    if err := s.repo.DeleteShafafiyaSettings(ctx, orgID); err != nil {
        return fmt.Errorf("failed to delete Shafafiya settings: %w", err)
    }

    return nil
}

// ValidateShafafiyaConfiguration validates that settings are complete
func (s *ShafafiyaService) ValidateShafafiyaConfiguration(ctx context.Context, orgID uuid.UUID) error {
    if orgID == uuid.Nil {
        return domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Get settings
    settings, err := s.repo.GetShafafiyaSettings(ctx, orgID)
    if err != nil {
        return fmt.Errorf("shafafiya settings not found: %w", err)
    }

    // Validate organization is healthcare and Abu Dhabi
    org, err := s.orgRepo.GetOrganizationByID(ctx, orgID)
    if err != nil {
        return fmt.Errorf("organization not found: %w", err)
    }

    if !org.IsHealthcare() {
        return domain.NewDomainError(
            "organization is not a healthcare facility",
            domain.ErrShafafiyaOrgNotAbuDhabi,
        )
    }

    if !org.IsInAbuDhabi() {
        return domain.NewDomainError(
            "organization is not in Abu Dhabi",
            domain.ErrShafafiyaOrgNotAbuDhabi,
        )
    }

    // Validate settings are configured
    if !settings.IsConfigured() {
        return domain.NewDomainError(
            "Shafafiya settings are incomplete",
            domain.ErrShafafiyaNotConfigured,
        )
    }

    // Validate can submit
    if !settings.CanSubmit() {
        return domain.NewDomainError(
            "Shafafiya settings are not ready for submission",
            domain.ErrShafafiyaNotConfigured,
        )
    }

    return nil
}

// GetDecryptedCredentials returns decrypted credentials with environment-aware endpoint (INTERNAL USE ONLY)
func (s *ShafafiyaService) GetDecryptedCredentials(ctx context.Context, orgID uuid.UUID) (*ShafafiyaCredentials, error) {
    if orgID == uuid.Nil {
        return nil, domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Get settings WITH encrypted password
    settings, err := s.repo.GetShafafiyaSettingsWithPassword(ctx, orgID)
    if err != nil {
        return nil, fmt.Errorf("shafafiya settings not found: %w", err)
    }

    // Decrypt password
    decryptedPassword, err := s.cryptoService.Decrypt(settings.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt password: %w", err)
    }

    // Get environment from config (default to staging for safety)
    environment := s.getEnvironment()

    return &ShafafiyaCredentials{
        Username:     settings.Username,
        Password:     decryptedPassword,
        ProviderCode: settings.ProviderCode,
        Endpoint:     settings.GetCostingSubmissionEndpoint(environment),
        WSDLUrl:      settings.GetWSDL(environment),
        Environment:  environment,
    }, nil
}

// GetDecryptedCredentialsForEnvironment allows explicit environment override
func (s *ShafafiyaService) GetDecryptedCredentialsForEnvironment(ctx context.Context, orgID uuid.UUID, environment string) (*ShafafiyaCredentials, error) {
    if orgID == uuid.Nil {
        return nil, domain.NewDomainError("organization ID is required", domain.ErrShafafiyaNotFound)
    }

    // Validate environment
    if !isValidEnvironment(environment) {
        return nil, domain.NewDomainErrorf(domain.ErrShafafiyaInvalidEnvironment, "invalid environment: %s", environment)
    }

    // Get settings WITH encrypted password
    settings, err := s.repo.GetShafafiyaSettingsWithPassword(ctx, orgID)
    if err != nil {
        return nil, fmt.Errorf("shafafiya settings not found: %w", err)
    }

    // Decrypt password
    decryptedPassword, err := s.cryptoService.Decrypt(settings.Password)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt password: %w", err)
    }

    return &ShafafiyaCredentials{
        Username:     settings.Username,
        Password:     decryptedPassword,
        ProviderCode: settings.ProviderCode,
        Endpoint:     settings.GetCostingSubmissionEndpoint(environment),
        WSDLUrl:      settings.GetWSDL(environment),
        Environment:  environment,
    }, nil
}

// GetSOAPConfig returns complete SOAP configuration for making web service calls
func (s *ShafafiyaService) GetSOAPConfig(ctx context.Context, orgID uuid.UUID) (*ShafafiyaSOAPConfig, error) {
    credentials, err := s.GetDecryptedCredentials(ctx, orgID)
    if err != nil {
        return nil, err
    }

    settings, err := s.repo.GetShafafiyaSettings(ctx, orgID)
    if err != nil {
        return nil, fmt.Errorf("failed to get settings: %w", err)
    }

    return &ShafafiyaSOAPConfig{
        Credentials:             *credentials,
        SOAPNamespace:           settings.GetSOAPNamespace(),
        SOAPEnvelopeNamespace:   settings.GetSOAPEnvelopeNamespace(),
        SOAPVersion:             settings.GetSOAPVersion(),
        DataDictionaryNamespace: settings.GetDataDictionaryNamespace(),
        Timeout:                 settings.GetTimeout(),
        MaxRetries:              settings.GetMaxRetries(),
        RetryDelay:              settings.GetRetryDelay(),
    }, nil
}

// GetSOAPConfigForEnvironment returns SOAP config for specific environment
func (s *ShafafiyaService) GetSOAPConfigForEnvironment(ctx context.Context, orgID uuid.UUID, environment string) (*ShafafiyaSOAPConfig, error) {
    credentials, err := s.GetDecryptedCredentialsForEnvironment(ctx, orgID, environment)
    if err != nil {
        return nil, err
    }

    settings, err := s.repo.GetShafafiyaSettings(ctx, orgID)
    if err != nil {
        return nil, fmt.Errorf("failed to get settings: %w", err)
    }

    return &ShafafiyaSOAPConfig{
        Credentials:             *credentials,
        SOAPNamespace:           settings.GetSOAPNamespace(),
        SOAPEnvelopeNamespace:   settings.GetSOAPEnvelopeNamespace(),
        SOAPVersion:             settings.GetSOAPVersion(),
        DataDictionaryNamespace: settings.GetDataDictionaryNamespace(),
        Timeout:                 settings.GetTimeout(),
        MaxRetries:              settings.GetMaxRetries(),
        RetryDelay:              settings.GetRetryDelay(),
    }, nil
}

// TestConnection tests the connection to Shafafiya web service
func (s *ShafafiyaService) TestConnection(ctx context.Context, orgID uuid.UUID) error {
    credentials, err := s.GetDecryptedCredentials(ctx, orgID)
    if err != nil {
        return err
    }

    // Validate credentials are complete
    if credentials.Username == "" || credentials.Password == "" {
        return domain.NewDomainError("invalid credentials", domain.ErrShafafiyaPasswordRequired)
    }

    // TODO: Implement actual SOAP connection test
    // For now, just validate that we can get credentials
    return nil
}

// ListFailedSubmissions retrieves organizations with failed submissions
func (s *ShafafiyaService) ListFailedSubmissions(ctx context.Context, limit int) ([]domain.ShafafiyaOrgSettings, error) {
    if limit <= 0 {
        limit = 100
    }

    settings, err := s.repo.ListBySubmissionStatus(ctx, domain.SubmissionStatusFailure, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to list failed submissions: %w", err)
    }

    return settings, nil
}

// GetAvailableOperations returns list of available Shafafiya operations
func (s *ShafafiyaService) GetAvailableOperations() []string {
    return []string{
        domain.OperationUploadTransaction,
        domain.OperationGetNewTransactions,
        domain.OperationDownloadTransactionFile,
        domain.OperationSetTransactionDownloaded,
        domain.OperationSearchTransactions,
        domain.OperationAddDRGToEClaim,
        domain.OperationGetDRGDetails,
        domain.OperationCheckForNewPriorAuthorizationTransactions,
        domain.OperationGetNewPriorAuthorizationTransactions,
        domain.OperationGetPrescriptions,
        domain.OperationGetPersonInsuranceHistory,
        domain.OperationGetClaimCountReconciliation,
        domain.OperationGetInsuranceContinuityCertificate,
        domain.OperationCancelInsuranceContinuityCertificate,
    }
}

// GetEndpointInfo returns information about configured endpoints
func (s *ShafafiyaService) GetEndpointInfo() map[string]string {
    settings := &domain.ShafafiyaOrgSettings{}

    return map[string]string{
        "production_v3": settings.GetCostingSubmissionEndpoint(domain.EnvironmentProduction),
        "staging_v3":    settings.GetCostingSubmissionEndpoint(domain.EnvironmentStaging),
        "production_v2": settings.GetLegacyEndpoint(domain.EnvironmentProduction),
        "staging_v2":    settings.GetLegacyEndpoint(domain.EnvironmentStaging),
    }
}

// getEnvironment retrieves the current environment from config
func (s *ShafafiyaService) getEnvironment() string {
    if s.config != nil && s.config.ShafafiyaEnvironment != "" {
        return s.config.ShafafiyaEnvironment
    }

    // Default to staging for safety
    return domain.EnvironmentStaging
}

// isValidEnvironment validates the environment string
func isValidEnvironment(env string) bool {
    validEnvs := map[string]bool{
        domain.EnvironmentProduction: true,
        domain.EnvironmentStaging:    true,
        domain.EnvironmentTest:       true,
    }
    return validEnvs[env]
}
