// backend/settings/internal/routes/shafafiya_routes.go
package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/settings/handler"
	"github.com/gin-gonic/gin"
)

// RegisterShafafiyaRoutes registers all Shafafiya-related routes
func RegisterShafafiyaRoutes(router *gin.RouterGroup, handler *handler.ShafafiyaHandler) {
	// Failed submissions endpoint (global)
	router.GET("/shafafiya/failed-submissions", handler.ListFailedSubmissions)

	// Organization-specific Shafafiya settings
	orgs := router.Group("/organizations/:org_id/shafafiya")
	{
		// CRUD operations
		orgs.POST("", handler.CreateShafafiyaSettings)
		orgs.GET("", handler.GetShafafiyaSettings)
		orgs.DELETE("", handler.DeleteShafafiyaSettings)

		// Configuration updates
		orgs.PUT("/credentials", handler.UpdateShafafiyaCredentials)
		orgs.PUT("/costing", handler.UpdateShafafiyaCosting)
		orgs.PUT("/submission", handler.UpdateShafafiyaSubmission)

		// Validation
		orgs.GET("/validate", handler.ValidateShafafiyaConfiguration)
	}
}
