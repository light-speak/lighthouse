package directive

import (
	"fmt"
	"strings"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
)

func Cache(ctx *context.Context, field *ast.Field, directive *ast.Directive, store *ast.NodeStore, parent ast.Node, result interface{}) errors.GraphqlErrorInterface {
	log.Warn().Msgf("cache directive: %+v", result)
	var argsKey strings.Builder
	for _, arg := range field.Args {
		argsKey.WriteString(arg.Name)
		argsKey.WriteString(fmt.Sprintf("%v", arg.Value))
	}
	log.Warn().Msgf("cache directive argsKey: %s", argsKey.String())
	return nil
}

func init() {
	ast.AddFieldRuntimeAfterDirective("cache", Cache)
}
