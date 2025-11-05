package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/chaitu35/costeasy/backend/internal/auth/domain"
	"github.com/chaitu35/costeasy/backend/internal/auth/handler/dto"
	"github.com/chaitu35/costeasy/backend/internal/auth/repository"
	"github.com/chaitu35/costeasy/backend/pkg/jwt"
)

type AuthService struct {
	userRepo    repository.UserRepositoryInterface
	roleRepo    repository.RoleRepositoryInterface
	permRepo    repository.PermissionRepositoryInterface
	refreshRepo repository.RefreshTokenRepositoryInterface // ✅ add this
	jwtUtil     *jwt.JWTUtil
}

func NewAuthService(
	userRepo repository.UserRepositoryInterface,
	roleRepo repository.RoleRepositoryInterface,
	permRepo repository.PermissionRepositoryInterface,
	refreshRepo repository.RefreshTokenRepositoryInterface, // ✅ add this
	jwtUtil *jwt.JWTUtil,
) AuthServiceInterface {
	return &AuthService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		permRepo:    permRepo,
		refreshRepo: refreshRepo, // ✅ assign here
		jwtUtil:     jwtUtil,
	}
}

//
// ──────────────────────────────────────────────
//  AUTHENTICATION
// ──────────────────────────────────────────────
//

// Register a new user
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:                uuid.New(),
		Email:             req.Email,
		Username:          req.Username,
		PasswordHash:      string(hashed),
		IsActive:          true,
		IsVerified:        false,
		CreatedAt:         now,
		UpdatedAt:         now,
		MFAEnabled:        false,
		MFAMethod:         domain.MFAMethodNone,
		AllowRemoteAccess: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	tokenPair, err := s.jwtUtil.GenerateTokenPair(user.ID, user.Email, nil, nil)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		IsActive: user.IsActive,
		Token:    tokenPair.AccessToken,
	}, nil
}

// Login authenticates an existing user
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Account lock and status check
	if !user.CanLogin() {
		return nil, errors.New("account locked or inactive")
	}

	// Password verification
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		user.IncrementFailedAttempts()
		_ = s.userRepo.Update(ctx, user)
		return nil, errors.New("invalid credentials")
	}

	user.ResetFailedAttempts()
	user.UpdateLastLogin(req.IPAddress)
	_ = s.userRepo.Update(ctx, user)

	// Future: handle MFA (if user.MFAEnabled)

	// Load user permissions
	userPerms, _ := s.permRepo.GetUserPermissions(ctx, user.ID)
	var permClaims []jwt.PermissionClaim
	for _, p := range userPerms {
		claim := jwt.PermissionClaim{
			Module: p.Module,
			Page:   p.Resource, // map Resource -> Page
		}
		switch p.Action {
		case "view":
			claim.CanView = true
		case "create":
			claim.CanAdd = true
		case "edit":
			claim.CanEdit = true
		case "delete":
			claim.CanDelete = true
		case "print":
			claim.CanPrint = true
		case "export":
			claim.CanExport = true
		}
		permClaims = append(permClaims, claim)
	}

	// Generate JWT pair
	tokens, err := s.jwtUtil.GenerateTokenPair(user.ID, user.Email, nil, permClaims)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtUtil.AccessExpiry().Seconds()),
		User:         dto.ToUserResponse(user),
	}, nil
}

// RefreshToken renews tokens using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenResponse, error) {
	claims, err := s.jwtUtil.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}

	newPair, err := s.jwtUtil.GenerateTokenPair(userID, "", nil, nil)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		AccessToken:  newPair.AccessToken,
		RefreshToken: newPair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwtUtil.AccessExpiry().Seconds()),
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	// Optional: blacklist tokens
	return nil
}

//
// ──────────────────────────────────────────────
//  USER MANAGEMENT
// ──────────────────────────────────────────────
//

func (s *AuthService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToUserResponse(user), nil
}

