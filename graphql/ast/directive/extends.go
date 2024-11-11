package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

func handlerExtends(o *ast.ObjectNode, d *ast.Directive, store *ast.NodeStore) errors.GraphqlErrorInterface {
	o.IsExtend = true
	return nil
}

func init() {
	ast.AddObjectDirective("extends", handlerExtends)
}
