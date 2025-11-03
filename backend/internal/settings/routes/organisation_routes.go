// backend/settings/internal/routes/organization_routes.go
package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/settings/handler"
	"github.com/gin-gonic/gin"
)

// RegisterOrganizationRoutes registers all organization-related routes
func RegisterOrganizationRoutes(router *gin.RouterGroup, handler *handler.OrganizationHandler) {
	orgs := router.Group("/organizations")
	{
		// Statistics endpoint (must be before /:id to avoid conflict)
		orgs.GET("/stats", handler.GetOrganizationStats)

		// CRUD operations
		orgs.POST("", handler.CreateOrganization)
		orgs.GET("", handler.ListOrganizations)
		orgs.GET("/:id", handler.GetOrganization)
		orgs.PUT("/:id", handler.UpdateOrganization)
		orgs.DELETE("/:id", handler.DeactivateOrganization)

		// Activation
		orgs.POST("/:id/activate", handler.ActivateOrganization)
	}
}
