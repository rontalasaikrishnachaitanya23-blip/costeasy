// backend/internal/auth/domain/permission.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// Permission represents a granular permission
type Permission struct {
	ID          uuid.UUID `json:"id"`
	Module      string    `json:"module"`   // 'gl', 'settings', 'auth', 'dashboard'
	Resource    string    `json:"resource"` // 'accounts', 'journal_entries', 'users'
	Action      string    `json:"action"`   // 'view', 'create', 'edit', 'delete', 'export', 'print'
	DisplayName string    `json:"display_name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// Key returns unique permission key
func (p *Permission) Key() string {
	return p.Module + ":" + p.Resource + ":" + p.Action
}

// PermissionAction defines standard permission actions
type PermissionAction string

const (
	ActionView   PermissionAction = "view"
	ActionCreate PermissionAction = "create"
	ActionEdit   PermissionAction = "edit"
	ActionDelete PermissionAction = "delete"
	ActionExport PermissionAction = "export"
	ActionPrint  PermissionAction = "print"
)

// Module defines system modules
type Module string

const (
	ModuleGL          Module = "gl"
	ModuleSettings    Module = "settings"
	ModuleAuth        Module = "auth"
	ModuleDashboard   Module = "dashboard"
	ModuleEMR         Module = "emr"
	ModulePayroll     Module = "payroll"
	ModuleReports     Module = "reports"
	ModuleSubmissions Module = "submissions"
	ModuleCosting     Module = "costing"
	ModuleAudit       Module = "audit"
)
