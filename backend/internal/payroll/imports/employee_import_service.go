package imports

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chaitu35/costeasy/backend/internal/payroll/domain"
	"github.com/chaitu35/costeasy/backend/internal/payroll/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xuri/excelize/v2"
)

type ImportResult struct {
	BatchID     uuid.UUID  `json:"batch_id"`
	TotalRows   int        `json:"total_rows"`
	ValidRows   int        `json:"valid_rows"`
	InvalidRows int        `json:"invalid_rows"`
	Errors      []RowError `json:"errors"`
	Status      string     `json:"status"` // success / partial / failed
}

type RowError struct {
	Row    int    `json:"row"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type EmployeeImportService struct {
	db      *pgxpool.Pool
	empRepo repository.EmployeeRepository
}

func NewEmployeeImportService(db *pgxpool.Pool, empRepo repository.EmployeeRepository) *EmployeeImportService {
	return &EmployeeImportService{db: db, empRepo: empRepo}
}

// Parse & Import Excel
func (s *EmployeeImportService) Import(ctx context.Context, orgID, uploadedBy uuid.UUID, fileName string, r io.Reader) (*ImportResult, error) {
	// create batch
	batchID := uuid.New()
	if err := s.createBatch(ctx, batchID, orgID, uploadedBy, fileName); err != nil {
		return nil, err
	}

	res := &ImportResult{BatchID: batchID}

	xf, err := excelize.OpenReader(r)
	if err != nil {
		_ = s.failBatch(ctx, batchID, "failed")
		return nil, fmt.Errorf("invalid excel file: %w", err)
	}
	defer xf.Close()

	rows, err := xf.GetRows(employeeTemplateSheet)
	if err != nil {
		_ = s.failBatch(ctx, batchID, "failed")
		return nil, fmt.Errorf("missing sheet '%s'", employeeTemplateSheet)
	}
	if len(rows) < 2 {
		_ = s.completeBatch(ctx, batchID, 0, 0, 0, "failed")
		return nil, errors.New("no data rows")
	}

	seenCodes := map[string]bool{}

	for i := 1; i < len(rows); i++ {
		res.TotalRows++
		rowno := i + 1
		row := normalizeRow(rows[i], rowno)

		// In-excel duplicate check
		if seenCodes[strings.ToLower(row.EmployeeCode)] {
			s.addError(ctx, batchID, rowno, "DUPLICATE_IN_FILE", "duplicate employee_code within file")
			res.InvalidRows++
			continue
		}
		seenCodes[strings.ToLower(row.EmployeeCode)] = true

		vr, vErrs := ValidateRow(row)
		if len(vErrs) > 0 {
			for _, e := range vErrs {
				s.addError(ctx, batchID, rowno, "VALIDATION", e)
			}
			res.InvalidRows++
			continue
		}

		// Lookups (country, dept, designation). For simplicity, by code/name.
		countryID, err := s.lookupCountry(ctx, vr.CountryCode)
		if err != nil {
			s.addError(ctx, batchID, rowno, "LOOKUP", "country_code not found: "+vr.CountryCode)
			res.InvalidRows++
			continue
		}
		deptID, err := s.lookupDepartment(ctx, orgID, vr.DepartmentName)
		if err != nil {
			s.addError(ctx, batchID, rowno, "LOOKUP", "department not found: "+vr.DepartmentName)
			res.InvalidRows++
			continue
		}
		desigID, err := s.lookupDesignation(ctx, orgID, vr.DesignationName)
		if err != nil {
			s.addError(ctx, batchID, rowno, "LOOKUP", "designation not found: "+vr.DesignationName)
			res.InvalidRows++
			continue
		}

		// Uniqueness in DB
		if existing, _ := s.empRepo.GetByCode(ctx, orgID, vr.EmployeeCode); existing != nil {
			s.addError(ctx, batchID, rowno, "DUPLICATE_IN_DB", "employee_code already exists")
			res.InvalidRows++
			continue
		}

		// Build domain employee
		now := time.Now()
		emp := &domain.Employee{
			ID:               uuid.New(),
			OrganizationID:   orgID,
			EmployeeCode:     vr.EmployeeCode,
			FirstName:        vr.FirstName,
			LastName:         vr.LastName,
			Email:            vr.Email,
			Phone:            vr.Phone,
			CountryID:        countryID,
			DepartmentID:     deptID,
			DesignationID:    desigID,
			JoinedAt:         vr.JoinedAt,
			EmploymentStatus: mapStatus(vr.EmploymentStatus),
			CreatedAt:        now,
			UpdatedAt:        now,
		}

		// status-specific fields
		if vr.TerminationAt != nil {
			emp.RelievedAt = vr.TerminationAt
			if strings.ToLower(vr.EmploymentStatus) == "terminated" {
				tr := strings.TrimSpace(vr.TerminationReason)
				emp.TerminationReason = &tr
				emp.IsSalaryStopped = true
			}
		}

		// Persist employee
		if err := s.empRepo.Create(ctx, emp); err != nil {
			s.addError(ctx, batchID, rowno, "DB_ERROR", err.Error())
			res.InvalidRows++
			continue
		}

		// Optionally insert base salary into employee_salary_details here if you want:
		//   INSERT INTO employee_salary_details(employee_id, base_salary, currency, effective_from, created_at)
		if err := s.insertBaseSalary(ctx, emp.ID, vr.BaseSalary, vr.SalaryCurrency, vr.JoinedAt); err != nil {
			// Don't fail the employee creation, but log the error
			s.addError(ctx, batchID, rowno, "DB_ERROR_SALARY", "salary save failed: "+err.Error())
			// still count as valid employee row
		}

		res.ValidRows++
	}

	// finalize batch
	status := "success"
	if res.InvalidRows > 0 && res.ValidRows > 0 {
		status = "partial"
	} else if res.InvalidRows > 0 && res.ValidRows == 0 {
		status = "failed"
	}
	if err := s.completeBatch(ctx, batchID, res.TotalRows, res.ValidRows, res.InvalidRows, status); err != nil {
		return nil, err
	}
	res.Status = status
	return res, nil
}

func normalizeRow(cols []string, rowno int) Row {
	get := func(i int) string {
		if i < len(cols) {
			return strings.TrimSpace(cols[i])
		}
		return ""
	}
	return Row{
		RowNo:              rowno,
		EmployeeCode:       get(0),
		FirstName:          get(1),
		LastName:           get(2),
		Email:              get(3),
		Phone:              get(4),
		CountryCode:        get(5),
		DepartmentName:     get(6),
		DesignationName:    get(7),
		JoinedAtStr:        get(8),
		BaseSalaryStr:      get(9),
		SalaryCurrency:     get(10),
		EmploymentStatus:   get(11),
		TerminationDateStr: get(12),
		TerminationReason:  get(13),
	}
}

// lookups (simple versions; adjust to your schema)
func (s *EmployeeImportService) lookupCountry(ctx context.Context, code string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `SELECT id FROM countries WHERE code = $1 AND is_active = true`, strings.ToUpper(code)).Scan(&id)
	return id, err
}

func (s *EmployeeImportService) lookupDepartment(ctx context.Context, orgID uuid.UUID, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT id FROM departments 
		WHERE organization_id = $1 AND LOWER(name) = LOWER($2) AND is_active = true
	`, orgID, name).Scan(&id)
	return id, err
}

