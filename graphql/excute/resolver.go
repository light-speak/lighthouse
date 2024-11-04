package excute

import (
	"sync"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/resolve"
)

func executeResolver(ctx *context.Context, field *ast.Field) (interface{}, bool, errors.GraphqlErrorInterface) {
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
	realType := field.Type.GetRealType()
	if realType.TypeNode.GetKind() == ast.KindObject && realType.TypeNode.(*ast.ObjectNode).IsModel {
		if field.Type.IsList() {
			return processListResult(ctx, field, realType, r)
		}
		return processSingleResult(ctx, field, realType, r)
	}

	data := make(map[string]interface{})
	var wg sync.WaitGroup
	errChan := make(chan errors.GraphqlErrorInterface, len(field.Children))
	resultChan := make(chan struct {
		key   string
		value interface{}
	}, len(field.Children))

	for _, child := range field.Children {
		wg.Add(1)
		go func(c *ast.Field) {
			defer wg.Done()
			v, err := mergeData(ctx, c, r.(map[string]interface{}))
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- struct {
				key   string
				value interface{}
			}{c.Name, v}
		}(child)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(resultChan)
	}()

	if err := <-errChan; err != nil {
		return nil, true, err
	}

	for result := range resultChan {
		data[result.key] = result.value
	}

	return data, true, nil
}

func processListResult(ctx *context.Context, field *ast.Field, realType *ast.TypeRef, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	result, err := model.GetQuickList(realType.Name)(ctx, r.([]map[string]interface{}))
	if err != nil {
		return nil, true, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	data := make([]map[string]interface{}, len(result))
	var wg sync.WaitGroup
	errChan := make(chan errors.GraphqlErrorInterface, len(result))

	for i, ri := range result {
		wg.Add(1)
		go func(index int, item map[string]interface{}) {
			defer wg.Done()
			riData := make(map[string]interface{})
			for _, child := range field.Children {
				v, err := mergeData(ctx, child, item)
				if err != nil {
					errChan <- err
					return
				}
				riData[child.Name] = v
			}
			data[index] = riData
		}(i, ri)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, true, <-errChan
	}

	return data, true, nil
}

func processSingleResult(ctx *context.Context, field *ast.Field, realType *ast.TypeRef, r interface{}) (interface{}, bool, errors.GraphqlErrorInterface) {
	result, err := model.GetQuickFirst(realType.Name)(ctx, r.(map[string]interface{}))
	if err != nil {
		return nil, true, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	data := make(map[string]interface{})
	var wg sync.WaitGroup
	errChan := make(chan errors.GraphqlErrorInterface, len(field.Children))
	resultChan := make(chan struct {
		key   string
		value interface{}
	}, len(field.Children))

	for _, child := range field.Children {
		wg.Add(1)
		go func(c *ast.Field) {
			defer wg.Done()
			v, err := mergeData(ctx, c, result)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- struct {
				key   string
				value interface{}
			}{c.Name, v}
		}(child)
	}

	go func() {
		wg.Wait()
		close(errChan)
		close(resultChan)
	}()

	if err := <-errChan; err != nil {
		return nil, true, err
	}

	for result := range resultChan {
		data[result.key] = result.value
	}

	return data, true, nil
}
