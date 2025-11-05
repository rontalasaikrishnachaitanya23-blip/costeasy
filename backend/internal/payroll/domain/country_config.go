package domain

import (
	"time"

	"github.com/google/uuid"
)

// Country represents a master list of countries
type Country struct {
	ID                 uuid.UUID `json:"id"`
	Code               string    `json:"code"` // e.g. "AE"
	Name               string    `json:"name"`
	CurrencyCode       string    `json:"currency_code"`
	PhoneCode          string    `json:"phone_code"`
	TimeZone           string    `json:"time_zone"`
	WorkingDaysPerWeek int       `json:"working_days_per_week"`
	WeekendDays        []string  `json:"weekend_days"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CountryPayrollConfig stores payroll rules per country
type CountryPayrollConfig struct {
	ID                  uuid.UUID `json:"id"`
	CountryID           uuid.UUID `json:"country_id"`
	HasIncomeTax        bool      `json:"has_income_tax"`
	HasSocialSecurity   bool      `json:"has_social_security"`
	HasProfessionalTax  bool      `json:"has_professional_tax"`
	HasGratuity         bool      `json:"has_gratuity"`
	MinimumWage         float64   `json:"minimum_wage"`
	OvertimeMultiplier  float64   `json:"overtime_multiplier"`
	ProbationPeriodDays int       `json:"probation_period_days"`
	NoticePeriodDays    int       `json:"notice_period_days"`
	AnnualLeaveDays     int       `json:"annual_leave_days"`
	SickLeaveDays       int       `json:"sick_leave_days"`
	MaternityLeaveDays  int       `json:"maternity_leave_days"`
	PaternityLeaveDays  int       `json:"paternity_leave_days"`
	ConfigJSON          []byte    `json:"config_json,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
