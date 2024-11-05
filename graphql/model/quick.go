package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

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

var (
	quickListMap     sync.Map
	quickFirstMap    sync.Map
	quickLoadMap     sync.Map
	quickLoadListMap sync.Map
	quickCountMap    sync.Map
)

func AddQuickList(name string, fn func(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error)) {
	quickListMap.Store(name, fn)
}

func AddQuickFirst(name string, fn func(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error)) {
	quickFirstMap.Store(name, fn)
}

func AddQuickLoad(name string, fn func(ctx *context.Context, key int64, field string) (*sync.Map, error)) {
	quickLoadMap.Store(name, fn)
}

func AddQuickLoadList(name string, fn func(ctx *context.Context, key int64, field string) ([]*sync.Map, error)) {
	quickLoadListMap.Store(name, fn)
}

func AddQuickCount(name string, fn func(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error)) {
	quickCountMap.Store(name, fn)
}

func GetQuickFirst(name string) func(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
	if fn, ok := quickFirstMap.Load(name); ok {
		return fn.(func(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error))
	}
	return nil
}

func GetQuickList(name string) func(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
	if fn, ok := quickListMap.Load(name); ok {
		return fn.(func(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error))
	}
	return nil
}

func GetQuickLoad(name string) func(ctx *context.Context, key int64, field string) (*sync.Map, error) {
	if fn, ok := quickLoadMap.Load(name); ok {
		return fn.(func(ctx *context.Context, key int64, field string) (*sync.Map, error))
	}
	return nil
}

func GetQuickLoadList(name string) func(ctx *context.Context, key int64, field string) ([]*sync.Map, error) {
	if fn, ok := quickLoadListMap.Load(name); ok {
		return fn.(func(ctx *context.Context, key int64, field string) ([]*sync.Map, error))
	}
	return nil
}

func GetQuickCount(name string) func(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
	if fn, ok := quickCountMap.Load(name); ok {
		return fn.(func(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error))
	}
	return nil
}

func StructToMap(m ModelInterface) (*sync.Map, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err = json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	syncMap := &sync.Map{}
	for key, value := range result {
		if v, ok := value.(float64); ok && v == float64(int64(v)) {
			syncMap.Store(key, int64(v))
		} else {
			syncMap.Store(key, value)
		}
	}

	syncMap.Store("__typename", m.TypeName())
	return syncMap, nil
}

func TypeToMap(m interface{}) (*sync.Map, error) {
	jsonData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err = json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	syncMap := &sync.Map{}
	for key, value := range result {
		if v, ok := value.(float64); ok && v == float64(int64(v)) {
			syncMap.Store(key, int64(v))
		} else {
			syncMap.Store(key, value)
		}
	}

	return syncMap, nil
}

func MapToStruct[T any](data *sync.Map) (T, error) {
	var m T
	tempMap := make(map[string]interface{})
	data.Range(func(key, value interface{}) bool {
		tempMap[key.(string)] = value
		return true
	})

	jsonData, err := json.Marshal(tempMap)
	if err != nil {
		return m, err
	}
	if err = json.Unmarshal(jsonData, &m); err != nil {
		return m, err
	}
	return m, nil
}

func FetchRelation(ctx *context.Context, data *sync.Map, relation *ast.Relation) (interface{}, error) {
	switch relation.RelationType {
	case ast.RelationTypeBelongsTo:
		fieldValue, ok := data.Load(relation.ForeignKey)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.ForeignKey, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		return fetchBelongsTo(ctx, relation.Name, relation.Reference, fieldValue)

	case ast.RelationTypeHasMany:
		fieldValue, ok := data.Load(relation.Reference)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Reference, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		return fetchHasMany(ctx, relation.Name, relation.ForeignKey, fieldValue)
	case ast.RelationTypeMorphTo:
		fieldValue, ok := data.Load(relation.MorphKey)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.MorphKey, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		relationName, ok := data.Load(relation.MorphType)
		if !ok {
			return nil, fmt.Errorf("field %s is not found", relation.MorphType)
		}
		relationNameStr, ok := relationName.(string)
		if !ok {
			return nil, fmt.Errorf("field %s is not a string", relation.MorphType)
		}
		return fetchBelongsTo(ctx, relationNameStr, relation.Reference, fieldValue)
	}
	return nil, nil
}

func fetchHasMany(ctx *context.Context, relationName string, foreignKey string, fieldValue interface{}) ([]*sync.Map, error) {
	key, err := convertToInt64(relationName, foreignKey, fieldValue)
	if err != nil {
		return nil, err
	}

	loadListFn := GetQuickLoadList(utils.UcFirst(utils.CamelCase(relationName)))
	if loadListFn == nil {
		return nil, fmt.Errorf("relation %s not found", relationName)
	}

	datas, err := loadListFn(ctx, key, foreignKey)
	if err != nil {
		return nil, err
	}

	listFn := GetQuickList(utils.UcFirst(utils.CamelCase(relationName)))
	if listFn == nil {
		return nil, fmt.Errorf("relation %s not found", relationName)
	}

	datas, err = listFn(ctx, datas)
	if err != nil {
		return nil, err
	}

	for _, data := range datas {
		data.Store("__typename", relationName)
	}
	return datas, nil
}

func fetchBelongsTo(ctx *context.Context, relationName string, foreignKey string, fieldValue interface{}) (*sync.Map, error) {
	key, err := convertToInt64(relationName, foreignKey, fieldValue)
	if err != nil {
		return nil, err
	}

	loadFn := GetQuickLoad(utils.UcFirst(utils.CamelCase(relationName)))
	if loadFn == nil {
		return nil, fmt.Errorf("relation %s not found", relationName)
	}

	data, err := loadFn(ctx, key, foreignKey)
	if err != nil {
		return nil, err
	}

	firstFn := GetQuickFirst(utils.UcFirst(utils.CamelCase(relationName)))
	if firstFn == nil {
		return nil, fmt.Errorf("relation %s not found", relationName)
	}

	data, err = firstFn(ctx, data)
	if err != nil {
		return nil, err
	}

	data.Store("__typename", relationName)
	return data, nil
}

func convertToInt64(relationName string, fieldName string, value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case string:
		key, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return key, nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("relation %s field %s value is not int64, got %v, type %T", relationName, fieldName, value, value)
	}
}
