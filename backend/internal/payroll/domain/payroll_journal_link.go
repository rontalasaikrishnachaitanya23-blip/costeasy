package domain

import (
	"time"
	"github.com/google/uuid"
)

// PayrollJournalLink links payroll runs to GL journal entries
type PayrollJournalLink struct {
	ID             uuid.UUID  `json:"id"`
	PayrollRunID   uuid.UUID  `json:"payroll_run_id"`
	JournalEntryID uuid.UUID  `json:"journal_entry_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	LinkedAt       time.Time  `json:"linked_at"`
	LinkedBy       uuid.UUID  `json:"linked_by"`
}
