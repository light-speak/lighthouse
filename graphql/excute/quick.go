package excute

import (
	"encoding/json"
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/log"
)

func QuickExecute(field *ast.Field) (interface{}, error) {
	paginate := ast.GetDirective("paginate", field.DefinitionDirectives)
	if paginate != nil {
		return executePaginate(field)
	}
	first := ast.GetDirective("first", field.DefinitionDirectives)
	if first != nil {
		return executeFirst(field)
	}
	create := ast.GetDirective("create", field.DefinitionDirectives)
	if create != nil {
		return nil, nil
	}
	update := ast.GetDirective("update", field.DefinitionDirectives)
	if update != nil {
		return nil, nil
	}
	delete := ast.GetDirective("delete", field.DefinitionDirectives)
	if delete != nil {
		return nil, nil
	}

	return nil, nil
}

var quickListMap = make(map[string]func() ([]map[string]interface{}, error))

var quickFirstMap = make(map[string]func() (map[string]interface{}, error))

func AddQuickList(name string, fn func() ([]map[string]interface{}, error)) {
	quickListMap[name] = fn
}

func AddQuickFirst(name string, fn func() (map[string]interface{}, error)) {
	quickFirstMap[name] = fn
}

func executeFirst(field *ast.Field) (interface{}, error) {
	fn, ok := quickFirstMap[field.Name]
	if !ok {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}

	column, err := getColumn(field)
	if err != nil {
		return nil, err
	}
	columnJson, err := json.Marshal(column)
	if err != nil {
		return nil, err
	}
	log.Error().Msgf("%s", columnJson)

	d, err := fn()
	if err != nil {
		return nil, err
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

func executePaginate(field *ast.Field) (interface{}, error) {
	res := make(map[string]interface{})
	info := &model.PaginateInfo{}
	res["paginateInfo"] = info
	datas, err := quickListMap[field.Name]()
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

func mergeData(field *ast.Field, datas map[string]interface{}) (interface{}, error) {
	v, ok := datas[field.Name]
	if !ok {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}
	if field.Children != nil {
		cv := make(map[string]interface{})
		for _, child := range field.Children {
			c, err := mergeData(child, v.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			cv[child.Name] = c
		}
		return cv, nil
	}
	return v, nil
}

func getColumn(field *ast.Field) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if len(field.Children) == 0 {
		return nil, nil
	}
	for _, child := range field.Children {
		if child.IsFragment || child.IsUnion {
			for _, c := range child.Children {
				column, err := getColumn(c)
				if err != nil {
					return nil, err
				}
				res[c.Name] = column
			}
		} else if child.Children != nil {
			cRes := make(map[string]interface{})
			for _, c := range child.Children {
				column, err := getColumn(c)
				if err != nil {
					return nil, err
				}
				cRes[c.Name] = column
			}
			res[child.Name] = cRes
		} else {
			res[child.Name] = nil
		}
	}
	return res, nil
}
