package middleware

import (
	"context"
)

type contextKey string

const userIDContextKey contextKey = "USER"

func UserID(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDContextKey).(int64)
	if !ok {
		return 0, NewAuthError("请登录后重试")
	}
	return userID, nil
}

type AuthError struct {
	Message string
}

func NewAuthError(message string) *AuthError {
	return &AuthError{Message: message}
}

func (e *AuthError) Error() string {
	return e.Message
}
