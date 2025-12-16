package auth

import (
	"context"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/light-speak/lighthouse/logs"
)

var userContextKey = &contextKey{"user"}
var sessionContextKey = &contextKey{"session"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token != "" {
				// Remove "Bearer " prefix if present
				if len(token) > 7 && token[:7] == "Bearer " {
					token = token[7:]
				}
				userId, err := GetUserId(token)
				logs.Debug().Msgf("request token: %s, user id: %d", token, userId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userContextKey, uint(userId))
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func AdminAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Header.Get("X-Session-Id")
			if session != "" {
				ctx := context.WithValue(r.Context(), sessionContextKey, session)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func XUserMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Header.Get("X-User-Id")
			if userId != "" {
				userId, err := strconv.ParseUint(userId, 10, 64)
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userContextKey, uint(userId))
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func WebSocketInitFunc(ctx context.Context, initPayload transport.InitPayload) (context.Context, *transport.InitPayload, error) {
	if authHeader, ok := initPayload["Authorization"].(string); ok {
		token := authHeader
		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
		userId, err := GetUserId(token)
		if err != nil {
			return ctx, nil, err
		}
		logs.Debug().Msgf("init payload: %v, user id: %d", initPayload, userId)
		ctx = context.WithValue(ctx, userContextKey, uint(userId))
	}
	if userIdStr, ok := initPayload["X-User-Id"].(string); ok {
		userId, err := strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			return ctx, nil, err
		}
		logs.Debug().Msgf("init payload: %v, user id: %d", initPayload, userId)
		ctx = context.WithValue(ctx, userContextKey, uint(userId))
	}
	return ctx, &initPayload, nil
}

func GetCtxUserId(ctx context.Context) uint {
	if userId, ok := ctx.Value(userContextKey).(uint); ok {
		return userId
	}
	return 0
}

func IsLogin(ctx context.Context) bool {
	return GetCtxUserId(ctx) != 0
}

func IsCurrentUser(ctx context.Context, userId uint) bool {
	return GetCtxUserId(ctx) == userId
}
