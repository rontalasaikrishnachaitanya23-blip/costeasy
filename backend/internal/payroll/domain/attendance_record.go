package domain

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceRecord stores daily attendance per employee
type AttendanceRecord struct {
	ID             uuid.UUID  `json:"id"`
	EmployeeID     uuid.UUID  `json:"employee_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	Date           time.Time  `json:"date"`
	Status         string     `json:"status"` // PRESENT, ABSENT, LEAVE, etc.
	CheckInTime    *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime   *time.Time `json:"check_out_time,omitempty"`
	HoursWorked    float64    `json:"hours_worked"`
	BatchID        *uuid.UUID `json:"batch_id,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
