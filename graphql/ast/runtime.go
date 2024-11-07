package ast

import (
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
)

func ExecuteFieldBeforeDirectives(ctx *context.Context, f *Field, store *NodeStore, parent Node) errors.GraphqlErrorInterface {
	for _, directive := range f.DefinitionDirectives {
		if fn, ok := fieldRuntimeBeforeDirectiveMap[directive.Name]; ok {
			if err := fn(ctx, f, directive, store, parent); err != nil {
				return err
			}
		}
	}
	return nil
}

func ExecuteFieldAfterDirectives(ctx *context.Context, f *Field, store *NodeStore, parent Node) errors.GraphqlErrorInterface {
	for _, directive := range f.DefinitionDirectives {
		if fn, ok := fieldRuntimeAfterDirectiveMap[directive.Name]; ok {
			if err := fn(ctx, f, directive, store, parent); err != nil {
				return err
			}
		}
	}
	return nil
}
