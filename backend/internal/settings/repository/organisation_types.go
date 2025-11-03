// backend/settings/internal/repository/organization_types.go
package repository

import (
	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
)

// OrganizationFilters defines filtering options for listing organizations
type OrganizationFilters struct {
	Type     *domain.OrganizationType
	Emirate  *domain.UAEmirate
	Country  *string
	IsActive *bool
	Limit    int
	Offset   int
}

// DefaultFilters returns default filter values
func DefaultFilters() *OrganizationFilters {
	isActive := true
	return &OrganizationFilters{
		IsActive: &isActive,
		Limit:    100,
		Offset:   0,
	}
}

// WithType adds type filter
func (f *OrganizationFilters) WithType(t domain.OrganizationType) *OrganizationFilters {
	f.Type = &t
	return f
}

// WithEmirate adds emirate filter
func (f *OrganizationFilters) WithEmirate(e domain.UAEmirate) *OrganizationFilters {
	f.Emirate = &e
	return f
}

// WithCountry adds country filter
func (f *OrganizationFilters) WithCountry(c string) *OrganizationFilters {
	f.Country = &c
	return f
}

// WithActive sets active status filter
func (f *OrganizationFilters) WithActive(active bool) *OrganizationFilters {
	f.IsActive = &active
	return f
}

// WithPagination sets limit and offset
func (f *OrganizationFilters) WithPagination(limit, offset int) *OrganizationFilters {
	f.Limit = limit
	f.Offset = offset
	return f
}
