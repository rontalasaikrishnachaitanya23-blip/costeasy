package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/repository"
	"github.com/google/uuid"
)

const (
	ErrAccountCodeRequired = "account code is required"
	ErrAccountNameRequired = "account name is required"
	ErrAccountTypeRequired = "account type is required"
)

type AccountServiceInterface interface {
	CreateAccount(ctx context.Context, account domain.GLAccount) (domain.GLAccount, error)
	GetAccountByID(ctx context.Context, id uuid.UUID, includeIsActive bool) (domain.GLAccount, error)
	GetAccountByCode(ctx context.Context, code string, includeIsActive bool) (domain.GLAccount, error)
	ListAccounts(ctx context.Context, includeIsActive bool) ([]domain.GLAccount, error)
	DeactivateAccount(ctx context.Context, id uuid.UUID) error
	ActivateAccount(ctx context.Context, id uuid.UUID) error
	UpdateAccount(ctx context.Context, account domain.GLAccount) (domain.GLAccount, error)
	SoftDeleteAccount(ctx context.Context, id uuid.UUID) error
	SearchAccounts(ctx context.Context, params repository.GLAccountSearchParams) ([]domain.GLAccount, error)

	//helper methods can be added here

}

type AccountService struct {
	repo repository.GLAccountRepositoryInterface
}

func NewAccountService(repo repository.GLAccountRepositoryInterface) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, account domain.GLAccount) (domain.GLAccount, error) {
	if account.Code == "" {
		return domain.GLAccount{}, errors.New(ErrAccountCodeRequired)
	}
	if account.Name == "" {
		return domain.GLAccount{}, errors.New(ErrAccountNameRequired)
	}
	if account.Type == "" {
		return domain.GLAccount{}, errors.New(ErrAccountTypeRequired)
	}

	// Check for duplicates (including inactive records)
	existing, err := s.repo.GetGLAccountByCode(ctx, account.Code, false)
	if err == nil && existing.ID != uuid.Nil {
		return domain.GLAccount{}, fmt.Errorf("account with code %s already exists", account.Code)
	}

	return s.repo.CreateGLAccount(ctx, account)
}

func (s *AccountService) GetAccountByID(ctx context.Context, id uuid.UUID, includeIsActive bool) (domain.GLAccount, error) {
	return s.repo.GetGLAccountByID(ctx, id, includeIsActive)
}

func (s *AccountService) GetAccountByCode(ctx context.Context, code string, includeIsActive bool) (domain.GLAccount, error) {
	return s.repo.GetGLAccountByCode(ctx, code, includeIsActive)
}

func (s *AccountService) ListAccounts(ctx context.Context, includeIsActive bool) ([]domain.GLAccount, error) {
	return s.repo.ListGLAccounts(ctx, includeIsActive)
}

func (s *AccountService) DeactivateAccount(ctx context.Context, id uuid.UUID) error {
	// Can pass false here since we only deactivate active accounts
	account, err := s.repo.GetGLAccountByID(ctx, id, false)
	if err != nil {
		return fmt.Errorf("account with id %s does not exist", id)
	}

	if !account.IsActive {
		return fmt.Errorf("account is already inactive")
	}

	// Check for transactions and child accounts
	hasTransactions, err := s.repo.HasActiveTransactions(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking account transactions: %w", err)
	}
	if hasTransactions {
		return fmt.Errorf("cannot deactivate account with existing transactions")
	}

	hasChildren, err := s.repo.HasChildGLAccounts(ctx, id)
	if err != nil {
		return fmt.Errorf("error checking child accounts: %w", err)
	}
	if hasChildren {
		return fmt.Errorf("cannot deactivate account with active child accounts")
	}

	return s.repo.DeactivateGLAccount(ctx, id)
}

func (s *AccountService) ActivateAccount(ctx context.Context, id uuid.UUID) error {
	// Pass true to include inactive accounts
	account, err := s.repo.GetGLAccountByID(ctx, id, true) // ✅ true here
	if err != nil {
		return fmt.Errorf("account with id %s does not exist", id)
	}

	if account.IsActive {
		return fmt.Errorf("account is already active")
	}

	return s.repo.ActivateGLAccount(ctx, id)
}

func (s *AccountService) UpdateAccount(ctx context.Context, account domain.GLAccount) (domain.GLAccount, error) {
	// Validate required fields
	if account.Code == "" {
		return domain.GLAccount{}, errors.New(ErrAccountCodeRequired)
	}
	if account.Name == "" {
		return domain.GLAccount{}, errors.New(ErrAccountNameRequired)
	}
	if account.Type == "" {
		return domain.GLAccount{}, errors.New(ErrAccountTypeRequired)
	}

	// Check if account exists (allow updating inactive accounts)
	existing, err := s.repo.GetGLAccountByID(ctx, account.ID, false)
	if err != nil || existing.ID == uuid.Nil {
		return domain.GLAccount{}, fmt.Errorf("account with id %s does not exist", account.ID)
	}

	// Check for duplicate code if code is being changed
	if existing.Code != account.Code {
		duplicate, err := s.repo.GetGLAccountByCode(ctx, account.Code, false)
		if err == nil && duplicate.ID != uuid.Nil && duplicate.ID != account.ID {
			return domain.GLAccount{}, fmt.Errorf("account with code %s already exists", account.Code)
		}
	}

	return s.repo.UpdateGLAccount(ctx, account)
}

func (s *AccountService) SoftDeleteAccount(ctx context.Context, id uuid.UUID) error {
	// Check if account exists (including inactive)
	existing, err := s.repo.GetGLAccountByID(ctx, id, false)
	if err != nil || existing.ID == uuid.Nil {
		return fmt.Errorf("account with id %s does not exist", id)
	}
	return s.repo.SoftDeleteGLAccount(ctx, id)
}

func (s *AccountService) SearchAccounts(ctx context.Context, params repository.GLAccountSearchParams) ([]domain.GLAccount, error) {
	return s.repo.SearchGLAccounts(ctx, params)
}
