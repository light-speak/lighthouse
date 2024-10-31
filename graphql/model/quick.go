package model

import (
	"encoding/json"
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
	"github.com/light-speak/lighthouse/context"
	"gorm.io/gorm"
)

type SelectRelation struct {
	Relation      *ast.Relation
	SelectColumns map[string]interface{}
}

var quickListMap = make(map[string]func(ctx *context.Context, columns map[string]interface{}, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error))

var quickFirstMap = make(map[string]func(ctx *context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error))

var quickLoadMap = make(map[string]func(ctx *context.Context, key int64, field string) (map[string]interface{}, error))

var quickLoadListMap = make(map[string]func(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error))

func AddQuickList(name string, fn func(ctx *context.Context, columns map[string]interface{}, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error)) {
	quickListMap[name] = fn
}
func AddQuickFirst(name string, fn func(ctx *context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error)) {
	quickFirstMap[name] = fn
}
func AddQuickLoad(name string, fn func(ctx *context.Context, key int64, field string) (map[string]interface{}, error)) {
	quickLoadMap[name] = fn
}
func AddQuickLoadList(name string, fn func(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error)) {
	quickLoadListMap[name] = fn
}

func GetQuickFirst(name string) func(ctx *context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
	fn, ok := quickFirstMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickList(name string) func(ctx *context.Context, columns map[string]interface{}, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
	fn, ok := quickListMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickLoad(name string) func(ctx *context.Context, key int64, field string) (map[string]interface{}, error) {
	fn, ok := quickLoadMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickLoadList(name string) func(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error) {
	fn, ok := quickLoadListMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetSelectInfo(columns map[string]interface{}, provide map[string]*ast.Relation) (selectColumns []string, selectRelations map[string]*SelectRelation) {
	selectColumns = make([]string, 0)
	selectRelations = make(map[string]*SelectRelation, 0)

	for key, value := range columns {
		if value != nil && len(value.(map[string]interface{})) > 0 {
			relation := provide[key]
			selectRelations[key] = &SelectRelation{Relation: relation, SelectColumns: value.(map[string]interface{})}
			switch relation.RelationType {
			case ast.RelationTypeBelongsTo:
				selectColumns = append(selectColumns, relation.ForeignKey)
			case ast.RelationTypeHasMany:
				selectColumns = append(selectColumns, relation.Reference)
			}
		} else {
			selectColumns = append(selectColumns, key)
		}
	}
	return selectColumns, selectRelations
}

func StructToMap(m ModelInterface) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	for key, value := range result {
		switch v := value.(type) {
		case float64:
			if v == float64(int64(v)) {
				result[key] = int64(v)
			}
		}
	}

	result["__typename"] = m.TypeName()
	return result, nil
}

func MapToStruct[T any](data map[string]interface{}) (T, error) {
	var m T
	jsonData, err := json.Marshal(data)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(jsonData, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}

func FetchRelation(ctx *context.Context, data map[string]interface{}, relation *SelectRelation) (interface{}, error) {
	switch relation.Relation.RelationType {
	case ast.RelationTypeBelongsTo:
		fieldValue, ok := data[relation.Relation.ForeignKey]
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Relation.ForeignKey, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		data, err := fetchBelongsTo(ctx, relation, fieldValue)
		if err != nil {
			return nil, err
		}
		return data, nil
	case ast.RelationTypeHasMany:
		fieldValue, ok := data[relation.Relation.Reference]
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Relation.Reference, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		datas, err := fetchHasMany(ctx, relation, fieldValue)
		if err != nil {
			return nil, err
		}
		return datas, nil
	}
	return nil, nil
}

func fetchHasMany(ctx *context.Context, relation *SelectRelation, fieldValue interface{}) ([]map[string]interface{}, error) {
	key, ok := fieldValue.(int64)
	if !ok {
		return nil, fmt.Errorf("relation %s field %s value is not int64", relation.Relation.Name, relation.Relation.ForeignKey)
	}
	datas, err := GetQuickLoadList(utils.UcFirst(relation.Relation.Name))(ctx, key, relation.Relation.ForeignKey)
	if err != nil {
		return nil, err
	}
	datas, err = GetQuickList(utils.UcFirst(relation.Relation.Name))(ctx, relation.SelectColumns, datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func fetchBelongsTo(ctx *context.Context, relation *SelectRelation, fieldValue interface{}) (map[string]interface{}, error) {
	key, ok := fieldValue.(int64)
	if !ok {
		return nil, fmt.Errorf("relation %s field %s value is not int64", relation.Relation.Name, relation.Relation.ForeignKey)
	}
	data, err := GetQuickLoad(utils.UcFirst(relation.Relation.Name))(ctx, key, relation.Relation.Reference)
	if err != nil {
		return nil, err
	}
	data, err = GetQuickFirst(utils.UcFirst(relation.Relation.Name))(ctx, relation.SelectColumns, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
