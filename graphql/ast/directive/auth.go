package directive

import (
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

func Auth(ctx *context.Context, f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node, result interface{}) errors.GraphqlErrorInterface {
	msg := env.LighthouseConfig.Auth.UnauthorizedMessage
	if msgArg := d.GetArg("msg"); msgArg != nil {
		msg = msgArg.Value.(string)
	}
	if ctx.UserId == nil {
		return &errors.GraphQLError{
			Message:   msg,
			Locations: []*errors.GraphqlLocation{f.GetLocation()},
		}
	}
	return nil
}

func init() {
	ast.AddFieldRuntimeBeforeDirective("auth", Auth)
}
