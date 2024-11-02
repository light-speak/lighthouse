package middleware

import (
	"net/http"

	"github.com/light-speak/lighthouse/env"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if env.LighthouseConfig.App.Mode == env.Single {
			//TODO: 本地验证
		} else {
			//TODO: 从router获取
		}
		next.ServeHTTP(w, r)
	})
}
