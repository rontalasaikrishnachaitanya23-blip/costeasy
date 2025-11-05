package domain

import (
	"time"

	"github.com/google/uuid"
)

// PayrollRun represents a batch payroll execution
type PayrollRun struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	PayrollPeriodID uuid.UUID  `json:"payroll_period_id"`
	ReferenceCode   string     `json:"reference_code"`
	Status          string     `json:"status"` // DRAFT, POSTED, REVERSED
	TotalEmployees  int        `json:"total_employees"`
	TotalDebit      float64    `json:"total_debit"`
	TotalCredit     float64    `json:"total_credit"`
	GLJournalID     *uuid.UUID `json:"gl_journal_id,omitempty"`
	PostedBy        *uuid.UUID `json:"posted_by,omitempty"`
	PostedAt        *time.Time `json:"posted_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
