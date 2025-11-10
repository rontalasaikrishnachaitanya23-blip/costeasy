package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PayrollFrequency string

const (
	PayrollFrequencyMonthly  PayrollFrequency = "monthly"
	PayrollFrequencyWeekly   PayrollFrequency = "weekly"
	PayrollFrequencyBiWeekly PayrollFrequency = "biweekly"
)

type PayrollPeriod struct {
	ID             uuid.UUID        `json:"id"`
	OrganizationID uuid.UUID        `json:"organization_id"`
	Name           string           `json:"name"`
	StartDate      time.Time        `json:"start_date"`
	EndDate        time.Time        `json:"end_date"`
	Month          int              `json:"month"`
	Year           int              `json:"year"`
	Frequency      PayrollFrequency `json:"frequency,omitempty"`
	IsLocked       bool             `json:"is_locked"`
	IsClosed       bool             `json:"is_closed"`
	ProcessedAt    *time.Time       `json:"processed_at,omitempty"`
	LockedBy       *uuid.UUID       `json:"locked_by,omitempty"`
	LockedAt       *time.Time       `json:"locked_at,omitempty"`
	ClosedBy       *uuid.UUID       `json:"closed_by,omitempty"`
	ClosedAt       *time.Time       `json:"closed_at,omitempty"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}

// Validate ensures data consistency
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

// Helper to auto-fill month/year
func (p *PayrollPeriod) SetMonthYear() {
	p.Month = int(p.StartDate.Month())
	p.Year = p.StartDate.Year()
}
