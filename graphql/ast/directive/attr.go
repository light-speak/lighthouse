package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

func Attr(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	f.IsAttr = true
	return nil
}

func init() {
	ast.AddFieldDirective("attr", Attr)
}