func (s *AuthService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashed),
		IsActive:     true,
		IsVerified:   true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return dto.ToUserResponse(user), nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != "" {
		user.FirstName = &req.FirstName
	}
	if req.LastName != "" {
		user.LastName = &req.LastName
	}
	if req.Phone != "" {
		user.Phone = &req.Phone
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return dto.ToUserResponse(user), nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("old password incorrect")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(newHash)
	user.PasswordChangedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) ListUsers(ctx context.Context, orgID uuid.UUID, limit, offset int) ([]*dto.UserResponse, error) {
	users, err := s.userRepo.List(ctx, orgID, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []*dto.UserResponse
	for _, u := range users {
		res = append(res, dto.ToUserResponse(u))
	}
	return res, nil
}

func (s *AuthService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.Delete(ctx, userID)
}

//
// ──────────────────────────────────────────────
//  ROLE MANAGEMENT
// ──────────────────────────────────────────────
//

func (s *AuthService) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleRepo.AssignToUser(ctx, userID, roleID)
}

func (s *AuthService) ChangeRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleRepo.ChangeUserRole(ctx, userID, roleID)
}

func (s *AuthService) RemoveRole(ctx context.Context, userID, roleID uuid.UUID) error {
	return s.roleRepo.RemoveUserRole(ctx, userID, roleID)
}

func (s *AuthService) CreateRole(ctx context.Context, req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	role := &domain.Role{
		ID:             uuid.New(),
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Description:    &req.Description,
		OrganizationID: uuid.Nil,
		IsSystemRole:   false,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if req.OrganizationID != nil {
		role.OrganizationID = *req.OrganizationID
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return dto.ToRoleResponse(role, nil), nil
}

func (s *AuthService) UpdateRole(ctx context.Context, roleID uuid.UUID, req *dto.UpdateRoleRequest) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name
	role.DisplayName = req.DisplayName
	role.Description = &req.Description
	role.UpdatedAt = time.Now()

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return dto.ToRoleResponse(role, nil), nil
}

func (s *AuthService) GetRole(ctx context.Context, roleID uuid.UUID) (*dto.RoleResponse, error) {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return dto.ToRoleResponse(role, nil), nil
}

func (s *AuthService) ListRoles(ctx context.Context, orgID *uuid.UUID) ([]*dto.RoleResponse, error) {
	roles, err := s.roleRepo.List(ctx, *orgID, 0, 0)
	if err != nil {
		return nil, err
	}

	var res []*dto.RoleResponse
	for _, r := range roles {
		res = append(res, dto.ToRoleResponse(r, nil))
	}
	return res, nil
}

func (s *AuthService) DeleteRole(ctx context.Context, roleID uuid.UUID) error {
	return s.roleRepo.Delete(ctx, roleID)
}

//
// ──────────────────────────────────────────────
//  PERMISSIONS
// ──────────────────────────────────────────────
//

func (s *AuthService) AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	for _, pid := range permissionIDs {
		if err := s.roleRepo.AssignPermission(ctx, roleID, pid, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) RemovePermissionsFromRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	for _, pid := range permissionIDs {
		if err := s.roleRepo.RemovePermission(ctx, roleID, pid); err != nil {
			return err
		}
	}
	return nil
}

func (s *AuthService) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]*dto.PermissionResponse, error) {
	perms, err := s.permRepo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return nil, err
	}

	var res []*dto.PermissionResponse
	for _, p := range perms {
		res = append(res, &dto.PermissionResponse{
			Module:    p.Module,
			Page:      p.Resource,
			CanView:   p.Action == "view",
			CanAdd:    p.Action == "create",
			CanEdit:   p.Action == "edit",
			CanDelete: p.Action == "delete",
			CanPrint:  p.Action == "print",
			CanExport: p.Action == "export",
		})
	}
	return res, nil
}

//
// ──────────────────────────────────────────────
//  PERMISSION CHECK
// ──────────────────────────────────────────────
//

func (s *AuthService) CheckPermission(ctx context.Context, userID uuid.UUID, module, resource, action string) (bool, error) {
	perms, err := s.permRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, p := range perms {
		if p.Module == module && p.Resource == resource && p.Action == action {
			return true, nil
		}
	}
	return false, nil
}
