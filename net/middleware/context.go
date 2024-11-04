package middleware

import (
	"net/http"

	"github.com/light-speak/lighthouse/context"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(r.Context())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
