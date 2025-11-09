package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/chaitu35/costeasy/backend/internal/payroll/imports/types"
	"github.com/chaitu35/costeasy/backend/internal/payroll/repository"
	"github.com/google/uuid"
)

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// ========================================================
// CRUD
// ========================================================

func (s *employeeService) CreateEmployee(ctx context.Context, emp *domain.Employee) (*domain.Employee, error) {
	// Apply system defaults
	now := time.Now()
	emp.ID = uuid.New()
	emp.CreatedAt = now
	emp.UpdatedAt = now

	if emp.EmploymentStatus == "" {
		emp.EmploymentStatus = domain.EmploymentStatusActive
	}

	if emp.ContractType == "" {
		emp.ContractType = "permanent"
	}

	if emp.SalaryCurrency == "" {
		emp.SalaryCurrency = "AED"
	}

	if emp.JoinedAt.IsZero() {
		emp.JoinedAt = emp.DateOfJoining
	}

	err := s.repo.Create(ctx, emp)
	if err != nil {
		return nil, fmt.Errorf("create employee failed: %w", err)
	}

	return emp, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, emp *domain.Employee) error {
	emp.UpdatedAt = time.Now()
	return s.repo.Update(ctx, emp)
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *employeeService) GetEmployeeByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.Employee, error) {
	return s.repo.GetByCode(ctx, orgID, code)
}

func (s *employeeService) ListEmployees(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Employee, error) {
	return s.repo.List(ctx, orgID, limit, offset)
}

// ========================================================
// IMPORT (Excel validated rows)
// ========================================================

func (s *employeeService) CreateFromImport(
	ctx context.Context,
	orgID uuid.UUID,
	row types.RowValidated,
) error {
	return s.repo.CreateFromImport(ctx, orgID, row)
}

// ========================================================
// WORKFLOW / STATUS TRANSITIONS
// ========================================================

func (s *employeeService) TerminateEmployee(ctx context.Context, id uuid.UUID, reason string, date time.Time) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !emp.CanBeTerminated() {
		return fmt.Errorf("employee cannot be terminated in status: %s", emp.EmploymentStatus)
	}

	emp.Terminate(reason, date)
	return s.repo.Terminate(ctx, id, reason, date.Format("2006-01-02"))
}

func (s *employeeService) RelieveEmployee(ctx context.Context, id uuid.UUID, date time.Time) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !emp.CanBeRelieved() {
		return fmt.Errorf("employee cannot be relieved in status: %s", emp.EmploymentStatus)
	}

	emp.Relieve(date)
	return s.repo.Relieve(ctx, id, date.Format("2006-01-02"))
}

func (s *employeeService) StopSalary(ctx context.Context, id uuid.UUID) error {
	return s.repo.StopSalary(ctx, id)
}

func (s *employeeService) ResumeSalary(ctx context.Context, id uuid.UUID) error {
	return s.repo.ResumeSalary(ctx, id)
}

func (s *employeeService) GenerateFinalSettlement(ctx context.Context, id uuid.UUID) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if emp.EmploymentStatus != domain.EmploymentStatusRelieved &&
		emp.EmploymentStatus != domain.EmploymentStatusTerminated {
		return fmt.Errorf("final settlement allowed only for relieved or terminated employees")
	}

	emp.MarkFinalSettlement()
	return s.repo.MarkFinalSettlement(ctx, id)
}
