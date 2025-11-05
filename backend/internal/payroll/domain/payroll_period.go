package domain

import (
	"time"
	"github.com/google/uuid"
)

type PayrollPeriod struct {
	ID             uuid.UUID  `json:"id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	CountryID      uuid.UUID  `json:"country_id"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	Status         string     `json:"status"` // OPEN, CLOSED, POSTED
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
