package excute

import (
	"context"
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"gorm.io/gorm"
)

func executeFirst(ctx context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, error) {
	fn := model.GetQuickFirst(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Type.GetGoName()),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}
	columns, err := getColumns(field)
	if err != nil {
		return nil, err
	}
	d, err := fn(ctx, columns, nil, scopes...)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, nil
	}
	data := make(map[string]interface{})
	for _, child := range field.Children {
		v, err := mergeData(child, d)
		data[child.Name] = v
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func executePaginate(ctx context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, error) {
	res := make(map[string]interface{})
	info := &model.PaginateInfo{}
	res["paginateInfo"] = info
	columns, err := getColumns(field)
	if err != nil {
		return nil, err
	}
	fn := model.GetQuickList(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}
	datas, err := fn(ctx, columns, scopes...)
	if err != nil {
		return nil, err
	}
	values := make([]map[string]interface{}, 0)
	for _, data := range datas {
		d := make(map[string]interface{}, 0)
		for _, child := range field.Children {
			v, err := mergeData(child, data)
			d[child.Name] = v
			if err != nil {
				return nil, err
			}
		}
		values = append(values, d)
	}
	res["data"] = values

	return res, nil
}

func executeFind(ctx context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, error) {
	columns, err := getColumns(field)
	if err != nil {
		return nil, err
	}
	fn := model.GetQuickList(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}
	datas, err := fn(ctx, columns, scopes...)
	if err != nil {
		return nil, err
	}
	data := make([]interface{}, 0)
	for _, item := range datas {
		d := make(map[string]interface{})
		for _, child := range field.Children {
			v, err := mergeData(child, item)
			d[child.Name] = v
			if err != nil {
				return nil, err
			}
		}
		data = append(data, d)
	}
	return data, nil
}
