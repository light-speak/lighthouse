package middleware

import (
	"context"
)

type contextKey struct {
	name string
}

var userIdContextKey = &contextKey{
	name: "USER",
}

func UserId(ctx context.Context) (int64, error) {
	userId, ok := ctx.Value(userIdContextKey).(int64)
	if !ok {
		return 0, &AuthError{Message: "请登录后重试"}
	}
	return userId, nil
}

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
