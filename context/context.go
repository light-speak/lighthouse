package context

import (
	"context"

	"github.com/light-speak/lighthouse/errors"
)

type Context struct {
	context.Context
	UserId int64
	Errors []errors.GraphqlErrorInterface
}

func NewContext(ctx context.Context) *Context {
	return &Context{Context: ctx}
}
