// app/middleware/ip_whitelist_with_user_toggle.go
package middleware

import (
	"net"
	"net/http"
	"strings"

	//"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/chaitu35/costeasy/backend/internal/auth/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// IPWhitelistWithUserToggleMiddleware checks IP restrictions per user
// Users with allow_remote_access=true can access from anywhere
// Users with allow_remote_access=false must be in office
func IPWhitelistWithUserToggleMiddleware(
	config IPWhitelistConfig,
	userRepo repository.UserRepositoryInterface,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if IP whitelist is completely disabled
		if !config.Enabled {
			c.Next()
			return
		}

		// Get user ID from context (set by AuthMiddleware)
		userIDStr, exists := c.Get("user_id")
		if !exists {
			// No user in context, apply global IP restriction
			enforceIPRestriction(c, config)
			return
		}

		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			enforceIPRestriction(c, config)
			return
		}

		// Get user from database to check remote access flag
		user, err := userRepo.GetByID(c.Request.Context(), userID)
		if err != nil {
			// Error getting user, apply restriction for safety
			enforceIPRestriction(c, config)
			return
		}

		// KEY CHECK: Does user have remote access enabled?
		if user.AllowRemoteAccess {
			// User can access from anywhere - skip IP check
			c.Set("client_ip", getClientIP(c))
			c.Set("remote_access_allowed", true)
			c.Next()
			return
		}

		// User does NOT have remote access - enforce IP restriction
		c.Set("remote_access_allowed", false)
		enforceIPRestriction(c, config)
	}
}

// enforceIPRestriction checks if client IP is allowed
func enforceIPRestriction(c *gin.Context, config IPWhitelistConfig) {
	clientIP := getClientIP(c)
	if clientIP == "" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Unable to determine client IP address",
		})
		c.Abort()
		return
	}

	// Check if IP is allowed
	if isIPAllowed(clientIP, config) {
		c.Set("client_ip", clientIP)
		c.Next()
		return
	}

	// IP not allowed - block access
	c.JSON(http.StatusForbidden, gin.H{
		"error":   "Access denied: You can only access this system from office premises",
		"ip":      clientIP,
		"message": "Contact your administrator to enable remote access",
	})
	c.Abort()
}

// getClientIP extracts client IP from request
func getClientIP(c *gin.Context) string {
	// Try X-Forwarded-For header (if behind proxy/load balancer)
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Try X-Real-IP header
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}

	return ip
}

// isIPAllowed checks if client IP is in whitelist
func isIPAllowed(clientIP string, config IPWhitelistConfig) bool {
	ip := net.ParseIP(clientIP)
	if ip == nil {
		return false
	}

	// Check if localhost and bypass is enabled
	if config.BypassForLocal && isLocalhost(ip) {
		return true
	}

	// Check individual IPs
	for _, allowedIP := range config.AllowedIPs {
		if clientIP == allowedIP {
			return true
		}
	}

	// Check IP ranges (CIDR)
	for _, cidr := range config.AllowedRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// isLocalhost checks if IP is localhost
func isLocalhost(ip net.IP) bool {
	return ip.IsLoopback() || ip.Equal(net.ParseIP("127.0.0.1")) || ip.Equal(net.ParseIP("::1"))
}

// IPWhitelistConfig holds IP whitelist configuration
type IPWhitelistConfig struct {
	Enabled        bool
	AllowedIPs     []string
	AllowedRanges  []string
	BypassForLocal bool
}
