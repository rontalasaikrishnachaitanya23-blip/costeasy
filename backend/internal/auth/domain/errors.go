// backend/internal/auth/domain/errors.go
package domain

import "errors"

var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserInactive      = errors.New("user account is inactive")
	ErrUserLocked        = errors.New("user account is locked")
	ErrUserNotVerified   = errors.New("user email not verified")

	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenExpired       = errors.New("token has expired")
	ErrTokenRevoked       = errors.New("token has been revoked")

	// MFA errors
	ErrMFARequired       = errors.New("MFA verification required")
	ErrInvalidMFACode    = errors.New("invalid MFA code")
	ErrMFANotEnabled     = errors.New("MFA is not enabled for this user")
	ErrMFAAlreadyEnabled = errors.New("MFA is already enabled for this user")

	// Password errors
	ErrWeakPassword     = errors.New("password does not meet security requirements")
	ErrPasswordMismatch = errors.New("passwords do not match")
	ErrSamePassword     = errors.New("new password must be different from current password")

	// Role errors
	ErrRoleNotFound           = errors.New("role not found")
	ErrRoleAlreadyExists      = errors.New("role already exists")
	ErrSystemRoleModification = errors.New("cannot modify system role")

	// Permission errors
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionDenied        = errors.New("permission denied")
	ErrInsufficientPermissions = errors.New("insufficient permissions")

	// Validation errors
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidUsername = errors.New("invalid username format")
	ErrEmptyField      = errors.New("required field is empty")
)
