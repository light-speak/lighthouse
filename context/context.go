package context

import (
	"context"

	"github.com/light-speak/lighthouse/errors"
)

type Context struct {
	context.Context `json:"-"`
	UserId          *int64                         `json:"userId"`
	RemoteAddr      *string                        `json:"remoteAddr"`
	Errors          []errors.GraphqlErrorInterface `json:"-"`
	Inject          map[string]interface{}         `json:"-"`
}

func NewContext(ctx context.Context) *Context {
	return &Context{Context: ctx}
}
