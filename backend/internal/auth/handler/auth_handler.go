package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/chaitu35/costeasy/backend/internal/auth/handler/dto"
	"github.com/chaitu35/costeasy/backend/internal/auth/service"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangePassword(c.Request.Context(), uid, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// ForgotPassword handles forgot password request
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement forgot password logic (send email with reset token)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset instructions sent to email"})
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement reset password logic (validate token and reset password)
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get authenticated user ID from context (set by middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convert to UUID
	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse request body
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update profile for this specific user
	updatedUser, err := h.authService.UpdateProfile(c.Request.Context(), uid, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// EnableMFA enables MFA for user
func (h *AuthHandler) EnableMFA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement MFA enable logic
	c.JSON(http.StatusOK, gin.H{"message": "MFA enabled successfully", "user_id": uid})
}

// DisableMFA disables MFA for user
func (h *AuthHandler) DisableMFA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement MFA disable logic
	c.JSON(http.StatusOK, gin.H{"message": "MFA disabled successfully", "user_id": uid})
}

// VerifyMFA verifies MFA code
func (h *AuthHandler) VerifyMFA(c *gin.Context) {
	var req dto.VerifyMFARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement MFA verification logic
	c.JSON(http.StatusOK, gin.H{"message": "MFA verified successfully"})
}

// User Management Methods

// ListUsers lists all users
func (h *AuthHandler) ListUsers(c *gin.Context) {
	// TODO: Implement list users with pagination
	c.JSON(http.StatusOK, gin.H{"users": []interface{}{}})
}

// GetUser gets a specific user
func (h *AuthHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser updates a user
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement update user logic
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser deletes a user
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// TODO: Implement delete user logic
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// AssignRole assigns a role to a user (can have multiple roles)
func (h *AuthHandler) AssignRole(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		RoleID uuid.UUID `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.AssignRole(c.Request.Context(), userID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// ChangeRole changes user's role (removes all existing, assigns new one)
func (h *AuthHandler) ChangeRole(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		RoleID uuid.UUID `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ChangeRole(c.Request.Context(), userID, req.RoleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role changed successfully"})
}

// RemoveRole removes a specific role from a user
func (h *AuthHandler) RemoveRole(c *gin.Context) {
	userIDParam := c.Param("id")
	roleIDParam := c.Param("roleId")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleID, err := uuid.Parse(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := h.authService.RemoveRole(c.Request.Context(), userID, roleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}

// Role Management Methods

// ListRoles lists all roles
func (h *AuthHandler) ListRoles(c *gin.Context) {
	// TODO: Implement list roles
	c.JSON(http.StatusOK, gin.H{"roles": []interface{}{}})
}

// GetRole gets a specific role
func (h *AuthHandler) GetRole(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// TODO: Implement get role logic
	c.JSON(http.StatusOK, gin.H{"message": "Role details"})
}

// CreateRole creates a new role
func (h *AuthHandler) CreateRole(c *gin.Context) {
	// TODO: Implement create role logic
	c.JSON(http.StatusCreated, gin.H{"message": "Role created successfully"})
}

// UpdateRole updates a role
func (h *AuthHandler) UpdateRole(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// TODO: Implement update role logic
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

// DeleteRole deletes a role
func (h *AuthHandler) DeleteRole(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// TODO: Implement delete role logic
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// AssignPermission assigns a permission to a role
func (h *AuthHandler) AssignPermission(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// TODO: Implement assign permission logic
	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned successfully"})
}

// RemovePermission removes a permission from a role
func (h *AuthHandler) RemovePermission(c *gin.Context) {
	idParam := c.Param("id")
	_, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	// TODO: Implement remove permission logic
	c.JSON(http.StatusOK, gin.H{"message": "Permission removed successfully"})
}
