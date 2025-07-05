package with

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

type Fields struct {
	Field  string
	Level  int
	Fields map[string]*Fields
}

func With(ctx context.Context) *Fields {
	fields := graphql.CollectFieldsCtx(ctx, nil)
	result := &Fields{
		Fields: make(map[string]*Fields),
		Level:  0,
	}

	for _, field := range fields {
		fields := &Fields{
			Field:  field.Name,
			Level:  1,
			Fields: make(map[string]*Fields),
		}
		result.Fields[field.Name] = fields

		if len(field.SelectionSet) > 0 {
			collectFields(field, fields, 2)
		}
	}

	if result.Field == "" && len(result.Fields) == 0 {
		return nil
	}

	return result
}

func collectFields(field graphql.CollectedField, parent *Fields, level int) {
	for _, selection := range field.SelectionSet {
		if field, ok := selection.(*ast.Field); ok {
			fields := &Fields{
				Field:  field.Name,
				Level:  level,
				Fields: make(map[string]*Fields),
			}
			parent.Fields[field.Name] = fields

			if len(field.SelectionSet) > 0 {
				collectFields(graphql.CollectedField{
					Field: field,
				}, fields, level+1)
			}
		}
	}
}
