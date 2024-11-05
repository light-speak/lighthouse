package middleware

import (
	"net"
	"net/http"

	"github.com/light-speak/lighthouse/context"
)

func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.NewContext(r.Context())
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ctx.RemoteAddr = &host
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
