// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package resolver

import (
  "github.com/light-speak/lighthouse/resolve"
  "sync"
  "github.com/light-speak/lighthouse/context"
  "github.com/light-speak/lighthouse/graphql/excute"
)

func init() {
	excute.AddAttrResolver("FuckingAttrAttr", func(ctx *context.Context, data *sync.Map, resolve resolve.Resolve) (interface{}, error) {
		r := resolve.(*Resolver)
		res, err := r.FuckingAttrAttrResolver(ctx, data)
		return res, err
	})
}