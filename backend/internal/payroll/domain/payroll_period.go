package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PayrollPeriod struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	OrganizationID uuid.UUID  `db:"organization_id" json:"organization_id"`
	StartDate      time.Time  `db:"start_date" json:"start_date"`
	EndDate        time.Time  `db:"end_date" json:"end_date"`
	Month          int        `db:"month" json:"month"`
	Year           int        `db:"year" json:"year"`
	IsLocked       bool       `db:"is_locked" json:"is_locked"`
	ProcessedAt    *time.Time `db:"processed_at" json:"processed_at,omitempty"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

// Validate checks logical consistency of the payroll period.
func (p *PayrollPeriod) Validate() error {
	if p.OrganizationID == uuid.Nil {
		return errors.New("organization_id is required")
	}
	if p.StartDate.IsZero() || p.EndDate.IsZero() {
		return errors.New("start_date and end_date are required")
	}
	if p.EndDate.Before(p.StartDate) {
		return errors.New("end_date cannot be before start_date")
	}
	if p.Month < 1 || p.Month > 12 {
		return errors.New("month must be between 1 and 12")
	}
	if p.Year < 2000 {
		return errors.New("invalid year")
	}
	return nil
}
