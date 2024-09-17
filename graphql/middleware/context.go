package middleware

import (
	"context"
	"fmt"
	"github.com/light-speak/lighthouse/graphql"
	"net/http"
	"strconv"
)

func ContextMiddleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := context.WithValue(request.Context(), graphql.ContextKey, &graphql.Context{})

			userIdList := request.Header.Values("X-User-Id")
			if len(userIdList) > 0 {
				userId := userIdList[0]
				userIdInt64, err := strconv.ParseInt(userId, 10, 64)
				if err == nil {
					ctx = context.WithValue(ctx, userIdContextKey, userIdInt64)
				} else {
					fmt.Println(err)
				}
			}

			request = request.WithContext(ctx)
			next.ServeHTTP(writer, request)
		})
	}
}

func GetContext(ctx context.Context) *graphql.Context {
	return ctx.Value(graphql.ContextKey).(*graphql.Context)
}
