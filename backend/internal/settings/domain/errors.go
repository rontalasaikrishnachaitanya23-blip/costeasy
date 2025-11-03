// backend/settings/internal/domain/errors.go
package domain

import "fmt"

// DomainError represents a domain-level error with code for API responses
type DomainError struct {
	Message string // User-friendly message
	Code    string // Error code for programmatic handling
	Err     error  // Underlying error for debugging
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s (%s): %v", e.Message, e.Code, e.Err)
	}
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

// GetCode returns the error code
func (e *DomainError) GetCode() string {
	return e.Code
}

// GetMessage returns the error message
func (e *DomainError) GetMessage() string {
	return e.Message
}

// Unwrap returns the underlying error for error chaining
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error with message and code
func NewDomainError(message, code string) *DomainError {
	return &DomainError{
		Message: message,
		Code:    code,
		Err:     nil,
	}
}

// NewDomainErrorWithCause creates a new domain error with underlying cause
func NewDomainErrorWithCause(message, code string, cause error) *DomainError {
	return &DomainError{
		Message: message,
		Code:    code,
		Err:     cause,
	}
}

// NewDomainErrorf creates a new domain error with formatted message
func NewDomainErrorf(code, format string, args ...interface{}) *DomainError {
	return &DomainError{
		Message: fmt.Sprintf(format, args...),
		Code:    code,
		Err:     nil,
	}
}

// IsDomainError checks if an error is a DomainError
func IsDomainError(err error) bool {
	_, ok := err.(*DomainError)
	return ok
}

// AsDomainError attempts to cast error to DomainError
func AsDomainError(err error) (*DomainError, bool) {
	de, ok := err.(*DomainError)
	return de, ok
}

// ---- Error Codes ----

// Organization error codes
const (
	ErrOrgNameRequired       = "ORGANIZATION_NAME_REQUIRED"
	ErrOrgNameTooShort       = "ORGANIZATION_NAME_TOO_SHORT"
	ErrOrgNameTooLong        = "ORGANIZATION_NAME_TOO_LONG"
	ErrOrgTypeRequired       = "ORGANIZATION_TYPE_REQUIRED"
	ErrOrgTypeInvalid        = "ORGANIZATION_TYPE_INVALID"
	ErrOrgCountryRequired    = "ORGANIZATION_COUNTRY_REQUIRED"
	ErrOrgCurrencyRequired   = "ORGANIZATION_CURRENCY_REQUIRED"
	ErrOrgCurrencyInvalid    = "ORGANIZATION_CURRENCY_INVALID"
	ErrOrgEmirateRequired    = "ORGANIZATION_EMIRATE_REQUIRED"
	ErrOrgEmirateInvalid     = "ORGANIZATION_EMIRATE_INVALID"
	ErrOrgLicenseRequired    = "ORGANIZATION_LICENSE_REQUIRED"
	ErrOrgEstablishmentIDReq = "ORGANIZATION_ESTABLISHMENT_ID_REQUIRED"
	ErrOrgNotFound           = "ORGANIZATION_NOT_FOUND"
	ErrOrgAlreadyExists      = "ORGANIZATION_ALREADY_EXISTS"
	ErrOrgDeactivationFailed = "ORGANIZATION_DEACTIVATION_FAILED"
)

// Shafafiya error codes
const (
	ErrShafafiyaOrgIDRequired           = "SHAFAFIYA_ORGANIZATION_ID_REQUIRED"
	ErrShafafiyaUsernameRequired        = "SHAFAFIYA_USERNAME_REQUIRED"
	ErrShafafiyaUsernameTooShort        = "SHAFAFIYA_USERNAME_TOO_SHORT"
	ErrShafafiyaPasswordRequired        = "SHAFAFIYA_PASSWORD_REQUIRED"
	ErrShafafiyaPasswordTooShort        = "SHAFAFIYA_PASSWORD_TOO_SHORT"
	ErrShafafiyaProviderCodeRequired    = "SHAFAFIYA_PROVIDER_CODE_REQUIRED"
	ErrShafafiyaProviderCodeTooShort    = "SHAFAFIYA_PROVIDER_CODE_TOO_SHORT"
	ErrShafafiyaInvalidCurrency         = "SHAFAFIYA_INVALID_CURRENCY"
	ErrShafafiyaInvalidLanguage         = "SHAFAFIYA_INVALID_LANGUAGE"
	ErrShafafiyaInvalidCostingMethod    = "SHAFAFIYA_INVALID_COSTING_METHOD"
	ErrShafafiyaInvalidAllocationMethod = "SHAFAFIYA_INVALID_ALLOCATION_METHOD"
	ErrShafafiyaInvalidEnvironment      = "SHAFAFIYA_INVALID_ENVIRONMENT" // NEW
	ErrShafafiyaNotConfigured           = "SHAFAFIYA_NOT_CONFIGURED"
	ErrShafafiyaNotFound                = "SHAFAFIYA_NOT_FOUND"
	ErrShafafiyaOrgNotAbuDhabi          = "SHAFAFIYA_ORG_NOT_ABU_DHABI"
	ErrShafafiyaSubmissionFailed        = "SHAFAFIYA_SUBMISSION_FAILED"
)
