package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/chaitu35/costeasy/backend/internal/payroll/repository"

	"github.com/google/uuid"
)

type employeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) CreateEmployee(ctx context.Context, emp *domain.Employee) (*domain.Employee, error) {

	// ----- VALIDATIONS -----

	if emp.OrganizationID == uuid.Nil {
		return nil, errors.New("organization_id is required")
	}

	if emp.EmployeeCode == "" {
		return nil, errors.New("employee_code is required")
	}

	if emp.JoinedAt.IsZero() {
		return nil, errors.New("joined_at date is required")
	}

	// Check duplicate employee code
	existing, _ := s.repo.GetByCode(ctx, emp.OrganizationID, emp.EmployeeCode)
	if existing != nil {
		return nil, fmt.Errorf("employee_code '%s' already exists", emp.EmployeeCode)
	}

	// Default system metadata
	now := time.Now()
	emp.ID = uuid.New()
	emp.CreatedAt = now
	emp.UpdatedAt = now
	emp.EmploymentStatus = domain.EmploymentStatusActive
	emp.RelievedAt = nil
	emp.FinalSettlementGenerated = false
	emp.IsSalaryStopped = false

	// ----- SAVE -----
	if err := s.repo.Create(ctx, emp); err != nil {
		return nil, err
	}

	return emp, nil
}
func (s *employeeService) UpdateEmployee(ctx context.Context, emp *domain.Employee) error {
	existing, err := s.repo.GetByID(ctx, emp.ID)
	if err != nil {
		return fmt.Errorf("cannot update employee: %w", err)
	}

	if existing.EmploymentStatus == domain.EmploymentStatusTerminated ||
		existing.EmploymentStatus == domain.EmploymentStatusFinalized {
		return errors.New("cannot update employee after termination/final settlement")
	}

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
func (s *employeeService) TerminateEmployee(ctx context.Context, id uuid.UUID, reason string, date time.Time) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if emp.EmploymentStatus == domain.EmploymentStatusTerminated {
		return errors.New("employee already terminated")
	}

	// ✅ Mandatory reason check
	if len(strings.TrimSpace(reason)) == 0 {
		return errors.New("termination reason is required")
	}

	// ✅ Date must be after joining date
	if date.Before(emp.JoinedAt) {
		return errors.New("termination date cannot be before joining date")
	}

	return s.repo.Terminate(ctx, id, reason, date.Format("2006-01-02"))
}

func (s *employeeService) RelieveEmployee(ctx context.Context, id uuid.UUID, date time.Time) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if date.Before(emp.JoinedAt) {
		return errors.New("relieving date cannot be before joining date")
	}

	if emp.EmploymentStatus == domain.EmploymentStatusRelieved {
		return errors.New("cannot relieve a terminated employee")
	}

	return s.repo.Relieve(ctx, id, date.Format("2006-01-02"))
}
func (s *employeeService) StopSalary(ctx context.Context, id uuid.UUID) error {
	return s.repo.StopSalary(ctx, id)
}
func (s *employeeService) ResumeSalary(ctx context.Context, id uuid.UUID) error {
	emp, _ := s.repo.GetByID(ctx, id)

	if emp.EmploymentStatus != domain.EmploymentStatusActive {
		return errors.New("salary can only be resumed for active employees")
	}

	return s.repo.ResumeSalary(ctx, id)
}
func (s *employeeService) GenerateFinalSettlement(ctx context.Context, id uuid.UUID) error {
	emp, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if emp.EmploymentStatus != domain.EmploymentStatusRelieved &&
		emp.EmploymentStatus != domain.EmploymentStatusTerminated {
		return errors.New("final settlement allowed only after employee is relieved or terminated")
	}
	if emp.FinalSettlementGenerated {
		return errors.New("final settlement already processed")
	}

	return s.repo.MarkFinalSettlement(ctx, id)
}
