package domain

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceUploadBatch tracks Excel/manual attendance imports
type AttendanceUploadBatch struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	UploadedBy     uuid.UUID  `json:"uploaded_by"`
	FileName       string     `json:"file_name"`
	TotalRecords   int        `json:"total_records"`
	ProcessedCount int        `json:"processed_count"`
	Status         string     `json:"status"` // PENDING, PROCESSED, FAILED
	UploadedAt     time.Time  `json:"uploaded_at"`
	ProcessedAt    *time.Time `json:"processed_at,omitempty"`
	Remarks        *string    `json:"remarks,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
