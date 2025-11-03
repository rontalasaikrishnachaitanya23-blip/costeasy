package repository

import (
	"context"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/google/uuid"
)

type GLAccountRepositoryInterface interface {
	CreateGLAccount(ctx context.Context, GLAccount domain.GLAccount) (domain.GLAccount, error)
	GetGLAccountByID(ctx context.Context, id uuid.UUID, includeInactive bool) (domain.GLAccount, error)
	GetGLAccountByCode(ctx context.Context, code string, includeInactive bool) (domain.GLAccount, error)
	ListGLAccounts(ctx context.Context, includeInactive bool) ([]domain.GLAccount, error)
	UpdateGLAccount(ctx context.Context, GLAccount domain.GLAccount) (domain.GLAccount, error)
	DeactivateGLAccount(ctx context.Context, id uuid.UUID) error
	ActivateGLAccount(ctx context.Context, id uuid.UUID) error
	SoftDeleteGLAccount(ctx context.Context, id uuid.UUID) error
	SearchGLAccounts(ctx context.Context, params GLAccountSearchParams) ([]domain.GLAccount, error)
	HasActiveTransactions(ctx context.Context, GLAccountID uuid.UUID) (bool, error)
	HasChildGLAccounts(ctx context.Context, GLAccountID uuid.UUID) (bool, error)
}
