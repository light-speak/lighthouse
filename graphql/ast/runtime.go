package ast

import (
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
)

// result is channel, so directive can set result to it
func ExecuteFieldBeforeDirectives(ctx *context.Context, f *Field, store *NodeStore, parent Node, result interface{}) errors.GraphqlErrorInterface {
	for _, directive := range f.DefinitionDirectives {
		if fn, ok := fieldRuntimeBeforeDirectiveMap[directive.Name]; ok {
			if err := fn(ctx, f, directive, store, parent, result); err != nil {
				return err
			}
		}
	}
	return nil
}

// result is data[field.Name] , directive can read it
func ExecuteFieldAfterDirectives(ctx *context.Context, f *Field, store *NodeStore, parent Node, result interface{}) errors.GraphqlErrorInterface {
	for _, directive := range f.DefinitionDirectives {
		if fn, ok := fieldRuntimeAfterDirectiveMap[directive.Name]; ok {
			if err := fn(ctx, f, directive, store, parent, result); err != nil {
				return err
			}
		}
	}
	return nil
}
