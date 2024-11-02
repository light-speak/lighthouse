package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/utils"
)

func mergeData(ctx *context.Context, field *ast.Field, datas map[string]interface{}) (interface{}, errors.GraphqlErrorInterface) {
	fieldName := utils.SnakeCase(field.Name)

	v, ok := datas[fieldName]
	if !ok {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	if v == nil {
		if field.Relation != nil {
			cData, err := model.FetchRelation(ctx, datas, &model.SelectRelation{Relation: field.Relation})
			if err != nil {
				return nil, &errors.GraphQLError{
					Message:   err.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
			v = cData
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
			merged, err := mergeData(ctx, listField, m)
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
			c, err := mergeData(ctx, child, vMap)
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
