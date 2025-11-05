package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/settings/handler"
	"github.com/gin-gonic/gin"
)

func RegisterShafafiyaRoutes(router *gin.RouterGroup, handler *handler.ShafafiyaHandler) {
	shafafiya := router.Group("/shafafiya")
	{
		// Configuration/Setup routes only
		shafafiya.POST("/organizations/:org_id", handler.CreateShafafiyaSettings)
		shafafiya.GET("/organizations/:org_id", handler.GetShafafiyaSettings)
		shafafiya.PUT("/organizations/:org_id", handler.UpdateShafafiyaSettings)
		shafafiya.DELETE("/organizations/:org_id", handler.DeleteShafafiyaSettings)

		// Validation endpoint for testing configuration
		shafafiya.POST("/validate/:org_id", handler.ValidateShafafiyaConfig)

		// List all configurations (for admin)
		shafafiya.GET("/", handler.ListShafafiyaSettings)
	}
}
