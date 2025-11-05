package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"` // Never expose password hash
	FirstName      *string   `json:"first_name,omitempty"`
	LastName       *string   `json:"last_name,omitempty"`
	Phone          *string   `json:"phone,omitempty"`

	// MFA fields
	// MFA Settings
	MFAEnabled    bool       `json:"mfa_enabled"`
	MFAMethod     MFAMethod  `json:"mfa_method"` // User's chosen method
	TOTPSecret    *string    `json:"-"`          // TOTP secret (hidden from JSON)
	BackupCodes   []string   `json:"-"`          // Encrypted backup codes
	MFAVerifiedAt *time.Time `json:"mfa_verified_at"`

	// OrganizationID  *uuid.UUID `json:"organization_id"`

	// Status fields
	IsActive        bool       `json:"is_active"`
	IsVerified      bool       `json:"is_verified"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP     *string    `json:"last_login_ip,omitempty"`

	// Password management
	PasswordChangedAt   time.Time  `json:"password_changed_at"`
	FailedLoginAttempts int        `json:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"locked_until,omitempty"`

	//Remote access Control

	AllowRemoteAccess      bool       `json:"allow_remote_access"`
	RemoteAccessReason     *string    `json:"remote_access_reason,omitempty"`
	RemoteAccessApprovedBy *uuid.UUID `json:"remote_access_approved_by,omitempty"`
	RemoteAccessApprovedAt *time.Time `json:"remote_access_approved_at,omitempty"`

	//AllowRemoteAccess bool   `json: `

	// Audit fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`

	// Relations (not stored in DB, populated via joins)
	Roles []Role `json:"roles,omitempty"`
}

// MFAMethod type (reuse from organization)
type MFAMethod string

const (
	MFAMethodNone  MFAMethod = "none"
	MFAMethodSMS   MFAMethod = "sms"
	MFAMethodEmail MFAMethod = "email"
	MFAMethodTOTP  MFAMethod = "totp"
)

// IsLocked checks if user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// CanLogin checks if user can attempt login
func (u *User) CanLogin() bool {
	return u.IsActive && !u.IsLocked()
}

// FullName returns user's full name
func (u *User) FullName() string {
	firstName := ""
	lastName := ""

	if u.FirstName != nil {
		firstName = *u.FirstName
	}
	if u.LastName != nil {
		lastName = *u.LastName
	}

	if firstName == "" && lastName == "" {
		return u.Username
	}

	return firstName + " " + lastName
}

// IncrementFailedAttempts increments failed login counter
func (u *User) IncrementFailedAttempts() {
	u.FailedLoginAttempts++

	// Lock account after 5 failed attempts for 30 minutes
	if u.FailedLoginAttempts >= 5 {
		lockUntil := time.Now().Add(30 * time.Minute)
		u.LockedUntil = &lockUntil
	}
}

// ResetFailedAttempts resets failed login counter
func (u *User) ResetFailedAttempts() {
	u.FailedLoginAttempts = 0
	u.LockedUntil = nil
}

// UpdateLastLogin updates last login timestamp and IP
func (u *User) UpdateLastLogin(ipAddress string) {
	now := time.Now()
	u.LastLoginAt = &now
	u.LastLoginIP = &ipAddress
}

// HasRemoteAccess checks if user can access from outside office
func (u *User) HasRemoteAccess() bool {
	return u.AllowRemoteAccess
}
