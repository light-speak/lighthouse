package context

import (
	"context"

	"github.com/light-speak/lighthouse/errors"
)

type Context struct {
	context.Context `json:"-"`
	UserId          *int64                         `json:"userId"`
	RemoteAddr      *string                        `json:"remoteAddr"`
	OprationName    *string                        `json:"oprationName"`
	Errors          []errors.GraphqlErrorInterface `json:"-"`
	Inject          map[string]interface{}         `json:"-"`
}
