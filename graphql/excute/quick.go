package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/log"
)

func QuickExecute(field *ast.Field) (interface{}, bool, error) {
	paginate := ast.GetDirective("paginate", field.DefinitionDirectives)
	if paginate != nil {
		res, err := executePaginate(field)
		return res, true, err
	}
	first := ast.GetDirective("first", field.DefinitionDirectives)
	if first != nil {
		res, err := executeFirst(field)
		return res, true, err
	}
	create := ast.GetDirective("create", field.DefinitionDirectives)
	if create != nil {
		return nil, false, nil
	}
	update := ast.GetDirective("update", field.DefinitionDirectives)
	if update != nil {
		return nil, false, nil
	}
	delete := ast.GetDirective("delete", field.DefinitionDirectives)
	if delete != nil {
		return nil, false, nil
	}

	return nil, false, nil
}

func executeFirst(field *ast.Field) (interface{}, error) {
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
	d, err := fn(columns)
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

func executePaginate(field *ast.Field) (interface{}, error) {
	res := make(map[string]interface{})
	info := &model.PaginateInfo{}
	res["paginateInfo"] = info
	fn := model.GetQuickList(field.Type.GetGoName())
	if fn == nil {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s not found", field.Name),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}
	datas, err := fn()
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
	v, err := ValidateValue(field, v, false)
	log.Error().Msgf("mergeData %v", v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getColumns(field *ast.Field) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	if len(field.Children) == 0 {
		return nil, nil
	}
	for _, child := range field.Children {
		if child.IsFragment || child.IsUnion {
			for _, c := range child.Children {
				column, err := getColumns(c)
				if err != nil {
					return nil, err
				}
				res[c.Name] = column
			}
		} else if child.Children != nil {
			cRes := make(map[string]interface{})
			for _, c := range child.Children {
				column, err := getColumns(c)
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
