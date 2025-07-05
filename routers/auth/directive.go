package auth

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
)

func AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver, msg *string) (interface{}, error) {
	user := GetCtxUserId(ctx)
	if user == 0 {
		if msg != nil {
			return nil, errors.New(*msg)
		}
		return nil, errors.New("unauthorized")
	}
	return next(ctx)
}
