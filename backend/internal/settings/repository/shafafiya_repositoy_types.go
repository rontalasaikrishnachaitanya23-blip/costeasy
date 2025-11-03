// backend/settings/internal/repository/shafafiya_types.go
package repository

import (
	"github.com/chaitu35/costeasy/backend/internal/settings/domain"
)

// ShafafiyaFilters defines filtering options for listing Shafafiya settings
type ShafafiyaFilters struct {
	SubmissionStatus *string
	Emirate          *domain.UAEmirate
	Limit            int
	Offset           int
}

// DefaultShafafiyaFilters returns default filter values
func DefaultShafafiyaFilters() *ShafafiyaFilters {
	return &ShafafiyaFilters{
		Limit:  100,
		Offset: 0,
	}
}

// WithSubmissionStatus adds submission status filter
func (f *ShafafiyaFilters) WithSubmissionStatus(status string) *ShafafiyaFilters {
	f.SubmissionStatus = &status
	return f
}

// WithEmirate adds emirate filter
func (f *ShafafiyaFilters) WithEmirate(e domain.UAEmirate) *ShafafiyaFilters {
	f.Emirate = &e
	return f
}

// WithPagination sets limit and offset
func (f *ShafafiyaFilters) WithPagination(limit, offset int) *ShafafiyaFilters {
	f.Limit = limit
	f.Offset = offset
	return f
}
