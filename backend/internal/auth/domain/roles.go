// backend/internal/auth/domain/role.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a user role with permissions
type Role struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	DisplayName    string    `json:"display_name"`
	Description    *string   `json:"description,omitempty"`
	IsSystemRole   bool      `json:"is_system_role"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations (populated via joins)
	Permissions []Permission `json:"permissions,omitempty"`
}

// CanDelete checks if role can be deleted
func (r *Role) CanDelete() bool {
	return !r.IsSystemRole
}
