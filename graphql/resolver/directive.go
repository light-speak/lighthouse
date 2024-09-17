package resolver

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	gql "github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/graphql/middleware"
)

func getFieldKeyValue(ctx context.Context, obj interface{}) (string, string, interface{}, error) {
	if data, ok := obj.(map[string]interface{}); ok {
		p := graphql.GetPathContext(ctx).ParentField.Path()
		k := *graphql.GetPathContext(ctx).Field
		v := data[k]
		return p.String(), k, v, nil
	}
	return "", "", nil, errors.New("parse key value error")
}

var First = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return next(ctx)
}

var All = func(ctx context.Context, obj interface{}, next graphql.Resolver, scopes []*string) (res interface{}, err error) {
	return next(ctx)
}

var Scope = func(ctx context.Context, obj interface{}, next graphql.Resolver, scope string) (res interface{}, err error) {
	return next(ctx)
}

var Eq = func(ctx context.Context, obj interface{}, next graphql.Resolver, key *string) (interface{}, error) {
	lctx := middleware.GetContext(ctx)

	if p, k, v, err := getFieldKeyValue(ctx, obj); err != nil {
		return nil, err
	} else {
		if key != nil {
			k = *key
		}
		query := fmt.Sprintf("%s = ?", k)
		where := &gql.Where{
			Path:  p,
			Query: query,
			Value: v,
		}
		lctx.Wheres = append(lctx.Wheres, where)
	}

	return next(ctx)
}

var CreateOrUpdate = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	lctx := middleware.GetContext(ctx)
	args := graphql.GetFieldContext(ctx).Args
	lctx.Data = &args

	return next(ctx)
}

var Page = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	lctx := middleware.GetContext(ctx)
	if _, _, page, err := getFieldKeyValue(ctx, obj); err != nil {
		return nil, err
	} else {
		if lctx.Paginate == nil {
			lctx.Paginate = &gql.Paginate{Page: 1, Size: 10}
		}
		lctx.Paginate.Page = page.(int64)
	}
	return next(ctx)
}

var Size = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	lctx := middleware.GetContext(ctx)
	if _, _, size, err := getFieldKeyValue(ctx, obj); err != nil {
		return nil, err
	} else {
		if lctx.Paginate == nil {
			lctx.Paginate = &gql.Paginate{Page: 1, Size: 10}
		}
		lctx.Paginate.Size = size.(int64)
	}
	return next(ctx)
}

var Sum = func(ctx context.Context, obj interface{}, next graphql.Resolver, model string, column string, scopes []*string) (interface{}, error) {
	lctx := middleware.GetContext(ctx)
	lctx.Column = &column
	return next(ctx)
}

var Count = func(ctx context.Context, obj interface{}, next graphql.Resolver, model string, scopes []*string) (interface{}, error) {
	return next(ctx)
}

var Resolve = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	return next(ctx)
}

var Inject = func(ctx context.Context, obj interface{}, next graphql.Resolver, field string, target string) (interface{}, error) {
	return next(ctx)
}
