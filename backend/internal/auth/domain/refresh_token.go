// backend/internal/auth/domain/refresh_token.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a JWT refresh token
type RefreshToken struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	TokenHash  string     `json:"-"` // SHA256 hash of token
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	ReplacedBy *uuid.UUID `json:"replaced_by,omitempty"`
	DeviceInfo *string    `json:"device_info,omitempty"`
	IPAddress  *string    `json:"ip_address,omitempty"`
}

// IsExpired checks if token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked checks if token is revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsValid checks if token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}

// Revoke marks token as revoked
func (rt *RefreshToken) Revoke() {
	now := time.Now()
	rt.RevokedAt = &now
}
