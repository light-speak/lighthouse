package excute

import (
	"fmt"
	"math"
	"sync"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

// Execute first record query
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

	resultMap := &sync.Map{}
	var wg sync.WaitGroup
	var lastErr errors.GraphqlErrorInterface

	for _, child := range field.Children {
		wg.Add(1)
		go func(child *ast.Field) {
			defer wg.Done()
			if child.Name == "__typename" {
				resultMap.Store("__typename", utils.SnakeCase(field.Type.GetGoName()))
				return
			}
			v, err := mergeData(ctx, child, d)
			if err != nil {
				lastErr = err
				return
			}
			resultMap.Store(child.Name, v)
		}(child)
	}
	wg.Wait()

	if lastErr != nil {
		return nil, lastErr
	}

	data := make(map[string]interface{})
	resultMap.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	return data, nil
}

// Execute paginated query
func executePaginate(ctx *context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, errors.GraphqlErrorInterface) {
	page, err := field.Args["page"].GetValue()
	if err != nil {
		return nil, err
	}
	size, err := field.Args["size"].GetValue()
	if err != nil {
		return nil, err
	}
	sort, err := field.Args["sort"].GetValue()
	if err != nil {
		return nil, err
	}

	scope := func(db *gorm.DB) *gorm.DB {
		offset := (int(page.(int64)) - 1) * int(size.(int64))
		return db.Offset(offset).Limit(int(size.(int64))).Order(fmt.Sprintf("id %s", sort.(string)))
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
		needCount := false
		for _, child := range field.Children["paginateInfo"].Children {
			if child.Name == "totalPage" || child.Name == "hasNextPage" || child.Name == "totalCount" {
				needCount = true
				break
			}
		}

		if needCount {
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

	res := map[string]interface{}{
		"data": data,
	}

	if field.Children["paginateInfo"] != nil {
		paginateInfoMap := &sync.Map{}
		var wg sync.WaitGroup

		for _, child := range field.Children["paginateInfo"].Children {
			wg.Add(1)
			go func(child *ast.Field) {
				defer wg.Done()
				value := mergePaginateInfo(child, count, page.(int64), size.(int64))
				paginateInfoMap.Store(child.Name, value)
			}(child)
		}
		wg.Wait()

		paginateInfo := make(map[string]interface{})
		paginateInfoMap.Range(func(key, value interface{}) bool {
			paginateInfo[key.(string)] = value
			return true
		})
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
		return page < int64(math.Ceil(float64(count)/float64(size)))
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
	var lastErr errors.GraphqlErrorInterface

	for i, item := range datas {
		wg.Add(1)
		go func(i int, item *sync.Map) {
			defer wg.Done()

			resultMap := &sync.Map{}
			for _, child := range field.Children {
				if child.Name == "__typename" {
					resultMap.Store("__typename", utils.SnakeCase(field.Type.GetGoName()))
					continue
				}
				v, err := mergeData(ctx, child, item)
				if err != nil {
					lastErr = err
					return
				}
				resultMap.Store(child.Name, v)
			}

			d := make(map[string]interface{})
			resultMap.Range(func(key, value interface{}) bool {
				d[key.(string)] = value
				return true
			})

			data[i] = d
		}(i, item)
	}
	wg.Wait()

	if lastErr != nil {
		return nil, lastErr
	}
	return data, nil
}
