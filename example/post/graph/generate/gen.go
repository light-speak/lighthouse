// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.

package generate

import "context"

type contextKey struct {
	name string
}

var loadersKey = &contextKey{
	name: "loadersKey",
}

type Loaders struct {
    FindPostById PostLoader
    
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
