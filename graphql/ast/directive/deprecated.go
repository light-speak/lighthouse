package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

func handlerDeprecated(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	reason := d.GetArg("reason")
	if reason != nil {
		f.IsDeprecated = true
		v := "field is deprecated"
		if reason.Value != nil {
			v = reason.Value.(string)
		} else if reason.DefaultValue != nil {
			v = reason.DefaultValue.(string)
		}
		f.DeprecationReason = &v
	}
	return nil
}

func init() {
	ast.AddFieldDirective("deprecated", handlerDeprecated)
}
