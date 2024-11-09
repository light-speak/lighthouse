package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/bytedance/sonic"
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
	scopesMap        sync.Map
)

func AddScopes(name string, fn func(ctx *context.Context) func(db *gorm.DB) *gorm.DB) {
	scopesMap.Store(name, fn)
}

func GetScopes(name string) func(ctx *context.Context) func(db *gorm.DB) *gorm.DB {
	if fn, ok := scopesMap.Load(name); ok {
		return fn.(func(ctx *context.Context) func(db *gorm.DB) *gorm.DB)
	}
	return nil
}

func AddQuickList(name string, fn func(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error)) {
	quickListMap.Store(name, fn)
}

func AddQuickFirst(name string, fn func(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error)) {
	quickFirstMap.Store(name, fn)
}

func AddQuickLoad(name string, fn func(ctx *context.Context, key int64, field string, filters ...*Filter) (*sync.Map, error)) {
	quickLoadMap.Store(name, fn)
}

func AddQuickLoadList(name string, fn func(ctx *context.Context, key int64, field string, filters ...*Filter) ([]*sync.Map, error)) {
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

func GetQuickLoad(name string) func(ctx *context.Context, key int64, field string, filters ...*Filter) (*sync.Map, error) {
	if fn, ok := quickLoadMap.Load(name); ok {
		return fn.(func(ctx *context.Context, key int64, field string, filters ...*Filter) (*sync.Map, error))
	}
	return nil
}

func GetQuickLoadList(name string) func(ctx *context.Context, key int64, field string, filters ...*Filter) ([]*sync.Map, error) {
	if fn, ok := quickLoadListMap.Load(name); ok {
		return fn.(func(ctx *context.Context, key int64, field string, filters ...*Filter) ([]*sync.Map, error))
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
	jsonData, err := sonic.Marshal(m)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err = sonic.Unmarshal(jsonData, &result); err != nil {
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

	jsonData, err := sonic.Marshal(tempMap)
	if err != nil {
		return m, err
	}
	if err = sonic.Unmarshal(jsonData, &m); err != nil {
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
		return fetchSingleRelation(ctx, relation.Name, relation.Reference, fieldValue)

	case ast.RelationTypeHasMany:
		fieldValue, ok := data.Load(relation.Reference)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Reference, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		return fetchMultipleRelations(ctx, relation.Name, relation.ForeignKey, fieldValue)

	case ast.RelationTypeHasOne:
		fieldValue, ok := data.Load(relation.Reference)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Reference, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		return fetchSingleRelation(ctx, relation.Name, relation.ForeignKey, fieldValue)

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
		return fetchSingleRelation(ctx, relationNameStr, relation.Reference, fieldValue)

	case ast.RelationTypeMorphMany:
		fieldValue, ok := data.Load(relation.Reference)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.Reference, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		filters := []*Filter{
			{
				Field: relation.MorphType,
				Value: relation.CurrentType,
			},
		}
		return fetchMultipleRelations(ctx, relation.Name, relation.MorphKey, fieldValue, filters...)

	case ast.RelationTypeBelongsToMany:
		fieldValue, ok := data.Load(relation.ForeignKey)
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("field %s not found in function %s", relation.ForeignKey, "FetchRelation"),
				Locations: []*errors.GraphqlLocation{},
			}
		}
		return fetchManyToManyRelations(ctx, relation, fieldValue)
	}
	return nil, nil
}

func fetchMultipleRelations(ctx *context.Context, relationName string, foreignKey string, fieldValue interface{}, filters ...*Filter) ([]*sync.Map, error) {
	key, err := convertToInt64(relationName, foreignKey, fieldValue)
	if err != nil {
		return nil, err
	}

	loadListFn := GetQuickLoadList(utils.UcFirst(utils.CamelCase(relationName)))
	if loadListFn == nil {
		return nil, fmt.Errorf("relation %s not found", relationName)
	}

	datas, err := loadListFn(ctx, key, foreignKey, filters...)
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

func fetchSingleRelation(ctx *context.Context, relationName string, foreignKey string, fieldValue interface{}) (*sync.Map, error) {
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

func fetchManyToManyRelations(ctx *context.Context, relation *ast.Relation, fieldValue interface{}) ([]*sync.Map, error) {
	key, err := convertToInt64(relation.Name, relation.ForeignKey, fieldValue)
	if err != nil {
		return nil, err
	}

	loadListFn := GetQuickLoadList(utils.UcFirst(utils.CamelCase(relation.Pivot)))
	if loadListFn == nil {
		return nil, fmt.Errorf("pivot relation %s not found", relation.Pivot)
	}

	pivotDatas, err := loadListFn(ctx, key, relation.PivotForeignKey)
	if err != nil {
		return nil, err
	}

	var relatedIds []int64
	for _, pivotData := range pivotDatas {
		if relatedId, ok := pivotData.Load(relation.PivotReference); ok {
			if id, err := convertToInt64(relation.Name, relation.PivotReference, relatedId); err == nil {
				relatedIds = append(relatedIds, id)
			}
		}
	}
	var results []*sync.Map
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, relatedId := range relatedIds {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			loadFn := GetQuickLoad(utils.UcFirst(utils.CamelCase(relation.Name)))
			if loadFn == nil {
				return
			}
			data, err := loadFn(ctx, id, relation.RelationForeignKey)
			if err != nil {
				return
			}
			firstFn := GetQuickFirst(utils.UcFirst(utils.CamelCase(relation.Name)))
			if firstFn == nil {
				return
			}

			data, err = firstFn(ctx, data)
			if err != nil {
				return
			}
			data.Store("__typename", relation.Name)
			mu.Lock()
			results = append(results, data)
			mu.Unlock()
		}(relatedId)
	}

	wg.Wait()

	return results, nil
}
