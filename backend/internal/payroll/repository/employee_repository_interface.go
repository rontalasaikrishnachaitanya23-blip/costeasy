package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/chaitu35/costeasy/backend/internal/payroll/imports/types"
	"github.com/google/uuid"
)

type EmployeeRepository interface {
	Create(ctx context.Context, emp *domain.Employee) error
	Update(ctx context.Context, emp *domain.Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.Employee, error)
	List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Employee, error)

	// ✅ Workflow
	Terminate(ctx context.Context, id uuid.UUID, reason string, date string) error
	Relieve(ctx context.Context, id uuid.UUID, date string) error
	StopSalary(ctx context.Context, id uuid.UUID) error
	ResumeSalary(ctx context.Context, id uuid.UUID) error
	MarkFinalSettlement(ctx context.Context, id uuid.UUID) error

	// ✅ Import Support (NEW)
	CreateFromImport(ctx context.Context, orgID uuid.UUID, row types.RowValidated) error
}
