package domain

import (
	"time"

	"github.com/google/uuid"
)

type PayrollGLMapping struct {
	ID              uuid.UUID  `json:"id"`
	OrganizationID  uuid.UUID  `json:"organization_id"`
	CountryCode     string     `json:"country_code"`
	ComponentType   string     `json:"component_type"`
	ComponentName   string     `json:"component_name"`
	DebitAccountID  *uuid.UUID `json:"debit_account_id"`
	CreditAccountID *uuid.UUID `json:"credit_account_id"`
	Description     *string    `json:"description,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
