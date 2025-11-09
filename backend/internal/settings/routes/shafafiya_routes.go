package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/settings/handler"
	"github.com/gin-gonic/gin"
)

// RegisterShafafiyaRoutes registers all Shafafiya-related API routes.
func RegisterShafafiyaRoutes(router *gin.RouterGroup, handler *handler.ShafafiyaHandler) {
	shafafiya := router.Group("/shafafiya")

	// ─────────────────────────────────────────────
	// 🔹 CRUD: Shafafiya Settings per Organization
	// ─────────────────────────────────────────────
	{
		shafafiya.POST("/organizations/:org_id", handler.CreateShafafiyaSettings)
		shafafiya.GET("/organizations/:org_id", handler.GetShafafiyaSettings)
		shafafiya.PUT("/organizations/:org_id", handler.UpdateShafafiyaSettings)
		shafafiya.DELETE("/organizations/:org_id", handler.DeleteShafafiyaSettings)
	}

	// ─────────────────────────────────────────────
	// 🔹 Partial Updates for Credentials / Costing / Submission
	// ─────────────────────────────────────────────
	{
		shafafiya.PUT("/organizations/:org_id/credentials", handler.UpdateShafafiyaCredentials)
		shafafiya.PUT("/organizations/:org_id/costing", handler.UpdateShafafiyaCosting)
		shafafiya.PUT("/organizations/:org_id/submission", handler.UpdateShafafiyaSubmission)
	}

	// ─────────────────────────────────────────────
	// 🔹 Validation, Listing, & Monitoring
	// ─────────────────────────────────────────────
	{
		shafafiya.GET("/organizations/:org_id/validate", handler.ValidateShafafiyaConfiguration)
		shafafiya.GET("/", handler.ListShafafiyaSettings)
		shafafiya.GET("/failed-submissions", handler.ListFailedSubmissions)
	}
}
