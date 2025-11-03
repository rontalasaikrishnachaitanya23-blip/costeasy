package domain

import (
	"time"

	"github.com/google/uuid"
)

type GLAccount struct {
	ID         uuid.UUID   `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Type       AccountType `json:"type"`
	ParentCode *uuid.UUID  `json:"parent_id,omitempty"`
	CreateAt   time.Time   `json:"create_at"`
	UpdateAt   time.Time   `json:"update_at"`
	IsActive   bool        `json:"is_active"`
}
