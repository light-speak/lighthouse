package middleware

import (
	"net/http"
	"strings"

	"github.com/light-speak/lighthouse/auth"
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/env"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		if env.LighthouseConfig.App.Mode == env.Single {
			userId, err := auth.GetUserId(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := r.Context().(*context.Context)
			ctx.UserId = &userId
		} else {
			//TODO: manor获取
		}
		next.ServeHTTP(w, r)
	})
}
