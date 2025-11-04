// backend/internal/auth/repository/refresh_token_repository_interface.go
package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/google/uuid"
)

// RefreshTokenRepositoryInterface defines the contract for refresh token operations
type RefreshTokenRepositoryInterface interface {
	// Create creates a new refresh token
	Create(ctx context.Context, token *domain.RefreshToken) error

	// GetByID retrieves a refresh token by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.RefreshToken, error)

	// GetByTokenHash retrieves a refresh token by its hash
	GetByTokenHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)

	// GetByUserID retrieves all refresh tokens for a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.RefreshToken, error)

	// Revoke revokes a refresh token
	Revoke(ctx context.Context, id uuid.UUID, replacedBy *uuid.UUID) error

	// RevokeAll revokes all refresh tokens for a user
	RevokeAll(ctx context.Context, userID uuid.UUID) error

	// Delete deletes a refresh token
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteExpired deletes all expired refresh tokens
	DeleteExpired(ctx context.Context) error

	// DeleteByUserID deletes all refresh tokens for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}
