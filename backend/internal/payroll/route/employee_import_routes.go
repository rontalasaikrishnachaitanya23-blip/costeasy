package routes

import (
	"github.com/chaitu35/costeasy/backend/internal/payroll/handler"
	"github.com/gin-gonic/gin"
)

func RegisterEmployeeImportRoutes(rg *gin.RouterGroup, h *handler.EmployeeImportHandler) {
	grp := rg.Group("/employees/import")
	{
		grp.GET("/template", h.DownloadTemplate)
		grp.POST("", h.Import)
	}
}
