package hidden

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func HiddenDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	next(ctx)
	return nil, nil
}
