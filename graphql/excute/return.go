package excute

import (
	"fmt"
	"math"
	"sync"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"gorm.io/gorm"
)

func executeFirst(ctx *context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, errors.GraphqlErrorInterface) {
	fn := model.GetQuickFirst(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found in function %s", field.Type.GetGoName(), "executeFirst"),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	d, e := fn(ctx, nil, scopes...)
	if e != nil {
		return nil, &errors.GraphQLError{
			Message:   e.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	if d == nil {
		return nil, nil
	}
	data := make(map[string]interface{})
	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr errors.GraphqlErrorInterface

	for _, child := range field.Children {
		wg.Add(1)
		go func(child *ast.Field) {
			defer wg.Done()
			v, err := mergeData(ctx, child, d)
			if err != nil {
				mu.Lock()
				lastErr = err
				mu.Unlock()
				return
			}
			mu.Lock()
			data[child.Name] = v
			mu.Unlock()
		}(child)
	}
	wg.Wait()

	if lastErr != nil {
		return nil, lastErr
	}
	return data, nil
}

func executePaginate(ctx *context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, errors.GraphqlErrorInterface) {
	pageArg := field.Args["page"]
	sizeArg := field.Args["size"]
	sortArg := field.Args["sort"]
	page, err := pageArg.GetValue()
	if err != nil {
		return nil, err
	}
	size, err := sizeArg.GetValue()
	if err != nil {
		return nil, err
	}
	sort, err := sortArg.GetValue()
	if err != nil {
		return nil, err
	}
	scope := func(db *gorm.DB) *gorm.DB {
		return db.Offset((int(page.(int64)) - 1) * int(size.(int64))).Limit(int(size.(int64))).Order(fmt.Sprintf("%s %s", "id", sort.(string)))
	}

	var wg sync.WaitGroup
	var data interface{}
	var dataErr errors.GraphqlErrorInterface
	var count int64
	var countErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		data, dataErr = executeFind(ctx, field.Children["data"], append(scopes, scope)...)
	}()

	if field.Children["paginateInfo"] != nil {
		if field.Children["paginateInfo"].Children["totalPage"] != nil ||
			field.Children["paginateInfo"].Children["hasNextPage"] != nil ||
			field.Children["paginateInfo"].Children["totalCount"] != nil {

			wg.Add(1)
			go func() {
				defer wg.Done()
				countFn := model.GetQuickCount(field.Children["data"].Type.GetRealType().GetGoName())
				if countFn == nil {
					dataErr = &errors.GraphQLError{
						Message:   fmt.Sprintf("quick count function %s not found", field.Type.GetGoName()),
						Locations: []*errors.GraphqlLocation{field.GetLocation()},
					}
					return
				}
				count, countErr = countFn(scopes...)
			}()
		}
	}

	wg.Wait()

	if dataErr != nil {
		return nil, dataErr
	}
	if countErr != nil {
		return nil, &errors.GraphQLError{
			Message:   countErr.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	res := make(map[string]interface{})
	res["data"] = data

	if field.Children["paginateInfo"] != nil {
		paginateInfo := make(map[string]interface{})
		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, child := range field.Children["paginateInfo"].Children {
			wg.Add(1)
			go func(child *ast.Field) {
				defer wg.Done()
				value := mergePaginateInfo(child, count, page.(int64), size.(int64))
				mu.Lock()
				paginateInfo[child.Name] = value
				mu.Unlock()
			}(child)
		}
		wg.Wait()
		res["paginateInfo"] = paginateInfo
	}
	return res, nil
}

func mergePaginateInfo(field *ast.Field, count int64, page int64, size int64) interface{} {
	switch field.Name {
	case "totalCount":
		return count
	case "currentPage":
		return page
	case "hasNextPage":
		totalPage := int64(math.Ceil(float64(count) / float64(size)))
		return page < totalPage
	case "totalPage":
		return int64(math.Ceil(float64(count) / float64(size)))
	}
	return nil
}

func executeFind(ctx *context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, errors.GraphqlErrorInterface) {
	fn := model.GetQuickList(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found in function %s", field.Name, "executeFind"),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	datas, e := fn(ctx, nil, scopes...)
	if e != nil {
		return nil, &errors.GraphQLError{
			Message:   e.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}

	data := make([]interface{}, len(datas))
	var wg sync.WaitGroup
	var mu sync.Mutex
	var lastErr errors.GraphqlErrorInterface

	for i, item := range datas {
		wg.Add(1)
		go func(i int, item interface{}) {
			defer wg.Done()
			d := make(map[string]interface{})
			for _, child := range field.Children {
				v, err := mergeData(ctx, child, item.(map[string]interface{}))
				if err != nil {
					mu.Lock()
					lastErr = err
					mu.Unlock()
					return
				}
				d[child.Name] = v
			}
			mu.Lock()
			data[i] = d
			mu.Unlock()
		}(i, item)
	}
	wg.Wait()

	if lastErr != nil {
		return nil, lastErr
	}
	return data, nil
}
