package domain

import (
	"time"

	"github.com/google/uuid"
)

// EmployeeSalaryDetail defines component-level pay structure
type EmployeeSalaryDetail struct {
	ID             uuid.UUID  `json:"id"`
	EmployeeID     uuid.UUID  `json:"employee_id"`
	OrganizationID uuid.UUID  `json:"organization_id"`
	ComponentCode  string     `json:"component_code"`
	ComponentName  string     `json:"component_name"`
	ComponentType  string     `json:"component_type"` // EARNING / DEDUCTION
	Amount         float64    `json:"amount"`
	IsRecurring    bool       `json:"is_recurring"`
	EffectiveFrom  time.Time  `json:"effective_from"`
	EffectiveTo    *time.Time `json:"effective_to,omitempty"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
