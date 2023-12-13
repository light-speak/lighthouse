package resolver

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/light-speak/lighthouse"
	"github.com/light-speak/lighthouse/middleware"
)

func getFieldKeyValue(ctx context.Context, obj interface{}) (string, interface{}, error) {

	if data, ok := obj.(map[string]interface{}); ok {
		k := *graphql.GetPathContext(ctx).Field
		v := data[k]
		return k, v, nil
	}
	return "", nil, errors.New("parse key value error")
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

	if k, v, err := getFieldKeyValue(ctx, obj); err != nil {
		return nil, err
	} else {
		if key != nil {
			k = *key
		}
		query := fmt.Sprintf("%s = ?", k)
		where := &lighthouse.Where{
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
	if _, page, err := getFieldKeyValue(ctx, obj); err != nil {
		return nil, err
	} else {
		if lctx.Paginate == nil {
			lctx.Paginate = &lighthouse.Paginate{Page: 0, Size: 0}
		}
		lctx.Paginate.Page = page.(int64)
	}
	return next(ctx)
}

var Size = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	lctx := middleware.GetContext(ctx)
	if _, size, err := getFieldKeyValue(ctx, obj); err != nil {
		return nil, err
	} else {
		if lctx.Paginate == nil {
			lctx.Paginate = &lighthouse.Paginate{Page: 0, Size: 0}
		}
		lctx.Paginate.Size = size.(int64)
	}
	return next(ctx)
}
