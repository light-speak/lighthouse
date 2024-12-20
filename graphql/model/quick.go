package model

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

type SelectRelation struct {
	Relation      *ast.Relation
	SelectColumns map[string]interface{}
}

var quickListMap = make(map[string]func(ctx *context.Context, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error))

var quickFirstMap = make(map[string]func(ctx *context.Context, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error))

var quickLoadMap = make(map[string]func(ctx *context.Context, key int64, field string) (map[string]interface{}, error))

var quickLoadListMap = make(map[string]func(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error))

var quickCountMap = make(map[string]func(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error))

func AddQuickList(name string, fn func(ctx *context.Context, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error)) {
	quickListMap[name] = fn
}
func AddQuickFirst(name string, fn func(ctx *context.Context, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error)) {
	quickFirstMap[name] = fn
}
func AddQuickLoad(name string, fn func(ctx *context.Context, key int64, field string) (map[string]interface{}, error)) {
	quickLoadMap[name] = fn
}
func AddQuickLoadList(name string, fn func(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error)) {
	quickLoadListMap[name] = fn
}
func AddQuickCount(name string, fn func(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error)) {
	quickCountMap[name] = fn
}

func GetQuickFirst(name string) func(ctx *context.Context, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
	fn, ok := quickFirstMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickList(name string) func(ctx *context.Context, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
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

func GetQuickCount(name string) func(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	fn, ok := quickCountMap[name]
	if !ok {
		return nil
	}
	return fn
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

func TypeToMap(m interface{}) (map[string]interface{}, error) {
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

func FetchRelation(ctx *context.Context, data map[string]interface{}, relation *ast.Relation) (interface{}, error) {
	switch relation.RelationType {
	case ast.RelationTypeBelongsTo:
		fieldValue, ok := data[relation.ForeignKey]
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.ForeignKey, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		data, err := fetchBelongsTo(ctx, relation, fieldValue)
		if err != nil {
			return nil, err
		}
		return data, nil
	case ast.RelationTypeHasMany:
		fieldValue, ok := data[relation.Reference]
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Reference, "FetchRelation"),
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

func fetchHasMany(ctx *context.Context, relation *ast.Relation, fieldValue interface{}) ([]map[string]interface{}, error) {
	var err error
	var key int64
	switch v := fieldValue.(type) {
	case int64:
		key = v
	case string:
		key, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	case float64:
		key = int64(v)
	default:
		return nil, fmt.Errorf("relation %s field %s value is not int64, got %v, type %T", relation.Name, relation.ForeignKey, fieldValue, fieldValue)
	}
	datas, err := GetQuickLoadList(utils.UcFirst(relation.Name))(ctx, key, relation.ForeignKey)
	if err != nil {
		return nil, err
	}
	datas, err = GetQuickList(utils.UcFirst(relation.Name))(ctx, datas)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func fetchBelongsTo(ctx *context.Context, relation *ast.Relation, fieldValue interface{}) (map[string]interface{}, error) {
	var err error
	var key int64
	switch v := fieldValue.(type) {
	case int64:
		key = v
	case string:
		key, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	case float64:
		key = int64(v)
	default:
		return nil, fmt.Errorf("relation %s field %s value is not int64, got %v, type %T", relation.Name, relation.ForeignKey, fieldValue, fieldValue)
	}
	data, err := GetQuickLoad(utils.UcFirst(relation.Name))(ctx, key, relation.Reference)
	if err != nil {
		return nil, err
	}
	data, err = GetQuickFirst(utils.UcFirst(relation.Name))(ctx, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
