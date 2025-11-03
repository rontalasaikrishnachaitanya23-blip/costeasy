package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/gl-core/handler"
	"github.com/gin-gonic/gin"
)

// RegisterImportRoutes registers import-related routes
func RegisterImportRoutes(router *gin.RouterGroup, importHandler *handler.ImportHandler) {
	imports := router.Group("/import")
	{
		// Import endpoints
		imports.POST("/accounts", importHandler.ImportChartOfAccounts)
		imports.POST("/journal-entries", importHandler.ImportJournalEntries)

		// Template download
		imports.GET("/template/:type", importHandler.DownloadTemplate)
	}
}
