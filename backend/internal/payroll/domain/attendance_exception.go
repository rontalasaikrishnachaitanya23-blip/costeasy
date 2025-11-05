package domain

import (
	"time"
	"github.com/google/uuid"
)

// AttendanceException records manual HR overrides
type AttendanceException struct {
	ID             uuid.UUID  `json:"id"`
	EmployeeID     uuid.UUID  `json:"employee_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Date           time.Time  `json:"date"`
	Reason         string     `json:"reason"`
	ApprovedBy     *uuid.UUID `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
