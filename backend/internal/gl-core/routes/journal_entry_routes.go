// backend/internal/gl-core/routes/journal_entry_routes.go
package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler"
	"github.com/gin-gonic/gin"
)

// RegisterJournalEntryRoutes registers all journal entry routes
func RegisterJournalEntryRoutes(r *gin.RouterGroup, h *handler.JournalEntryHandler) {
	journalEntries := r.Group("/journal-entries")
	{
		journalEntries.POST("", h.CreateJournalEntry)              // Create entry
		journalEntries.GET("", h.ListJournalEntries)               // List entries
		journalEntries.GET("/:id", h.GetJournalEntry)              // Get entry by ID
		journalEntries.PUT("/:id", h.UpdateJournalEntry)           // Update entry
		journalEntries.DELETE("/:id", h.DeleteJournalEntry)        // Delete draft entry
		journalEntries.POST("/:id/post", h.PostJournalEntry)       // Post entry
		journalEntries.POST("/:id/void", h.VoidJournalEntry)       // Void entry
		journalEntries.POST("/:id/reverse", h.ReverseJournalEntry) // Reverse entry
	}
}
