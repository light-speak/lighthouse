package middleware

import (
	"context"
	"github.com/light-speak/lighthouse"
	"net/http"
)

func ContextMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			//TODO: Authorization
			ctx := context.WithValue(request.Context(), lighthouse.ContextKey, &lighthouse.Context{})

			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request)
		})
	}
}

func GetContext(ctx context.Context) *lighthouse.Context {
	return ctx.Value(lighthouse.ContextKey).(*lighthouse.Context)
}
