package excute

import (
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
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
		} else {
			columns, e := getColumns(field)
			if e != nil {
				return nil, true, &errors.GraphQLError{
					Message:   e.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
			realType := field.Type.GetRealType()
			if field.Type.IsList() {
				r, err := model.GetQuickList(realType.Name)(ctx, columns, r.([]map[string]interface{}))
				if err != nil {
					return nil, true, &errors.GraphQLError{
						Message:   err.Error(),
						Locations: []*errors.GraphqlLocation{field.GetLocation()},
					}
				}
				data := make([]map[string]interface{}, 0)
				for _, ri := range r {
					riData := make(map[string]interface{})
					for _, child := range field.Children {
						v, err := mergeData(child, ri)
						riData[child.Name] = v
						if err != nil {
							return nil, true, err
						}
					}
					data = append(data, riData)
				}
				return data, true, nil
			} else {
				r, err := model.GetQuickFirst(realType.Name)(ctx, columns, r.(map[string]interface{}))
				if err != nil {
					return nil, true, &errors.GraphQLError{
						Message:   err.Error(),
						Locations: []*errors.GraphqlLocation{field.GetLocation()},
					}
				}
				data := make(map[string]interface{})
				for _, child := range field.Children {
					v, err := mergeData(child, r)
					data[child.Name] = v
					if err != nil {
						return nil, true, err
					}
				}
				return data, true, nil
			}
		}

	}
	return nil, false, nil
}

// if resolverFunc, ok := resolverMap[field.Name]; ok {
// 	args := make(map[string]any)
// 	for _, arg := range field.Args {
// 		args[arg.Name] = arg.Value
// 	}
// 	r, e := resolverFunc(ctx, args)
// 	if e != nil {
// 		ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
// 			Message:   e.Error(),
// 			Locations: []*errors.GraphqlLocation{field.GetLocation()},
// 		})
// 		continue
// 	}

// 	if r == nil {
// 		if field.Type.Kind == ast.KindNonNull {
// 			ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
// 				Message:   "field is not nullable",
// 				Locations: []*errors.GraphqlLocation{field.GetLocation()},
// 			})
// 			continue
// 		}
// 		res[field.Name] = nil
// 		continue
// 	}

// 	if modelData, ok := r.(model.ModelInterface); ok {
// 		columns, e := getColumns(field)
// 		if e != nil {
// 			ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
// 				Message:   e.Error(),
// 				Locations: []*errors.GraphqlLocation{field.GetLocation()},
// 			})
// 			continue
// 		}
// 		modelMap, err := model.StructToMap(modelData)
// 		if err != nil {
// 			ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
// 				Message:   err.Error(),
// 				Locations: []*errors.GraphqlLocation{field.GetLocation()},
// 			})
// 			continue
// 		}
// 		isList := false

// 		returnType := field.Type
// 		if field.Type.Kind == ast.KindNonNull {
// 			returnType = field.Type.OfType
// 		}
// 		if returnType.Kind == ast.KindList {
// 			isList = true
// 			returnType = returnType.OfType
// 		}

// 		if returnType.Kind == ast.KindObject {
// 			if isList {
// 				model.GetQuickLoadList(returnType.Name)(ctx, modelMap)
// 			} else {
// 				model.GetQuickFirst(returnType.Name)(ctx, columns, modelMap)
// 			}
// 		}
// 		data := make(map[string]interface{})
// 		for _, child := range field.Children {
// 			d, err := mergeData(child, modelMap)
// 			if err != nil {
// 				ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
// 					Message:   err.Error(),
// 					Locations: []*errors.GraphqlLocation{child.GetLocation()},
// 				})
// 				continue
// 			}
// 			data[child.Name] = d
// 		}
// 		res[field.Name] = data
// 		continue
// 	}
// 	res[field.Name] = r
// 	continue
// }
