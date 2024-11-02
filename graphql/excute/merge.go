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

// mergeData merges the field data with the given data map based on GraphQL field definition
func mergeData(ctx *context.Context, field *ast.Field, datas map[string]interface{}) (interface{}, errors.GraphqlErrorInterface) {
	// Convert field name to snake case for database column mapping

	fieldName := utils.SnakeCase(field.Name)
	var v interface{}

	// Get value from data map
	ev, exists := datas[fieldName]
	if exists && ev != nil {
		v = ev
	} else {
		v = nil
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

	// Return nil if value is nil
	if v == nil {
		return nil, nil
	}

	// Get real type by unwrapping NonNull
	typeRef := field.Type
	if typeRef.Kind == ast.KindNonNull {
		typeRef = typeRef.OfType
	}

	// Handle list type
	if typeRef.Kind == ast.KindList {
		vList, ok := v.([]map[string]interface{})
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("Expected list type but got %T", v),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		// Process each item in list
		result := make([]interface{}, len(vList))
		var wg sync.WaitGroup
		errChan := make(chan errors.GraphqlErrorInterface, len(vList))

		for i, item := range vList {
			wg.Add(1)
			go func(index int, itemData map[string]interface{}) {
				defer wg.Done()
				listField := &ast.Field{
					Name:     fieldName,
					Children: field.Children,
					Type:     typeRef.OfType,
				}

				m := make(map[string]interface{})
				m[fieldName] = itemData
				merged, err := mergeData(ctx, listField, m)
				if err != nil {
					errChan <- err
					return
				}
				result[index] = merged
			}(i, item)
		}

		wg.Wait()
		close(errChan)

		if len(errChan) > 0 {
			return nil, <-errChan
		}
		return result, nil
	}

	// Handle object type with children
	if field.Children != nil {
		cv := make(map[string]interface{})
		vMap, ok := v.(map[string]interface{})
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("Expected map type but got %T", v),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}

		// Process each child field
		var wg sync.WaitGroup
		errChan := make(chan errors.GraphqlErrorInterface, len(field.Children))
		resultChan := make(chan struct {
			key   string
			value interface{}
		}, len(field.Children))

		for _, child := range field.Children {

			wg.Add(1)
			go func(childField *ast.Field) {
				defer wg.Done()
				c, err := mergeData(ctx, childField, vMap)
				if err != nil {
					errChan <- err
					return
				}
				resultChan <- struct {
					key   string
					value interface{}
				}{childField.Name, c}
			}(child)
		}

		go func() {
			wg.Wait()
			close(errChan)
			close(resultChan)
		}()

		if err := <-errChan; err != nil {
			return nil, err
		}

		for result := range resultChan {
			cv[result.key] = result.value
		}

		return cv, nil
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
