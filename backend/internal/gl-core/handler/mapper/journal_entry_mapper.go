// backend/internal/gl-core/handler/dto/mapper.go
package mapper

import (
	"github.com/chaitu35/costeasy/backend/internal/gl-core/domain"
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler/dto"
)

// ToJournalEntryResponse converts domain.JournalEntry to JournalEntryResponse
func ToJournalEntryResponse(entry *domain.JournalEntry) dto.JournalEntryResponse {
	var postingDate *string
	if entry.PostingDate != nil {
		pd := entry.PostingDate.Format("2006-01-02T15:04:05Z07:00")
		postingDate = &pd
	}

	lines := make([]dto.JournalLineResponse, len(entry.Lines))
	for i, line := range entry.Lines {
		lines[i] = dto.JournalLineResponse{
			ID:          line.ID.String(),
			AccountID:   line.AccountID.String(),
			LineNumber:  line.LineNumber,
			Reference:   line.Reference,
			Description: line.Description,
			Debit:       line.Debit,
			Credit:      line.Credit,
		}
	}

	return dto.JournalEntryResponse{
		ID:              entry.ID.String(),
		OrganizationID:  entry.OrganizationID.String(),
		EntryNumber:     entry.EntryNumber,
		TransactionDate: entry.TransactionDate.Format("2006-01-02"),
		PostingDate:     postingDate,
		Reference:       entry.Reference,
		Description:     entry.Description,
		Status:          string(entry.Status),
		TotalDebit:      entry.TotalDebit,
		TotalCredit:     entry.TotalCredit,
		Lines:           lines,
		CreatedAt:       entry.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       entry.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
