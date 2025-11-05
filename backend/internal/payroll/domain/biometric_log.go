package domain

import (
	"time"

	"github.com/google/uuid"
)

// BiometricLog represents a raw IN/OUT log from device
type BiometricLog struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	EmployeeID     *uuid.UUID `json:"employee_id,omitempty"`
	DeviceID       *string    `json:"device_id,omitempty"`
	LogType        string     `json:"log_type"` // IN / OUT
	LogTime        time.Time  `json:"log_time"`
	SourceFile     *string    `json:"source_file,omitempty"`
	SyncedAt       *time.Time `json:"synced_at,omitempty"`
	Processed      bool       `json:"processed"`
	CreatedAt      time.Time  `json:"created_at"`
}
