package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/chaitu35/costeasy/backend/app/middleware"
    "github.com/chaitu35/costeasy/backend/internal/auth/handler"
)

// RegisterAuthRoutes registers all authentication routes
func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
    auth := router.Group("/auth")
    {
        // Public routes (no authentication required)
        auth.POST("/register", authHandler.Register)
        auth.POST("/login", authHandler.Login)
        auth.POST("/refresh", authHandler.RefreshToken)
        auth.POST("/forgot-password", authHandler.ForgotPassword)
        auth.POST("/reset-password", authHandler.ResetPassword)
        
        // Protected routes (authentication required)
        protected := auth.Group("")
        protected.Use(authMiddleware.Authenticate())
        {
            protected.POST("/logout", authHandler.Logout)
            protected.POST("/change-password", authHandler.ChangePassword)
            protected.GET("/me", authHandler.GetCurrentUser)
            protected.PUT("/me", authHandler.UpdateProfile)
            
            // MFA routes
            protected.POST("/mfa/enable", authHandler.EnableMFA)
            protected.POST("/mfa/disable", authHandler.DisableMFA)
            protected.POST("/mfa/verify", authHandler.VerifyMFA)
        }
    }
    
    // User management routes (admin only)
    users := router.Group("/users")
    users.Use(authMiddleware.Authenticate())
    users.Use(authMiddleware.RequirePermission("users", "view"))
    {
        users.GET("", authHandler.ListUsers)
        users.GET("/:id", authHandler.GetUser)
        users.POST("", authMiddleware.RequirePermission("users", "add"), authHandler.CreateUser)
        users.PUT("/:id", authMiddleware.RequirePermission("users", "edit"), authHandler.UpdateUser)
        users.DELETE("/:id", authMiddleware.RequirePermission("users", "delete"), authHandler.DeleteUser)
        users.POST("/:id/roles", authMiddleware.RequirePermission("users", "edit"), authHandler.AssignRole)
        users.DELETE("/:id/roles/:roleId", authMiddleware.RequirePermission("users", "edit"), authHandler.RemoveRole)
    }
    
    // Role management routes (admin only)
    roles := router.Group("/roles")
    roles.Use(authMiddleware.Authenticate())
    roles.Use(authMiddleware.RequirePermission("roles", "view"))
    {
        roles.GET("", authHandler.ListRoles)
        roles.GET("/:id", authHandler.GetRole)
        roles.POST("", authMiddleware.RequirePermission("roles", "add"), authHandler.CreateRole)
        roles.PUT("/:id", authMiddleware.RequirePermission("roles", "edit"), authHandler.UpdateRole)
        roles.DELETE("/:id", authMiddleware.RequirePermission("roles", "delete"), authHandler.DeleteRole)
        roles.POST("/:id/permissions", authMiddleware.RequirePermission("roles", "edit"), authHandler.AssignPermission)
        roles.DELETE("/:id/permissions/:permissionId", authMiddleware.RequirePermission("roles", "edit"), authHandler.RemovePermission)
    }
}
