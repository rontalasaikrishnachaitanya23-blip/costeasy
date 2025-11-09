package repository

import (
	"context"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/chaitu35/costeasy/backend/internal/payroll/imports/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type employeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) EmployeeRepository {
	return &employeeRepository{db: db}
}

//
// ------------------------------------------------------------
// CREATE
// ------------------------------------------------------------
//

func (r *employeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	query := `
		INSERT INTO employees (
			id, organization_id, country_id,
			employee_code, first_name, last_name,
			email, phone, date_of_birth, gender, nationality,
			date_of_joining, date_of_exit,
			work_location, contract_type,
			salary_currency, base_salary,
			is_active, created_at, updated_at,
			department_id, designation_id,
			employment_status, joined_at, relieved_at,
			termination_reason, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			leave_policy_id
		)
		VALUES (
			$1,$2,$3,
			$4,$5,$6,
			$7,$8,$9,$10,$11,
			$12,$13,
			$14,$15,
			$16,$17,
			$18,$19,$20,
			$21,$22,
			$23,$24,$25,
			$26,$27,
			$28,$29,
			$30
		)
	`

	_, err := r.db.Exec(ctx, query,
		emp.ID, emp.OrganizationID, emp.CountryID,
		emp.EmployeeCode, emp.FirstName, emp.LastName,
		emp.Email, emp.Phone, emp.DateOfBirth, emp.Gender, emp.Nationality,
		emp.DateOfJoining, emp.DateOfExit,
		emp.WorkLocation, emp.ContractType,
		emp.SalaryCurrency, emp.BaseSalary,
		emp.IsActive, emp.CreatedAt, emp.UpdatedAt,
		emp.DepartmentID, emp.DesignationID,
		emp.EmploymentStatus, emp.JoinedAt, emp.RelievedAt,
		emp.TerminationReason, emp.IsSalaryStopped,
		emp.FinalSettlementGenerated, emp.FinalSettlementDate,
		emp.LeavePolicyID,
	)

	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	return nil
}

//
// ------------------------------------------------------------
// UPDATE
// ------------------------------------------------------------
//

func (r *employeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	query := `
		UPDATE employees SET
			first_name = $2,
			last_name = $3,
			email = $4,
			phone = $5,
			date_of_birth = $6,
			gender = $7,
			nationality = $8,
			date_of_joining = $9,
			date_of_exit = $10,
			work_location = $11,
			contract_type = $12,
			salary_currency = $13,
			base_salary = $14,
			is_active = $15,
			updated_at = $16,
			department_id = $17,
			designation_id = $18,
			employment_status = $19,
			joined_at = $20,
			relieved_at = $21,
			termination_reason = $22,
			is_salary_stopped = $23,
			final_settlement_generated = $24,
			final_settlement_date = $25,
			leave_policy_id = $26
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query,
		emp.ID,
		emp.FirstName, emp.LastName,
		emp.Email, emp.Phone,
		emp.DateOfBirth, emp.Gender, emp.Nationality,
		emp.DateOfJoining, emp.DateOfExit,
		emp.WorkLocation, emp.ContractType,
		emp.SalaryCurrency, emp.BaseSalary,
		emp.IsActive, emp.UpdatedAt,
		emp.DepartmentID, emp.DesignationID,
		emp.EmploymentStatus, emp.JoinedAt, emp.RelievedAt,
		emp.TerminationReason, emp.IsSalaryStopped,
		emp.FinalSettlementGenerated, emp.FinalSettlementDate,
		emp.LeavePolicyID,
	)

	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

//
// ------------------------------------------------------------
// GET BY ID
// ------------------------------------------------------------
//

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT 
			id, organization_id, country_id,
			employee_code, first_name, last_name,
			email, phone, date_of_birth, gender, nationality,
			date_of_joining, date_of_exit,
			work_location, contract_type,
			salary_currency, base_salary,
			is_active, created_at, updated_at,
			department_id, designation_id,
			employment_status, joined_at, relieved_at,
			termination_reason, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			leave_policy_id
		FROM employees WHERE id = $1
	`

	var emp domain.Employee
	err := r.db.QueryRow(ctx, query, id).Scan(
		&emp.ID, &emp.OrganizationID, &emp.CountryID,
		&emp.EmployeeCode, &emp.FirstName, &emp.LastName,
		&emp.Email, &emp.Phone, &emp.DateOfBirth, &emp.Gender, &emp.Nationality,
		&emp.DateOfJoining, &emp.DateOfExit,
		&emp.WorkLocation, &emp.ContractType,
		&emp.SalaryCurrency, &emp.BaseSalary,
		&emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
		&emp.DepartmentID, &emp.DesignationID,
		&emp.EmploymentStatus, &emp.JoinedAt, &emp.RelievedAt,
		&emp.TerminationReason, &emp.IsSalaryStopped,
		&emp.FinalSettlementGenerated, &emp.FinalSettlementDate,
		&emp.LeavePolicyID,
	)

	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return &emp, nil
}

