package mapper

import (
	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler/dto"
)

// ToAccountResponse converts domain.Account to AccountResponse
func ToAccountResponse(account domain.GLAccount) dto.AccountResponse {
	return dto.AccountResponse{
		ID:         account.ID,
		Code:       account.Code,
		Name:       account.Name,
		Type:       account.Type,
		ParentCode: account.ParentCode,
		IsActive:   account.IsActive,
		CreatedAt:  account.CreateAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  account.UpdateAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
