// routes/account_routes.go
package routes

import (
	handlers "github.com/chaitu35/costeasy/backend/internal/gl-core/handler"
	"github.com/gin-gonic/gin"
)

func RegisterAccountRoutes(router *gin.RouterGroup, handler *handlers.AccountHandler) {
	accounts := router.Group("/accounts")
	{
		accounts.POST("", handler.CreateAccount)
		accounts.GET("", handler.ListAccounts)
		accounts.GET("/search", handler.SearchAccounts)
		accounts.GET("/code/:code", handler.GetAccountByCode)
		accounts.GET("/:id", handler.GetAccount)
		accounts.PUT("/:id", handler.UpdateAccount)
		accounts.DELETE("/:id", handler.SoftDeleteAccount)
		accounts.POST("/:id/activate", handler.ActivateAccount)
		accounts.POST("/:id/deactivate", handler.DeactivateAccount)
	}
}
