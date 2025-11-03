// backend/settings/internal/domain/organization.go (Updated Validate method)
package domain

import (
	"time"

	"github.com/google/uuid"
)

// OrganizationType defines the type of organization/business
type OrganizationType string

const (
	OrganizationTypeHealthcare    OrganizationType = "HEALTHCARE"
	OrganizationTypeRetail        OrganizationType = "RETAIL"
	OrganizationTypeManufacturing OrganizationType = "MANUFACTURING"
	OrganizationTypeFinance       OrganizationType = "FINANCE"
	OrganizationTypeEducation     OrganizationType = "EDUCATION"
	OrganizationTypeHospitality   OrganizationType = "HOSPITALITY"
	OrganizationTypeLogistics     OrganizationType = "LOGISTICS"
	OrganizationTypeRealEstate    OrganizationType = "REAL_ESTATE"
	OrganizationTypeService       OrganizationType = "SERVICE"
	OrganizationTypeOther         OrganizationType = "OTHER"
)

// UAEmirate defines UAE emirates
type UAEmirate string

const (
	EmirateAbuDhabi     UAEmirate = "ABU_DHABI"
	EmirateDubai        UAEmirate = "DUBAI"
	EmirateSharjah      UAEmirate = "SHARJAH"
	EmirateRasAlKhaimah UAEmirate = "RAS_AL_KHAIMAH"
	EmirateUmAlQuwain   UAEmirate = "UMM_AL_QUWAIN"
	EmirateFujairah     UAEmirate = "FUJAIRAH"
	EmirateAjman        UAEmirate = "AJMAN"
)

// Organization represents a business organization
type Organization struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Type            OrganizationType `json:"type"`
	Country         string           `json:"country"`
	Emirate         UAEmirate        `json:"emirate"`
	Area            string           `json:"area"`
	Currency        string           `json:"currency"`
	TaxID           string           `json:"tax_id"`
	LicenseNumber   string           `json:"license_number"`
	EstablishmentID string           `json:"establishment_id"`
	Description     string           `json:"description"`
	IsActive        bool             `json:"is_active"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// Validate performs domain validation on Organization
func (o *Organization) Validate() *DomainError {
	// Name validation
	if o.Name == "" {
		return NewDomainError("name is required", ErrOrgNameRequired)
	}

	if len(o.Name) < 3 {
		return NewDomainErrorf(ErrOrgNameTooShort, "name must be at least 3 characters, got %d", len(o.Name))
	}

	if len(o.Name) > 255 {
		return NewDomainErrorf(ErrOrgNameTooLong, "name cannot exceed 255 characters, got %d", len(o.Name))
	}

	// Type validation
	if o.Type == "" {
		return NewDomainError("organization type is required", ErrOrgTypeRequired)
	}

	if !isValidOrganizationType(o.Type) {
		return NewDomainErrorf(ErrOrgTypeInvalid, "invalid organization type: %s", o.Type)
	}

	// Country validation
	if o.Country == "" {
		return NewDomainError("country is required", ErrOrgCountryRequired)
	}

	// Currency validation
	if o.Currency == "" {
		return NewDomainError("currency is required", ErrOrgCurrencyRequired)
	}

	if !IsValidCurrencyCode(o.Currency) {
		return NewDomainErrorf(ErrOrgCurrencyInvalid, "invalid currency code: %s", o.Currency)
	}

	// Healthcare-specific validations
	if o.Type == OrganizationTypeHealthcare {
		if err := o.validateHealthcare(); err != nil {
			return err
		}
	}

	return nil
}

// validateHealthcare performs healthcare-specific validation
func (o *Organization) validateHealthcare() *DomainError {
	if o.Emirate == "" {
		return NewDomainError("emirate is required for healthcare organizations", ErrOrgEmirateRequired)
	}

	if !isValidEmirate(o.Emirate) {
		return NewDomainErrorf(ErrOrgEmirateInvalid, "invalid emirate: %s", o.Emirate)
	}

	if o.LicenseNumber == "" {
		return NewDomainError("license number is required for healthcare organizations", ErrOrgLicenseRequired)
	}

	if o.EstablishmentID == "" {
		return NewDomainError("establishment ID (Shafafiya provider code) is required for healthcare organizations", ErrOrgEstablishmentIDReq)
	}

	return nil
}

// IsHealthcare returns true if organization is healthcare type
func (o *Organization) IsHealthcare() bool {
	return o.Type == OrganizationTypeHealthcare
}

// IsInAbuDhabi returns true if organization is in Abu Dhabi emirate
func (o *Organization) IsInAbuDhabi() bool {
	return o.Emirate == EmirateAbuDhabi
}

// Helper functions
func isValidOrganizationType(t OrganizationType) bool {
	validTypes := map[OrganizationType]bool{
		OrganizationTypeHealthcare:    true,
		OrganizationTypeRetail:        true,
		OrganizationTypeManufacturing: true,
		OrganizationTypeFinance:       true,
		OrganizationTypeEducation:     true,
		OrganizationTypeHospitality:   true,
		OrganizationTypeLogistics:     true,
		OrganizationTypeRealEstate:    true,
		OrganizationTypeService:       true,
		OrganizationTypeOther:         true,
	}
	return validTypes[t]
}

func isValidEmirate(e UAEmirate) bool {
	validEmirates := map[UAEmirate]bool{
		EmirateAbuDhabi:     true,
		EmirateDubai:        true,
		EmirateSharjah:      true,
		EmirateRasAlKhaimah: true,
		EmirateUmAlQuwain:   true,
		EmirateFujairah:     true,
		EmirateAjman:        true,
	}
	return validEmirates[e]
}
