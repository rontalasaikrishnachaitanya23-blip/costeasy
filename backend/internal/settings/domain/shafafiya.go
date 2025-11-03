// backend/settings/internal/domain/shafafiya.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// ShafafiyaOrgSettings stores Shafafiya-specific configuration for Abu Dhabi healthcare facilities
type ShafafiyaOrgSettings struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`

	// Shafafiya Credentials
	Username     string `json:"username"`           // Shafafiya portal username
	Password     string `json:"password,omitempty"` // Encrypted in storage, never returned to client
	ProviderCode string `json:"provider_code"`      // Provider code assigned by Shafafiya

	// XML Submission Settings
	DefaultCurrencyCode  string `json:"default_currency_code"`  // Default: "AED"
	DefaultLanguage      string `json:"default_language"`       // Default: "en", Alternative: "ar"
	IncludeSensitiveData bool   `json:"include_sensitive_data"` // Whether to include patient sensitive data in XML

	// Costing Configuration
	CostingMethod    string `json:"costing_method"`    // e.g., "DEPARTMENTAL", "ACTIVITY_BASED"
	AllocationMethod string `json:"allocation_method"` // e.g., "WEIGHTED", "PERCENTAGE"

	// Submission Tracking
	LastSubmissionAt     *time.Time `json:"last_submission_at,omitempty"`
	LastSubmissionStatus string     `json:"last_submission_status,omitempty"` // SUCCESS, FAILURE, PENDING
	LastSubmissionError  string     `json:"last_submission_error,omitempty"`  // Error message if failed

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate performs domain validation on ShafafiyaOrgSettings
func (s *ShafafiyaOrgSettings) Validate() *DomainError {
	// Organization ID validation
	if s.OrganizationID == uuid.Nil {
		return NewDomainError("organization ID is required", ErrShafafiyaOrgIDRequired)
	}

	// Username validation
	if s.Username == "" {
		return NewDomainError("Shafafiya username is required", ErrShafafiyaUsernameRequired)
	}

	if len(s.Username) < 3 {
		return NewDomainErrorf(ErrShafafiyaUsernameTooShort, "username must be at least 3 characters, got %d", len(s.Username))
	}

	if len(s.Username) > 100 {
		return NewDomainErrorf(ErrShafafiyaUsernameTooShort, "username cannot exceed 100 characters, got %d", len(s.Username))
	}

	// Password validation
	if s.Password == "" {
		return NewDomainError("Shafafiya password is required", ErrShafafiyaPasswordRequired)
	}

	if len(s.Password) < 6 {
		return NewDomainErrorf(ErrShafafiyaPasswordTooShort, "password must be at least 6 characters, got %d", len(s.Password))
	}

	if len(s.Password) > 255 {
		return NewDomainErrorf(ErrShafafiyaPasswordTooShort, "password cannot exceed 255 characters")
	}

	// Provider Code validation
	if s.ProviderCode == "" {
		return NewDomainError("provider code is required", ErrShafafiyaProviderCodeRequired)
	}

	if len(s.ProviderCode) < 3 {
		return NewDomainErrorf(ErrShafafiyaProviderCodeTooShort, "provider code must be at least 3 characters, got %d", len(s.ProviderCode))
	}

	if len(s.ProviderCode) > 100 {
		return NewDomainErrorf(ErrShafafiyaProviderCodeTooShort, "provider code cannot exceed 100 characters")
	}

	// Currency code validation
	if s.DefaultCurrencyCode == "" {
		s.DefaultCurrencyCode = "AED"
	} else if !IsValidCurrencyCode(s.DefaultCurrencyCode) {
		return NewDomainErrorf(ErrShafafiyaInvalidCurrency, "invalid currency code: %s", s.DefaultCurrencyCode)
	}

	// Language validation
	if s.DefaultLanguage == "" {
		s.DefaultLanguage = "en"
	} else if !isValidLanguage(s.DefaultLanguage) {
		return NewDomainErrorf(ErrShafafiyaInvalidLanguage, "invalid language code: %s, must be 'en' or 'ar'", s.DefaultLanguage)
	}

	// Costing method validation
	if s.CostingMethod == "" {
		s.CostingMethod = "DEPARTMENTAL"
	} else if !isValidCostingMethod(s.CostingMethod) {
		return NewDomainErrorf(ErrShafafiyaInvalidCostingMethod, "invalid costing method: %s", s.CostingMethod)
	}

	// Allocation method validation
	if s.AllocationMethod == "" {
		s.AllocationMethod = "WEIGHTED"
	} else if !isValidAllocationMethod(s.AllocationMethod) {
		return NewDomainErrorf(ErrShafafiyaInvalidAllocationMethod, "invalid allocation method: %s", s.AllocationMethod)
	}

	return nil
}

// IsConfigured returns true if all required credentials are present
func (s *ShafafiyaOrgSettings) IsConfigured() bool {
	return s.Username != "" && s.Password != "" && s.ProviderCode != ""
}

// CanSubmit returns true if settings are valid and configured for submission
func (s *ShafafiyaOrgSettings) CanSubmit() bool {
	return s.IsConfigured() && s.CostingMethod != "" && s.DefaultCurrencyCode != ""
}

// UpdateSubmissionStatus updates the submission tracking information
func (s *ShafafiyaOrgSettings) UpdateSubmissionStatus(status string, errorMsg string) *DomainError {
	if !isValidSubmissionStatus(status) {
		return NewDomainErrorf(ErrShafafiyaSubmissionFailed, "invalid submission status: %s", status)
	}

	now := time.Now()
	s.LastSubmissionAt = &now
	s.LastSubmissionStatus = status

	if errorMsg != "" {
		s.LastSubmissionError = errorMsg
	} else {
		s.LastSubmissionError = ""
	}

	return nil
}

// LastSubmissionWasSuccessful returns true if the last submission was successful
func (s *ShafafiyaOrgSettings) LastSubmissionWasSuccessful() bool {
	return s.LastSubmissionStatus == "SUCCESS"
}

// LastSubmissionWasFailed returns true if the last submission failed
func (s *ShafafiyaOrgSettings) LastSubmissionWasFailed() bool {
	return s.LastSubmissionStatus == "FAILURE"
}

// Helper functions for Shafafiya validation

// IsValidCurrencyCode validates currency codes
func IsValidCurrencyCode(code string) bool {
	validCurrencies := map[string]bool{
		"AED": true, // UAE Dirham
		"USD": true, // US Dollar
		"EUR": true, // Euro
		"GBP": true, // British Pound
		"INR": true, // Indian Rupee
		"SAR": true, // Saudi Riyal
		"KWD": true, // Kuwaiti Dinar
		"QAR": true, // Qatari Riyal
	}
	return validCurrencies[code]
}

// isValidLanguage validates language codes for XML submission
func isValidLanguage(lang string) bool {
	validLanguages := map[string]bool{
		"en": true, // English
		"ar": true, // Arabic
	}
	return validLanguages[lang]
}

// isValidCostingMethod validates costing methodology
func isValidCostingMethod(method string) bool {
	validMethods := map[string]bool{
		"DEPARTMENTAL":   true, // By department
		"ACTIVITY_BASED": true, // Activity-based costing
		"SERVICE_BASED":  true, // By service line
	}
	return validMethods[method]
}

// isValidAllocationMethod validates overhead allocation methods
func isValidAllocationMethod(method string) bool {
	validMethods := map[string]bool{
		"WEIGHTED":   true, // Weighted allocation
		"PERCENTAGE": true, // Percentage-based
		"FIXED":      true, // Fixed amount
	}
	return validMethods[method]
}

// isValidSubmissionStatus validates submission status values
func isValidSubmissionStatus(status string) bool {
	validStatuses := map[string]bool{
		"SUCCESS": true, // Submission successful
		"FAILURE": true, // Submission failed
		"PENDING": true, // Submission pending
	}
	return validStatuses[status]
}

// ---- Endpoint Methods ----

// GetCostingSubmissionEndpoint returns the SOAP endpoint for Shafafiya web services
func (s *ShafafiyaOrgSettings) GetCostingSubmissionEndpoint(environment string) string {
	switch environment {
	case EnvironmentProduction:
		return "https://shafafiya.doh.gov.ae/v3/webservices.asmx"
	case EnvironmentStaging, EnvironmentTest:
		return "https://shafafiyapte.doh.gov.ae/v3/webservices.asmx"
	default:
		// Default to staging for safety
		return "https://shafafiyapte.doh.gov.ae/v3/webservices.asmx"
	}
}

// GetWSDL returns the WSDL URL for the web service
func (s *ShafafiyaOrgSettings) GetWSDL(environment string) string {
	return s.GetCostingSubmissionEndpoint(environment) + "?WSDL"
}

// GetSOAPNamespace returns the SOAP namespace for XML
func (s *ShafafiyaOrgSettings) GetSOAPNamespace() string {
	return "https://www.shafafiya.org/v2"
}

// GetSOAPEnvelopeNamespace returns the SOAP envelope namespace
func (s *ShafafiyaOrgSettings) GetSOAPEnvelopeNamespace() string {
	return "http://www.w3.org/2003/05/soap-envelope"
}

// GetSOAPVersion returns the SOAP version (SOAP 1.2)
func (s *ShafafiyaOrgSettings) GetSOAPVersion() string {
	return "1.2"
}

// GetDataDictionaryNamespace returns the namespace for Common Types
func (s *ShafafiyaOrgSettings) GetDataDictionaryNamespace() string {
	return "http://www.haad.ae/DataDictionary/CommonTypes"
}

// GetVersion returns the API version being used
func (s *ShafafiyaOrgSettings) GetVersion() string {
	return "v3" // Latest version
}

// GetLegacyVersion returns the legacy API version (v2)
func (s *ShafafiyaOrgSettings) GetLegacyVersion() string {
	return "v2"
}

// GetLegacyEndpoint returns the v2 endpoint (for backward compatibility)
func (s *ShafafiyaOrgSettings) GetLegacyEndpoint(environment string) string {
	switch environment {
	case EnvironmentProduction:
		return "https://shafafiya.haad.ae/v2/webservices.asmx"
	case EnvironmentStaging, EnvironmentTest:
		return "https://shafafiyapte.doh.gov.ae/v2/webservices.asmx"
	default:
		return "https://shafafiyapte.doh.gov.ae/v2/webservices.asmx"
	}
}

// GetTimeout returns the recommended timeout for web service calls (in seconds)
func (s *ShafafiyaOrgSettings) GetTimeout() int {
	return 30 // 30 seconds
}

// GetMaxRetries returns the maximum number of retries for failed requests
func (s *ShafafiyaOrgSettings) GetMaxRetries() int {
	return 3
}

// GetRetryDelay returns the delay between retries (in seconds)
func (s *ShafafiyaOrgSettings) GetRetryDelay() int {
	return 5
}
