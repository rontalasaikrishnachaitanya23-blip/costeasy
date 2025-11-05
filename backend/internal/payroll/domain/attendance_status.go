package domain

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Code           string    `json:"code"`
	Description    string    `json:"description"`
	IsPaid         bool      `json:"is_paid"`
	IsWorkingDay   bool      `json:"is_working_day"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
