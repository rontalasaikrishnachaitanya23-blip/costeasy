package dto

import (
	"time"

	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
	"github.com/google/uuid"
)

// CreateOrganizationRequest represents the request to create an organization
type CreateOrganizationRequest struct {
	Name            string `json:"name" binding:"required"`
	Code            string `json:"code"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type" binding:"required"`
	Country         string `json:"country" binding:"required"`
	Emirate         string `json:"emirate"`
	Area            string `json:"area"`
	Address         string `json:"address"`
	City            string `json:"city"`
	State           string `json:"state"`
	PostalCode      string `json:"postal_code"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Website         string `json:"website"`
	Currency        string `json:"currency" binding:"required"`
	TaxID           string `json:"tax_id"`
	LicenseNumber   string `json:"license_number"`
	LicenseExpiry   string `json:"license_expiry"`
	EstablishmentID string `json:"establishment_id"`
	Description     string `json:"description"`
}

// UpdateOrganizationRequest represents the request to update an organization
type UpdateOrganizationRequest struct {
	Name            string `json:"name" binding:"required"`
	Code            string `json:"code"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type" binding:"required"`
	Country         string `json:"country" binding:"required"`
	Emirate         string `json:"emirate"`
	Area            string `json:"area"`
	Address         string `json:"address"`
	City            string `json:"city"`
	State           string `json:"state"`
	PostalCode      string `json:"postal_code"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Website         string `json:"website"`
	Currency        string `json:"currency" binding:"required"`
	TaxID           string `json:"tax_id"`
	LicenseNumber   string `json:"license_number"`
	LicenseExpiry   string `json:"license_expiry"`
	EstablishmentID string `json:"establishment_id"`
	Description     string `json:"description"`
}

// OrganizationResponse represents the response for organization operations
type OrganizationResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Code            string     `json:"code"`
	DisplayName     string     `json:"display_name"`
	Type            string     `json:"type"`
	Country         string     `json:"country"`
	Emirate         string     `json:"emirate"`
	Area            string     `json:"area"`
	Address         string     `json:"address"`
	City            string     `json:"city"`
	State           string     `json:"state"`
	PostalCode      string     `json:"postal_code"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	Website         string     `json:"website"`
	Currency        string     `json:"currency"`
	TaxID           string     `json:"tax_id"`
	LicenseNumber   string     `json:"license_number"`
	LicenseExpiry   *time.Time `json:"license_expiry,omitempty"`
	EstablishmentID string     `json:"establishment_id"`
	Description     string     `json:"description"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// ToOrganizationResponse converts domain organization to response DTO
func ToOrganizationResponse(org *domain.Organization) *OrganizationResponse {
	return &OrganizationResponse{
		ID:              org.ID,
		Name:            org.Name,
		Code:            domain.StringValue(org.Code),
		DisplayName:     domain.StringValue(org.DisplayName),
		Type:            string(org.Type),
		Country:         org.Country,
		Emirate:         domain.EmirateValue(org.Emirate),
		Area:            domain.StringValue(org.Area),
		Address:         domain.StringValue(org.Address),
		City:            domain.StringValue(org.City),
		State:           domain.StringValue(org.State),
		PostalCode:      domain.StringValue(org.PostalCode),
		Phone:           domain.StringValue(org.Phone),
		Email:           domain.StringValue(org.Email),
		Website:         domain.StringValue(org.Website),
		Currency:        org.Currency,
		TaxID:           domain.StringValue(org.TaxID),
		LicenseNumber:   domain.StringValue(org.LicenseNumber),
		LicenseExpiry:   org.LicenseExpiry,
		EstablishmentID: domain.StringValue(org.EstablishmentID),
		Description:     domain.StringValue(org.Description),
		IsActive:        org.IsActive,
		CreatedAt:       org.CreatedAt,
		UpdatedAt:       org.UpdatedAt,
	}
}

// ToOrganizationListResponse converts multiple organizations to response DTOs
func ToOrganizationListResponse(orgs []*domain.Organization) []*OrganizationResponse {
	responses := make([]*OrganizationResponse, len(orgs))
	for i, org := range orgs {
		responses[i] = ToOrganizationResponse(org)
	}
	return responses
}

// OrganizationStatsResponse represents organization statistics
type OrganizationStatsResponse struct {
	TotalOrganizations      int            `json:"total_organizations"`
	ActiveOrganizations     int            `json:"active_organizations"`
	InactiveOrganizations   int            `json:"inactive_organizations"`
	OrganizationsByType     map[string]int `json:"organizations_by_type"`
	OrganizationsByEmirate  map[string]int `json:"organizations_by_emirate"`
	HealthcareOrganizations int            `json:"healthcare_organizations"`
}
