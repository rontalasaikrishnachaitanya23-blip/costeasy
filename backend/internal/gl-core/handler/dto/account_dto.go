// backend/internal/gl-core/handler/dto/account_dto.go
package dto

import (
	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
)

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	Code     string             `json:"code" binding:"required"`
	Name     string             `json:"name" binding:"required"`
	Type     domain.AccountType `json:"type" binding:"required"`
	ParentCode *uuid.UUID         `json:"parent_id"`
}

// UpdateAccountRequest represents the request body for updating an account
type UpdateAccountRequest struct {
	Name     string     `json:"name" binding:"required"`
	ParentCode *uuid.UUID `json:"parent_id"`
}

// AccountResponse represents the response structure for account operations
type AccountResponse struct {
	ID        uuid.UUID          `json:"id"`
	Code      string             `json:"code"`
	Name      string             `json:"name"`
	Type      domain.AccountType `json:"type"`
	ParentCode  *uuid.UUID         `json:"parent_id,omitempty"`
	IsActive  bool               `json:"is_active"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}


