package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
)

func handlerScopes(o *ast.ObjectNode, d *ast.Directive, store *ast.NodeStore) errors.GraphqlErrorInterface {
	names := d.GetArg("names").Value.([]interface{})
	for _, name := range names {
		o.Scopes = append(o.Scopes, name.(string))
	}
	log.Warn().Interface("names", o.Scopes).Msg("scopes")
	return nil
}

func init() {
	ast.AddObjectDirective("scopes", handlerScopes)
}