func (s *EmployeeImportService) lookupDesignation(ctx context.Context, orgID uuid.UUID, name string) (uuid.UUID, error) {
	var id uuid.UUID
	err := s.db.QueryRow(ctx, `
		SELECT id FROM designations 
		WHERE organization_id = $1 AND LOWER(name) = LOWER($2) AND is_active = true
	`, orgID, name).Scan(&id)
	return id, err
}

// salary insert (optional – adjust table/columns as per your schema)
func (s *EmployeeImportService) insertBaseSalary(ctx context.Context, empID uuid.UUID, amount float64, currency string, effFrom time.Time) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO employee_salary_details (employee_id, base_salary, currency_code, effective_from, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`, empID, amount, strings.ToUpper(currency), effFrom)
	return err
}

// batch helpers
func (s *EmployeeImportService) createBatch(ctx context.Context, id, orgID, uploadedBy uuid.UUID, file string) error {
	_, err := s.db.Exec(ctx, `
		INSERT INTO employee_import_batches (id, organization_id, uploaded_by, file_name, status, created_at)
		VALUES ($1,$2,$3,$4,'processing',NOW())
	`, id, orgID, uploadedBy, file)
	return err
}

func (s *EmployeeImportService) completeBatch(ctx context.Context, id uuid.UUID, total, valid, invalid int, status string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE employee_import_batches
		SET total_rows=$2, valid_rows=$3, invalid_rows=$4, status=$5, completed_at=NOW()
		WHERE id=$1
	`, id, total, valid, invalid, status)
	return err
}

func (s *EmployeeImportService) failBatch(ctx context.Context, id uuid.UUID, status string) error {
	_, err := s.db.Exec(ctx, `
		UPDATE employee_import_batches SET status=$2, completed_at=NOW() WHERE id=$1
	`, id, status)
	return err
}

func (s *EmployeeImportService) addError(ctx context.Context, batchID uuid.UUID, row int, code, msg string) {
	_, _ = s.db.Exec(ctx, `
		INSERT INTO employee_import_errors (batch_id, row_number, error_code, message)
		VALUES ($1,$2,$3,$4)
	`, batchID, row, code, msg)
}

func mapStatus(s string) domain.EmploymentStatus {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "resigned":
		return domain.EmploymentStatusRelieved
	case "terminated":
		return domain.EmploymentStatusTerminated
	default:
		return domain.EmploymentStatusActive
	}
}
