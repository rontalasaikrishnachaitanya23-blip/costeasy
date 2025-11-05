package repository

import (
	"context"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type employeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return &employeeRepository{db: db}
}
func (r *employeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	query := `
		INSERT INTO employees (
			id, organization_id, employee_code,
			first_name, last_name, email, phone,
			country_id, department_id, designation_id,
			joined_at, relieved_at, termination_reason,
			employment_status, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			created_at, updated_at, created_by, updated_by
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21
		)
	`

	_, err := r.db.Exec(ctx, query,
		emp.ID,
		emp.OrganizationID,
		emp.EmployeeCode,
		emp.FirstName,
		emp.LastName,
		emp.Email,
		emp.Phone,
		emp.CountryID,
		emp.DepartmentID,
		emp.DesignationID,
		emp.JoinedAt,
		emp.RelievedAt,
		emp.TerminationReason,
		emp.EmploymentStatus,
		emp.IsSalaryStopped,
		emp.FinalSettlementGenerated,
		emp.FinalSettlementDate,
		emp.CreatedAt,
		emp.UpdatedAt,
		emp.CreatedBy,
		emp.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	return nil
}
func (r *employeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	query := `
		UPDATE employees SET
			first_name = $2,
			last_name = $3,
			email = $4,
			phone = $5,
			country_id = $6,
			department_id = $7,
			designation_id = $8,
			relieved_at = $9,
			termination_reason = $10,
			employment_status = $11,
			is_salary_stopped = $12,
			final_settlement_generated = $13,
			final_settlement_date = $14,
			updated_at = CURRENT_TIMESTAMP,
			updated_by = $15
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		emp.ID,
		emp.FirstName,
		emp.LastName,
		emp.Email,
		emp.Phone,
		emp.CountryID,
		emp.DepartmentID,
		emp.DesignationID,
		emp.RelievedAt,
		emp.TerminationReason,
		emp.EmploymentStatus,
		emp.IsSalaryStopped,
		emp.FinalSettlementGenerated,
		emp.FinalSettlementDate,
		emp.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}
func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, organization_id, employee_code,
			first_name, last_name, email, phone,
			country_id, department_id, designation_id,
			joined_at, relieved_at, termination_reason,
			employment_status, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			created_at, updated_at, created_by, updated_by
		FROM employees
		WHERE id = $1
	`

	var emp domain.Employee

	err := r.db.QueryRow(ctx, query, id).Scan(
		&emp.ID,
		&emp.OrganizationID,
		&emp.EmployeeCode,
		&emp.FirstName,
		&emp.LastName,
		&emp.Email,
		&emp.Phone,
		&emp.CountryID,
		&emp.DepartmentID,
		&emp.DesignationID,
		&emp.JoinedAt,
		&emp.RelievedAt,
		&emp.TerminationReason,
		&emp.EmploymentStatus,
		&emp.IsSalaryStopped,
		&emp.FinalSettlementGenerated,
		&emp.FinalSettlementDate,
		&emp.CreatedAt,
		&emp.UpdatedAt,
		&emp.CreatedBy,
		&emp.UpdatedBy,
	)

	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return &emp, nil
}
func (r *employeeRepository) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.Employee, error) {
	query := `
		SELECT id, organization_id, employee_code,
			first_name, last_name, email, phone,
			country_id, department_id, designation_id,
			joined_at, relieved_at, termination_reason,
			employment_status, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			created_at, updated_at, created_by, updated_by
		FROM employees
		WHERE organization_id = $1 AND employee_code = $2
	`

	var emp domain.Employee
	err := r.db.QueryRow(ctx, query, orgID, code).Scan(
		&emp.ID,
		&emp.OrganizationID,
		&emp.EmployeeCode,
		&emp.FirstName,
		&emp.LastName,
		&emp.Email,
		&emp.Phone,
		&emp.CountryID,
		&emp.DepartmentID,
		&emp.DesignationID,
		&emp.JoinedAt,
		&emp.RelievedAt,
		&emp.TerminationReason,
		&emp.EmploymentStatus,
		&emp.IsSalaryStopped,
		&emp.FinalSettlementGenerated,
		&emp.FinalSettlementDate,
		&emp.CreatedAt,
		&emp.UpdatedAt,
		&emp.CreatedBy,
		&emp.UpdatedBy,
	)

	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return &emp, nil
}
func (r *employeeRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Employee, error) {
	query := `
		SELECT id, organization_id, employee_code,
			first_name, last_name, email, phone,
			country_id, department_id, designation_id,
			joined_at, relieved_at, termination_reason,
			employment_status, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			created_at, updated_at, created_by, updated_by
		FROM employees
		WHERE organization_id = $1
		ORDER BY employee_code
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}
	defer rows.Close()

	var list []*domain.Employee

	for rows.Next() {
		var emp domain.Employee
		err := rows.Scan(
			&emp.ID,
			&emp.OrganizationID,
			&emp.EmployeeCode,
			&emp.FirstName,
			&emp.LastName,
			&emp.Email,
			&emp.Phone,
			&emp.CountryID,
			&emp.DepartmentID,
			&emp.DesignationID,
			&emp.JoinedAt,
			&emp.RelievedAt,
			&emp.TerminationReason,
			&emp.EmploymentStatus,
			&emp.IsSalaryStopped,
			&emp.FinalSettlementGenerated,
			&emp.FinalSettlementDate,
			&emp.CreatedAt,
			&emp.UpdatedAt,
			&emp.CreatedBy,
			&emp.UpdatedBy,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &emp)
	}

	return list, nil
}
func (r *employeeRepository) Terminate(ctx context.Context, id uuid.UUID, reason string, date string) error {
	query := `
		UPDATE employees SET
			employment_status = 'terminated',
			termination_reason = $2,
			relieved_at = $3,
			is_salary_stopped = true,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, reason, date)
	return err
}
func (r *employeeRepository) Relieve(ctx context.Context, id uuid.UUID, date string) error {
	query := `
		UPDATE employees SET
			employment_status = 'resigned',
			relieved_at = $2,
			updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id, date)
	return err
}
func (r *employeeRepository) StopSalary(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE employees SET
			is_salary_stopped = true,
			updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
func (r *employeeRepository) ResumeSalary(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE employees SET
			is_salary_stopped = false,
			updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
func (r *employeeRepository) MarkFinalSettlement(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE employees SET
			final_settlement_generated = true,
			final_settlement_date = NOW(),
			employment_status = 'finalized',
			updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}
