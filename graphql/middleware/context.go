package middleware

import (
	"context"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/log"
	"net/http"
	"strconv"
)

func ContextMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), graphql.ContextKey, &graphql.Context{})

			if userID := r.Header.Get("X-User-Id"); userID != "" {
				if userIDInt64, err := strconv.ParseInt(userID, 10, 64); err == nil {
					ctx = context.WithValue(ctx, userIDContextKey, userIDInt64)
				} else {
					log.Error("解析用户ID失败: %v", err)
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetContext(ctx context.Context) *graphql.Context {
	if gqlCtx, ok := ctx.Value(graphql.ContextKey).(*graphql.Context); ok {
		return gqlCtx
	}
	return &graphql.Context{}
}
