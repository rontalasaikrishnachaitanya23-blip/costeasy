package domain

import (
	"time"
	"github.com/google/uuid"
)

// PayrollPeriodLock prevents reprocessing after posting
type PayrollPeriodLock struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	PayrollPeriodID uuid.UUID  `json:"payroll_period_id"`
	IsLocked        bool       `json:"is_locked"`
	LockedBy        *uuid.UUID `json:"locked_by,omitempty"`
	LockedAt        *time.Time `json:"locked_at,omitempty"`
	Reason          *string    `json:"reason,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
