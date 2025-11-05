package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/chaitu35/costeasy/backend/internal/auth/service"
	"github.com/chaitu35/costeasy/backend/pkg/contextx"
	"github.com/chaitu35/costeasy/backend/pkg/jwt"
)

type AuthMiddleware struct {
	jwtUtil     *jwt.JWTUtil
	authService service.AuthServiceInterface
}

func NewAuthMiddleware(jwtUtil *jwt.JWTUtil, authService service.AuthServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		jwtUtil:     jwtUtil,
		authService: authService,
	}
}

// Authenticate validates JWT token and sets user context
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.jwtUtil.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// ✅ Parse user ID from JWT claims
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// ✅ Build user context
		userCtx := &contextx.UserContext{
			UserID: userID,
			Email:  claims.Email,
			Roles:  claims.Roles,
		}

		// ✅ Store in both Gin and Request Context
		c.Set("user", userCtx)
		ctx := context.WithValue(c.Request.Context(), contextx.UserContextKey, userCtx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// RequirePermission checks if user has a specific permission
func (m *AuthMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, ok := contextx.Get(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		hasPermission, err := m.authService.CheckPermission(ctx, userCtx.UserID, resource, action, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission check failed"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole checks if user has a specific role
func (m *AuthMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, ok := contextx.Get(c.Request.Context())
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// ✅ Direct role check
		for _, r := range userCtx.Roles {
			if r == roleName {
				c.Next()
				return
			}
		}

		// Optional: fallback to permission-based check
		ctx := c.Request.Context()
		hasPermission, err := m.authService.CheckPermission(ctx, userCtx.UserID, "system", "admin", "")
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role: " + roleName + " required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
