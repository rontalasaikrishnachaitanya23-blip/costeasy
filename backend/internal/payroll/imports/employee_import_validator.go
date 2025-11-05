package imports

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var reEmail = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
var rePhone = regexp.MustCompile(`^[0-9+\-\s()]{6,20}$`)

type Row struct {
	EmployeeCode      string
	FirstName         string
	LastName          string
	Email             string
	Phone             string
	CountryCode       string
	DepartmentName    string
	DesignationName   string
	JoinedAtStr       string
	BaseSalaryStr     string
	SalaryCurrency    string
	EmploymentStatus  string
	TerminationDateStr string
	TerminationReason string
	RowNo             int
}

type RowValidated struct {
	Row
	JoinedAt       time.Time
	TerminationAt  *time.Time
	BaseSalary     float64
}

func ValidateRow(r Row) (RowValidated, []string) {
	var errs []string
	// required
	if strings.TrimSpace(r.EmployeeCode) == "" {
		errs = append(errs, "employee_code is required")
	}
	if strings.TrimSpace(r.FirstName) == "" {
		errs = append(errs, "first_name is required")
	}
	if strings.TrimSpace(r.Email) == "" || !reEmail.MatchString(r.Email) {
		errs = append(errs, "email is invalid")
	}
	if strings.TrimSpace(r.CountryCode) == "" {
		errs = append(errs, "country_code is required")
	}
	if strings.TrimSpace(r.DepartmentName) == "" {
		errs = append(errs, "department_name is required")
	}
	if strings.TrimSpace(r.DesignationName) == "" {
		errs = append(errs, "designation_name is required")
	}
	if strings.TrimSpace(r.JoinedAtStr) == "" {
		errs = append(errs, "joined_at is required")
	}
	if r.Phone != "" && !rePhone.MatchString(r.Phone) {
		errs = append(errs, "phone is invalid")
	}
	// joined_at
	var joined time.Time
	if r.JoinedAtStr != "" {
		t, err := time.Parse("2006-01-02", r.JoinedAtStr)
		if err != nil {
			errs = append(errs, "joined_at must be YYYY-MM-DD")
		}
		joined = t
	}
	// base_salary
	var sal float64
	if strings.TrimSpace(r.BaseSalaryStr) == "" {
		errs = append(errs, "base_salary is required")
	} else {
		_, err := fmt.Sscanf(r.BaseSalaryStr, "%f", &sal)
		if err != nil {
			errs = append(errs, "base_salary must be numeric")
		}
	}
	// employment_status
	stat := strings.ToLower(strings.TrimSpace(r.EmploymentStatus))
	if stat == "" {
		stat = "active"
	}
	if stat != "active" && stat != "resigned" && stat != "terminated" {
		errs = append(errs, "employment_status must be one of: active, resigned, terminated")
	}
	// termination
	var term *time.Time
	if stat == "terminated" || stat == "resigned" {
		if strings.TrimSpace(r.TerminationDateStr) == "" {
			errs = append(errs, "termination_date required when status is resigned/terminated")
		} else {
			t, err := time.Parse("2006-01-02", r.TerminationDateStr)
			if err != nil {
				errs = append(errs, "termination_date must be YYYY-MM-DD")
			} else if !joined.IsZero() && t.Before(joined) {
				errs = append(errs, "termination_date cannot be before joined_at")
			} else {
				term = &t
			}
		}
		if stat == "terminated" && strings.TrimSpace(r.TerminationReason) == "" {
			errs = append(errs, "termination_reason required when status is terminated")
		}
	}

	return RowValidated{
		Row:           r,
		JoinedAt:      joined,
		TerminationAt: term,
		BaseSalary:    sal,
	}, errs
}
