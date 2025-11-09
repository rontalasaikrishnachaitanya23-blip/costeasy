package service

import (
	"context"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/chaitu35/costeasy/backend/internal/payroll/imports/types"
	"github.com/google/uuid"
)

type EmployeeService interface {
	// CRUD
	CreateEmployee(ctx context.Context, input *domain.Employee) (*domain.Employee, error)
	UpdateEmployee(ctx context.Context, input *domain.Employee) error
	GetEmployeeByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	GetEmployeeByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.Employee, error)
	ListEmployees(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Employee, error)

	// Import (Excel)
	CreateFromImport(ctx context.Context, orgID uuid.UUID, row types.RowValidated) error

	// Workflow
	TerminateEmployee(ctx context.Context, id uuid.UUID, reason string, date time.Time) error
	RelieveEmployee(ctx context.Context, id uuid.UUID, date time.Time) error
	StopSalary(ctx context.Context, id uuid.UUID) error
	ResumeSalary(ctx context.Context, id uuid.UUID) error
	GenerateFinalSettlement(ctx context.Context, id uuid.UUID) error
}
