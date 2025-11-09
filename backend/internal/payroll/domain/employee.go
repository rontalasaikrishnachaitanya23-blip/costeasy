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
	EmploymentStatusFinalized  EmploymentStatus = "finalized"
)

type Employee struct {
	// Primary Info
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`

	// Identity & Contact
	EmployeeCode string  `json:"employee_code"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`

	// Personal Info
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Gender      *string    `json:"gender,omitempty"`
	Nationality *string    `json:"nationality,omitempty"`
	CountryID   *uuid.UUID `json:"country_id,omitempty"`

	// HR Structure
	DepartmentID  *uuid.UUID `json:"department_id,omitempty"`
	DesignationID *uuid.UUID `json:"designation_id,omitempty"`
	WorkLocation  *string    `json:"work_location,omitempty"`

	// Employment Dates
	DateOfJoining time.Time  `json:"date_of_joining"`
	JoinedAt      time.Time  `json:"joined_at"`
	DateOfExit    *time.Time `json:"date_of_exit,omitempty"`
	RelievedAt    *time.Time `json:"relieved_at,omitempty"`

	// Employment Status
	EmploymentStatus EmploymentStatus `json:"employment_status"`
	TerminationReason *string         `json:"termination_reason,omitempty"`

	// Payroll Info
	ContractType    string     `json:"contract_type"`
	SalaryCurrency  string     `json:"salary_currency"`
	BaseSalary      float64    `json:"base_salary"`
	LeavePolicyID   *uuid.UUID `json:"leave_policy_id,omitempty"`

	// Flags / Controls
	IsActive                bool  `json:"is_active"`
	IsSalaryStopped         bool  `json:"is_salary_stopped"`
	FinalSettlementGenerated bool `json:"final_settlement_generated"`
	FinalSettlementDate      *time.Time `json:"final_settlement_date,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Helper Methods -----------------------------------------------------

func (e *Employee) FullName() string {
	return e.FirstName + " " + e.LastName
}

// State Checks
func (e *Employee) CanBeTerminated() bool {
	return e.EmploymentStatus == EmploymentStatusActive
}

func (e *Employee) CanBeRelieved() bool {
	return e.EmploymentStatus == EmploymentStatusActive
}

// State Mutations
func (e *Employee) Terminate(reason string, date time.Time) {
	e.EmploymentStatus = EmploymentStatusTerminated
	e.TerminationReason = &reason
	e.RelievedAt = &date
	e.IsSalaryStopped = true
	e.UpdatedAt = time.Now()
}

func (e *Employee) Relieve(date time.Time) {
	e.EmploymentStatus = EmploymentStatusRelieved
	e.RelievedAt = &date
	e.UpdatedAt = time.Now()
}

func (e *Employee) MarkFinalSettlement() {
	now := time.Now()
	e.FinalSettlementGenerated = true
	e.FinalSettlementDate = &now
	e.EmploymentStatus = EmploymentStatusFinalized
	e.UpdatedAt = now
}
