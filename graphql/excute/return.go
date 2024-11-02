package excute

import (
	"fmt"
	"math"

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
			Message:   fmt.Sprintf("field %s not found", field.Type.GetGoName()),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	columns, err := getColumns(field)
	if err != nil {
		return nil, err
	}
	d, e := fn(ctx, columns, nil, scopes...)
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
	for _, child := range field.Children {
		v, err := mergeData(ctx, child, d)
		data[child.Name] = v
		if err != nil {
			return nil, err
		}
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
	data, err := executeFind(ctx, field.Children["data"], append(scopes, scope)...)
	if err != nil {
		return nil, err
	}
	res := make(map[string]interface{})
	res["data"] = data
	if field.Children["paginateInfo"] != nil {
		count := int64(0)
		var e error
		if field.Children["paginateInfo"].Children["totalPage"] != nil ||
			field.Children["paginateInfo"].Children["hasNextPage"] != nil ||
			field.Children["paginateInfo"].Children["totalCount"] != nil {
			countFn := model.GetQuickCount(field.Children["data"].Type.GetRealType().GetGoName())
			if countFn == nil {
				return nil, &errors.GraphQLError{
					Message:   fmt.Sprintf("quick count function %s not found", field.Type.GetGoName()),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
			count, e = countFn(scopes...)
			if e != nil {
				return nil, &errors.GraphQLError{
					Message:   e.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
		}
		paginateInfo := make(map[string]interface{})
		for _, child := range field.Children["paginateInfo"].Children {
			paginateInfo[child.Name] = mergePaginateInfo(child, count, page.(int64), size.(int64))
		}
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
	columns, err := getColumns(field)
	if err != nil {
		return nil, err
	}
	fn := model.GetQuickList(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	datas, e := fn(ctx, columns, nil, scopes...)
	if e != nil {
		return nil, &errors.GraphQLError{
			Message:   e.Error(),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
	}
	data := make([]interface{}, 0)
	for _, item := range datas {
		d := make(map[string]interface{})
		for _, child := range field.Children {
			v, err := mergeData(ctx, child, item)
			d[child.Name] = v
			if err != nil {
				return nil, err
			}
		}
		data = append(data, d)
	}
	return data, nil
}
