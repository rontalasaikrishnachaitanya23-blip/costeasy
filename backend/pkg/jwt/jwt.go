package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTUtil struct {
	accessSecret  string
	refreshSecret string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// TokenPair represents both access and refresh tokens.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type PermissionClaim struct {
	Module    string `json:"module"`
	Page      string `json:"page"`
	CanView   bool   `json:"can_view"`
	CanAdd    bool   `json:"can_add"`
	CanEdit   bool   `json:"can_edit"`
	CanDelete bool   `json:"can_delete"`
	CanPrint  bool   `json:"can_print"`
	CanExport bool   `json:"can_export"`
}

type AccessTokenClaims struct {
	UserID      string            `json:"user_id"`
	Email       string            `json:"email"`
	Roles       []string          `json:"roles"`
	Permissions []PermissionClaim `json:"permissions"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTUtil(secret, accessExpiry, refreshExpiry string) *JWTUtil {
	// Parse duration strings
	accessDuration, err := parseDuration(accessExpiry)
	if err != nil {
		accessDuration = 15 * time.Minute
	}

	refreshDuration, err := parseDuration(refreshExpiry)
	if err != nil {
		refreshDuration = 7 * 24 * time.Hour
	}

	return &JWTUtil{
		accessSecret:  secret,
		refreshSecret: secret, // Use same secret
		accessExpiry:  accessDuration,
		refreshExpiry: refreshDuration,
	}
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) > 1 && s[len(s)-1] == 'd' {
		days, err := time.ParseDuration(s[:len(s)-1] + "h")
		if err != nil {
			return 0, err
		}
		return days * 24, nil
	}
	return time.ParseDuration(s)
}

// GenerateAccessToken generates JWT access token with permissions
func (j *JWTUtil) GenerateAccessToken(userID uuid.UUID, email string, roles []string, permissions []PermissionClaim) (string, error) {
	now := time.Now()
	claims := AccessTokenClaims{
		UserID:      userID.String(),
		Email:       email,
		Roles:       roles,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.accessSecret))
}

// GenerateRefreshToken generates JWT refresh token
func (j *JWTUtil) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := RefreshTokenClaims{
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

// ValidateAccessToken validates and parses access token
func (j *JWTUtil) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.accessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateRefreshToken validates and parses refresh token
func (j *JWTUtil) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.refreshSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateTokenPair generates both access and refresh tokens.
func (j *JWTUtil) GenerateTokenPair(
	userID uuid.UUID,
	email string,
	roles []string,
	permissions []PermissionClaim,
) (*TokenPair, error) {
	access, err := j.GenerateAccessToken(userID, email, roles, permissions)
	if err != nil {
		return nil, err
	}
	refresh, err := j.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// AccessExpiry returns the configured access token expiry.
func (j *JWTUtil) AccessExpiry() time.Duration {
	return j.accessExpiry
}

// RefreshExpiry returns the configured refresh token expiry.
func (j *JWTUtil) RefreshExpiry() time.Duration {
	return j.refreshExpiry
}
