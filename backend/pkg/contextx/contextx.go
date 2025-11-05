package contextx

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const UserContextKey contextKey = "user_context"

// UserContext holds lightweight user info from JWT
type UserContext struct {
	UserID uuid.UUID
	Email  string
	Roles  []string
}

// Helper: extract UserContext from request context
func Get(ctx context.Context) (*UserContext, bool) {
	val := ctx.Value(UserContextKey)
	if val == nil {
		return nil, false
	}
	uc, ok := val.(*UserContext)
	return uc, ok
}

// Helper: return UserID pointer if present
func GetUserID(ctx context.Context) *uuid.UUID {
	if uc, ok := Get(ctx); ok {
		return &uc.UserID
	}
	return nil
}
