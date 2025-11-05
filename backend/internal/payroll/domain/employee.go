package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmploymentStatus string

const (
	EmploymentStatusActive     EmploymentStatus = "active"
	EmploymentStatusOnLeave    EmploymentStatus = "on_leave"
	EmploymentStatusRelieved   EmploymentStatus = "relieved"
	EmploymentStatusTerminated EmploymentStatus = "terminated"
	EmploymentStatusFinalized  EmploymentStatus = "finalized" // after final settlement
)

// Employee represents an organization employee
type Employee struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	EmployeeCode   string    `json:"employee_code"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`

	CountryID     uuid.UUID `json:"country_id"`
	DepartmentID  uuid.UUID `json:"department_id"`
	DesignationID uuid.UUID `json:"designation_id"`

	JoinedAt          time.Time  `json:"joined_at"` // mandatory
	RelievedAt        *time.Time `json:"relieved_at"`
	TerminationReason *string    `json:"termination_reason"`

	EmploymentStatus EmploymentStatus `json:"employment_status"`

	// Salary control
	IsSalaryStopped bool `json:"is_salary_stopped"`

	// Final Settlement
	FinalSettlementGenerated bool       `json:"final_settlement_generated"`
	FinalSettlementDate      *time.Time `json:"final_settlement_date"`

	// Audit fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by"`
	UpdatedBy *uuid.UUID `json:"updated_by"`
}

// FullName returns the full name of employee
func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// IsActive returns true if employee is active and salary not stopped
func (e *Employee) IsActive() bool {
	return e.EmploymentStatus == EmploymentStatusActive && !e.IsSalaryStopped
}

// CanBeTerminated checks if termination is allowed
func (e *Employee) CanBeTerminated() bool {
	return e.EmploymentStatus == EmploymentStatusActive
}

// CanBeRelieved checks if relieving is possible
func (e *Employee) CanBeRelieved() bool {
	return e.EmploymentStatus == EmploymentStatusActive
}

// Terminate employee
func (e *Employee) Terminate(reason string, date time.Time) {
	e.EmploymentStatus = EmploymentStatusTerminated
	e.TerminationReason = &reason
	e.RelievedAt = &date
	e.IsSalaryStopped = true
	e.UpdatedAt = time.Now()
}

// Relieve employee
func (e *Employee) Relieve(date time.Time) {
	e.EmploymentStatus = EmploymentStatusRelieved
	e.RelievedAt = &date
	e.UpdatedAt = time.Now()
}

// MarkFinalSettlement marks final payroll closure
func (e *Employee) MarkFinalSettlement() {
	now := time.Now()
	e.FinalSettlementGenerated = true
	e.FinalSettlementDate = &now
	e.EmploymentStatus = EmploymentStatusFinalized
	e.UpdatedAt = now
}
