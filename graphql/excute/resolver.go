package excute

import (
	"fmt"
	"sync"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/resolve"
	"github.com/light-speak/lighthouse/utils"
)

func executeResolver(ctx *context.Context, field *ast.Field) (interface{}, bool, errors.GraphqlErrorInterface) {
	defer func() {
		if r := recover(); r != nil {
			ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
				Message:   fmt.Sprintf("panic: %v", r),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			})
		}
	}()
	if resolverFunc, ok := resolverMap[field.Name]; ok {
		args := make(map[string]any)
		for _, arg := range field.Args {
			args[arg.Name] = arg.Value
		}
		r, e := resolverFunc(ctx, args, resolve.R)
		if e != nil {
			return nil, true, &errors.GraphQLError{
				Message:   e.Error(),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		if r == nil {
			if field.Type.Kind == ast.KindNonNull {
				return nil, true, &errors.GraphQLError{
					Message:   fmt.Sprintf("field %s is not nullable", field.Name),
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
	realType := field.Type.GetRealType()
	if realType.TypeNode.GetKind() == ast.KindObject && realType.TypeNode.(*ast.ObjectNode).IsModel {
		if field.Type.IsList() {
			return processListResult(ctx, field, realType, r)
		}
		return processSingleResult(ctx, field, realType, r)
	}

	dataMap := r.(*sync.Map)

	totalFields := countTotalFields(field.Children)
	var wg sync.WaitGroup
	errChan := make(chan errors.GraphqlErrorInterface, totalFields)
	resultChan := make(chan struct {
		key   string
		value interface{}
	}, totalFields)

	processFields(ctx, field.Children, dataMap, &wg, errChan, resultChan)

	go func() {
		wg.Wait()
		close(errChan)
		close(resultChan)
	}()

	if err := <-errChan; err != nil {
		return nil, true, err
	}

	data := make(map[string]interface{})
	for result := range resultChan {
		data[result.key] = result.value
	}

	return data, true, nil
}

func processListResult(ctx *context.Context, field *ast.Field, realType *ast.TypeRef, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	rSlice, ok := r.([]*sync.Map)
	if !ok {
		return nil, true, &errors.GraphQLError{
			Message:   "invalid input type for list result",
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	result, err := model.GetQuickList(realType.Name)(ctx, rSlice)
	if err != nil {
		return nil, true, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	data := make([]map[string]interface{}, len(result))
	var wg sync.WaitGroup
	var lastErr errors.GraphqlErrorInterface
	resultMapSlice := make([]*sync.Map, len(result))

	// Process each item concurrently
	for i, item := range result {
		wg.Add(1)
		go func(i int, item *sync.Map) {
			defer wg.Done()
			resultMap := &sync.Map{}

			for _, child := range field.Children {
				if child.Name == "__typename" {
					resultMap.Store(child.Name, realType.Name)
					continue
				}

				value := getValueFromSyncMap(item, child.Name)
				if value == nil {
					var err errors.GraphqlErrorInterface
					value, err = mergeData(ctx, child, item)
					if err != nil {
						lastErr = err
						return
					}
				}
				resultMap.Store(child.Name, value)
			}
			resultMapSlice[i] = resultMap
		}(i, item)
	}
	wg.Wait()

	if lastErr != nil {
		return nil, true, lastErr
	}

	// Convert sync.Map to regular map
	for i, resultMap := range resultMapSlice {
		m := make(map[string]interface{})
		resultMap.Range(func(key, value interface{}) bool {
			m[key.(string)] = value
			return true
		})
		data[i] = m
	}

	return data, true, nil
}

func processSingleResult(ctx *context.Context, field *ast.Field, realType *ast.TypeRef, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	syncMap, ok := r.(*sync.Map)
	if !ok {
		return nil, true, &errors.GraphQLError{
			Message:   "invalid input type for single result",
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	result, err := model.GetQuickFirst(realType.Name)(ctx, syncMap)
	if err != nil {
		return nil, true, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	resultMap := &sync.Map{}
	var wg sync.WaitGroup
	var lastErr errors.GraphqlErrorInterface

	// Process each field concurrently
	for _, child := range field.Children {
		wg.Add(1)
		go func(child *ast.Field) {
			defer wg.Done()
			if child.Name == "__typename" {
				resultMap.Store(child.Name, realType.Name)
				return
			}

			value := getValueFromSyncMap(result, child.Name)
			if value == nil {
				var err errors.GraphqlErrorInterface
				value, err = mergeData(ctx, child, result)
				if err != nil {
					lastErr = err
					return
				}
			}
			resultMap.Store(child.Name, value)
		}(child)
	}
	wg.Wait()

	if lastErr != nil {
		return nil, true, lastErr
	}

	// Convert sync.Map to regular map
	data := make(map[string]interface{})
	resultMap.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	return data, true, nil
}

// Helper function to get value from sync.Map using snake case key
func getValueFromSyncMap(m *sync.Map, key string) interface{} {
	var value interface{}
	snakeName := utils.SnakeCase(key)
	m.Range(func(k, v interface{}) bool {
		if k.(string) == snakeName {
			value = v
			return false
		}
		return true
	})
	return value
}
