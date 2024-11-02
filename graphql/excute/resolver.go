package excute

import (
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/log"
)

func executeResolver(ctx *context.Context, field *ast.Field) (interface{}, bool, errors.GraphqlErrorInterface) {
	if resolverFunc, ok := resolverMap[field.Name]; ok {
		args := make(map[string]any)
		for _, arg := range field.Args {
			args[arg.Name] = arg.Value
		}
		r, e := resolverFunc(ctx, args)
		if e != nil {
			return nil, true, &errors.GraphQLError{
				Message:   e.Error(),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		if r == nil {
			if field.Type.Kind == ast.KindNonNull {
				return nil, true, &errors.GraphQLError{
					Message:   "field is not nullable",
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
			return nil, true, nil
		}

		if field.Type.IsScalar() {
			return r, true, nil
		}

		return processObjectResult(ctx, field, r)
	}
	return nil, false, nil
}

func processObjectResult(ctx *context.Context, field *ast.Field, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	columns, e := getColumns(field)
	if e != nil {
		return nil, true, &errors.GraphQLError{
			Message:   e.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	realType := field.Type.GetRealType()
	if realType.TypeNode.GetKind() == ast.KindObject && realType.TypeNode.(*ast.ObjectNode).IsModel {
		if field.Type.IsList() {
			return processListResult(ctx, field, realType, columns, r)
		}
		return processSingleResult(ctx, field, realType, columns, r)
	}

	data := make(map[string]interface{})
	for _, child := range field.Children {
		v, err := mergeData(ctx, child, r.(map[string]interface{}))
		if err != nil {
			return nil, true, err
		}
		data[child.Name] = v
	}

	log.Info().Msgf("processObjectResult: %v", data)
	return data, true, nil
}

func processListResult(ctx *context.Context, field *ast.Field, realType *ast.TypeRef, columns map[string]interface{}, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	result, err := model.GetQuickList(realType.Name)(ctx, columns, r.([]map[string]interface{}))
	if err != nil {
		return nil, true, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	data := make([]map[string]interface{}, 0)
	for _, ri := range result {
		riData := make(map[string]interface{})
		for _, child := range field.Children {
			v, err := mergeData(ctx, child, ri)
			if err != nil {
				return nil, true, err
			}
			riData[child.Name] = v
		}
		data = append(data, riData)
	}
	return data, true, nil
}

func processSingleResult(ctx *context.Context, field *ast.Field, realType *ast.TypeRef, columns map[string]interface{}, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	result, err := model.GetQuickFirst(realType.Name)(ctx, columns, r.(map[string]interface{}))
	if err != nil {
		return nil, true, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	data := make(map[string]interface{})
	for _, child := range field.Children {
		v, err := mergeData(ctx, child, result)
		if err != nil {
			return nil, true, err
		}
		data[child.Name] = v
	}
	return data, true, nil
}
