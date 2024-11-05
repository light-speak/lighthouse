package excute

import (
	"fmt"
	"sync"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/utils"
)

// countTotalFields recursively counts total fields including nested fragments
func countTotalFields(fields map[string]*ast.Field) int {
	total := 0
	for _, field := range fields {
		if field.IsFragment || field.IsUnion {
			total += countTotalFields(field.Children)
		} else {
			total++
		}
	}
	return total
}

// processFields handles field processing and sends results to channels
func processFields(
	ctx *context.Context,
	fields map[string]*ast.Field,
	data *sync.Map,
	wg *sync.WaitGroup,
	errChan chan<- errors.GraphqlErrorInterface,
	resultChan chan<- struct {
		key   string
		value interface{}
	},
) {
	for _, field := range fields {
		if field.IsFragment || field.IsUnion {
			// For union types, check the actual type before processing fields
			if field.IsUnion {
				typename, ok := data.Load("__typename")
				if !ok {
					continue
				}
				typeName, ok := typename.(string)
				if !ok {
					continue
				}
				// Only process fields if they match the actual type
				if field.Type.GetRealType().Name != utils.UcFirst(typeName) {
					continue
				}
			}
			processFields(ctx, field.Children, data, wg, errChan, resultChan)
		} else {
			wg.Add(1)
			go func(f *ast.Field) {
				defer wg.Done()
				c, err := mergeData(ctx, f, data)
				if err != nil {
					errChan <- err
					return
				}
				resultChan <- struct {
					key   string
					value interface{}
				}{f.Name, c}
			}(field)
		}
	}
}

// mergeData merges the field data with the given data map based on GraphQL field definition
func mergeData(ctx *context.Context, field *ast.Field, datas *sync.Map) (interface{}, errors.GraphqlErrorInterface) {
	defer func() {
		if r := recover(); r != nil {
			ctx.Errors = append(ctx.Errors, &errors.GraphQLError{
				Message:   fmt.Sprintf("mergeData panic: %v", r),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			})
		}
	}()

	fieldName := utils.SnakeCase(field.Name)
	var v interface{}

	// Get value from sync.Map
	if ev, exists := datas.Load(fieldName); exists && ev != nil {
		v = ev
	}

	// Handle relation fields
	if v == nil && field.Relation != nil {
		cData, err := model.FetchRelation(ctx, datas, field.Relation)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("Failed to fetch relation: %v", err),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}
		v = cData
	}

	if v == nil {
		return nil, nil
	}

	typeRef := field.Type
	if typeRef.Kind == ast.KindNonNull {
		typeRef = typeRef.OfType
	}

	// Handle list type
	if typeRef.Kind == ast.KindList {
		vList, ok := v.([]*sync.Map)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("Expected list type but got %T", v),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		result := make([]interface{}, len(vList))
		errChan := make(chan error, len(vList))
		var wg sync.WaitGroup

		for i, item := range vList {
			wg.Add(1)
			go func(index int, itemData *sync.Map) {
				defer wg.Done()

				if itemData == nil {
					errChan <- fmt.Errorf("nil item data at index %d", index)
					return
				}

				if field.Children != nil {
					totalFields := countTotalFields(field.Children)
					childErrChan := make(chan errors.GraphqlErrorInterface, totalFields)
					childResultChan := make(chan struct {
						key   string
						value interface{}
					}, totalFields)
					var childWg sync.WaitGroup
					processFields(ctx, field.Children, itemData, &childWg, childErrChan, childResultChan)

					go func() {
						childWg.Wait()
						close(childErrChan)
						close(childResultChan)
					}()

					if err := <-childErrChan; err != nil {
						errChan <- err
						return
					}

					merged := make(map[string]interface{})
					for result := range childResultChan {
						merged[result.key] = result.value
					}

					result[index] = merged
				}
			}(i, item)
		}

		wg.Wait()
		close(errChan)

		if err := <-errChan; err != nil {
			return nil, &errors.GraphQLError{
				Message:   err.Error(),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		return result, nil
	}

	// Handle object type with children
	if field.Children != nil {
		vMap, ok := v.(map[string]interface{})
		if !ok {
			vSyncMap, ok := v.(*sync.Map)
			if !ok {
				return nil, &errors.GraphQLError{
					Message:   fmt.Sprintf("Expected map or sync.Map type but got %T", v),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
			vMap = make(map[string]interface{})
			vSyncMap.Range(func(key, value interface{}) bool {
				vMap[key.(string)] = value
				return true
			})
		}

		sMap := &sync.Map{}
		for k, v := range vMap {
			sMap.Store(k, v)
		}

		totalFields := countTotalFields(field.Children)
		var wg sync.WaitGroup
		errChan := make(chan errors.GraphqlErrorInterface, totalFields)
		resultChan := make(chan struct {
			key   string
			value interface{}
		}, totalFields)

		processFields(ctx, field.Children, sMap, &wg, errChan, resultChan)

		go func() {
			wg.Wait()
			close(errChan)
			close(resultChan)
		}()

		if err := <-errChan; err != nil {
			return nil, err
		}

		finalResult := make(map[string]interface{})

		for result := range resultChan {
			finalResult[result.key] = result.value
		}

		return finalResult, nil
	}

	// Validate scalar value
	v, err := ValidateValue(field, v, false)
	if err != nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("Value validation failed: %v", err),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	return v, nil
}
