package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
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
func mergeData(field *ast.Field, datas map[string]interface{}) (interface{}, errors.GraphqlErrorInterface) {
	fieldName := utils.SnakeCase(field.Name)
	v, ok := datas[fieldName]
	if !ok {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	if v == nil {
		return nil, nil
	}

	typeRef := field.Type
	if typeRef.Kind == ast.KindNonNull {
		typeRef = typeRef.OfType
	}

	if typeRef.Kind == ast.KindList {
		vList, ok := v.([]map[string]interface{})
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("expected list but got %T", v),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		result := make([]interface{}, len(vList))
		for i, item := range vList {
			listField := &ast.Field{
				Name:     field.Name,
				Children: field.Children,
				Type:     typeRef.OfType,
			}

			m := make(map[string]interface{})
			m[field.Name] = item
			merged, err := mergeData(listField, m)
			if err != nil {
				return nil, err
			}
			result[i] = merged
		}
		return result, nil
	}

	if field.Children != nil {
		cv := make(map[string]interface{})
		vMap := v.(map[string]interface{})
		for _, child := range field.Children {
			c, err := mergeData(child, vMap)
			if err != nil {
				return nil, err
			}
			cv[child.Name] = c
		}
		return cv, nil
	}

	v, err := ValidateValue(field, v, false)
	if err != nil {
		return nil, err
	}
	return v, nil
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
