package domain

import (
	"time"

	"github.com/google/uuid"
)

// PayrollEntryLine stores computed pay for each component
type PayrollEntryLine struct {
	ID            uuid.UUID `json:"id"`
	PayrollRunID  uuid.UUID `json:"payroll_run_id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	ComponentCode string    `json:"component_code"`
	ComponentName string    `json:"component_name"`
	ComponentType string    `json:"component_type"` // EARNING / DEDUCTION
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}