//
// ------------------------------------------------------------
// GET BY CODE
// ------------------------------------------------------------
//

func (r *employeeRepository) GetByCode(ctx context.Context, orgID uuid.UUID, code string) (*domain.Employee, error) {
	query := `
		SELECT 
			id, organization_id, country_id,
			employee_code, first_name, last_name,
			email, phone, date_of_birth, gender, nationality,
			date_of_joining, date_of_exit,
			work_location, contract_type,
			salary_currency, base_salary,
			is_active, created_at, updated_at,
			department_id, designation_id,
			employment_status, joined_at, relieved_at,
			termination_reason, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			leave_policy_id
		FROM employees 
		WHERE organization_id = $1 AND employee_code = $2
	`

	var emp domain.Employee
	err := r.db.QueryRow(ctx, query, orgID, code).Scan(
		&emp.ID, &emp.OrganizationID, &emp.CountryID,
		&emp.EmployeeCode, &emp.FirstName, &emp.LastName,
		&emp.Email, &emp.Phone, &emp.DateOfBirth, &emp.Gender, &emp.Nationality,
		&emp.DateOfJoining, &emp.DateOfExit,
		&emp.WorkLocation, &emp.ContractType,
		&emp.SalaryCurrency, &emp.BaseSalary,
		&emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
		&emp.DepartmentID, &emp.DesignationID,
		&emp.EmploymentStatus, &emp.JoinedAt, &emp.RelievedAt,
		&emp.TerminationReason, &emp.IsSalaryStopped,
		&emp.FinalSettlementGenerated, &emp.FinalSettlementDate,
		&emp.LeavePolicyID,
	)

	if err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	return &emp, nil
}

//
// ------------------------------------------------------------
// LIST
// ------------------------------------------------------------
//

func (r *employeeRepository) List(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*domain.Employee, error) {
	query := `
		SELECT
			id, organization_id, country_id,
			employee_code, first_name, last_name,
			email, phone, date_of_birth, gender, nationality,
			date_of_joining, date_of_exit,
			work_location, contract_type,
			salary_currency, base_salary,
			is_active, created_at, updated_at,
			department_id, designation_id,
			employment_status, joined_at, relieved_at,
			termination_reason, is_salary_stopped,
			final_settlement_generated, final_settlement_date,
			leave_policy_id
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
			&emp.ID, &emp.OrganizationID, &emp.CountryID,
			&emp.EmployeeCode, &emp.FirstName, &emp.LastName,
			&emp.Email, &emp.Phone, &emp.DateOfBirth, &emp.Gender, &emp.Nationality,
			&emp.DateOfJoining, &emp.DateOfExit,
			&emp.WorkLocation, &emp.ContractType,
			&emp.SalaryCurrency, &emp.BaseSalary,
			&emp.IsActive, &emp.CreatedAt, &emp.UpdatedAt,
			&emp.DepartmentID, &emp.DesignationID,
			&emp.EmploymentStatus, &emp.JoinedAt, &emp.RelievedAt,
			&emp.TerminationReason, &emp.IsSalaryStopped,
			&emp.FinalSettlementGenerated, &emp.FinalSettlementDate,
			&emp.LeavePolicyID,
		)

		if err != nil {
			return nil, err
		}
		list = append(list, &emp)
	}

	return list, nil
}

//
// ------------------------------------------------------------
// WORKFLOW ACTIONS
// ------------------------------------------------------------
//

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
			employment_status = 'relieved',
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

//
// ------------------------------------------------------------
// IMPORT SUPPORT
// ------------------------------------------------------------
//

func (r *employeeRepository) CreateFromImport(
	ctx context.Context,
	orgID uuid.UUID,
	row types.RowValidated,
) error {

	query := `
		INSERT INTO employees (
			organization_id,
			employee_code,
			first_name,
			last_name,
			email,
			phone,
			country_id,
			department_id,
			designation_id,
			leave_policy_id,
			date_of_joining,
			joined_at,
			date_of_exit,
			termination_reason,
			employment_status,
			contract_type,
			salary_currency,
			base_salary,
			is_active,
			is_salary_stopped,
			final_settlement_generated
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
			$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21
		)
	`

	_, err := r.db.Exec(ctx, query,
		orgID,
		row.EmployeeCode,
		row.FirstName,
		row.LastName,
		row.Email,
		row.Phone,
		row.CountryID,
		row.DepartmentID,
		row.DesignationID,
		row.LeavePolicyID,
		row.JoinedAt,
		row.JoinedAt,
		row.TerminationDate,
		row.TerminationReason,
		row.EmploymentStatus,
		"permanent",
		row.SalaryCurrency,
		row.BaseSalary,
		true,
		false,
		false,
	)

	if err != nil {
		return fmt.Errorf("failed to import employee (%s): %w", row.EmployeeCode, err)
	}

	return nil
}
