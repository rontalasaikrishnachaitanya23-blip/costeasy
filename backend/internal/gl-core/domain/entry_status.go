// backend/internal/gl-core/domain/entry_status.go
package domain

// EntryStatus represents the status of a journal entry
type EntryStatus string

const (
	EntryStatusDraft    EntryStatus = "DRAFT"    // Editable, not posted
	EntryStatusPosted   EntryStatus = "POSTED"   // Posted to ledger, immutable
	EntryStatusVoid     EntryStatus = "VOID"     // Voided entry
	EntryStatusReversed EntryStatus = "REVERSED" // Reversed entry
)

// IsParentAccount checks if this account is a parent/control account
// Parent accounts typically have shorter codes and should not be used for direct posting
func (a *GLAccount) IsParentAccount() bool {
	// Simple heuristic: parent accounts have shorter codes
	// Adjust based on your account numbering scheme:
	// - "1000" = parent (4 digits)
	// - "1000-001" = child (8+ characters)
	// - "1100" = parent (4 digits)
	// - "1100.01" = child (7+ characters)

	return len(a.Code) <= 4
}

// IsValid checks if the status is valid
func (es EntryStatus) IsValid() bool {
	validStatuses := map[EntryStatus]bool{
		EntryStatusDraft:    true,
		EntryStatusPosted:   true,
		EntryStatusVoid:     true,
		EntryStatusReversed: true,
	}
	return validStatuses[es]
}

// String returns the string representation
func (es EntryStatus) String() string {
	return string(es)
}

// CanTransitionTo checks if status can transition to another status
func (es EntryStatus) CanTransitionTo(newStatus EntryStatus) bool {
	validTransitions := map[EntryStatus][]EntryStatus{
		EntryStatusDraft: {
			EntryStatusPosted,
			EntryStatusVoid,
		},
		EntryStatusPosted: {
			EntryStatusVoid,
			EntryStatusReversed,
		},
		EntryStatusVoid:     {}, // No transitions from VOID
		EntryStatusReversed: {}, // No transitions from REVERSED
	}

	allowedStatuses, exists := validTransitions[es]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

// IsFinal checks if this is a final status (no further transitions)
func (es EntryStatus) IsFinal() bool {
	return es == EntryStatusVoid || es == EntryStatusReversed
}

// IsEditable checks if entry with this status can be edited
func (es EntryStatus) IsEditable() bool {
	return es == EntryStatusDraft
}

// IsPosted checks if entry with this status is posted
func (es EntryStatus) IsPosted() bool {
	return es == EntryStatusPosted || es == EntryStatusReversed
}
