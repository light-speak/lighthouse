package excute

import (
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"gorm.io/gorm"
)

var quickExcuteMap = map[string]func(ctx *context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, errors.GraphqlErrorInterface){
	"find":     executeFind,
	"first":    executeFirst,
	"paginate": executePaginate,
	// TODO: create, update, delete,
}

func QuickExecute(ctx *context.Context, field *ast.Field) (interface{}, bool, errors.GraphqlErrorInterface) {
	scopes := make([]func(db *gorm.DB) *gorm.DB, 0)
	for _, arg := range field.DefinitionArgs {
		if len(arg.Directives) > 0 {
			scope, err := executeFilter(ctx, arg, field.Args[arg.Name].Value)
			if err != nil {
				return nil, false, err
			}
			scopes = append(scopes, scope)
		}
	}
	for _, directive := range field.DefinitionDirectives {
		if fn, ok := quickExcuteMap[directive.Name]; ok {
			res, err := fn(ctx, field, scopes...)
			if err != nil {
				return nil, true, err
			}
			if err := field.Type.ValidateValue(res, true); err != nil {
				return nil, true, err
			}
			return res, true, nil
		}
	}

	return nil, false, nil
}

func getColumns(field *ast.Field) (map[string]interface{}, errors.GraphqlErrorInterface) {
	res := make(map[string]interface{})
	if len(field.Children) == 0 {
		return nil, nil
	}
	for _, child := range field.Children {
		if child.IsFragment || child.IsUnion {
			for _, c := range child.Children {
				column, err := getColumns(c)
				if err != nil {
					return nil, err
				}
				res[c.Name] = column
			}
		} else if child.Children != nil {
			cRes := make(map[string]interface{})
			for _, c := range child.Children {
				column, err := getColumns(c)
				if err != nil {
					return nil, err
				}
				cRes[c.Name] = column
			}
			res[child.Name] = cRes
		} else {
			res[child.Name] = nil
		}
	}
	return res, nil
}
