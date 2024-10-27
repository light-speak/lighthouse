package model

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

type SelectRelation struct {
	Relation      *ast.Relation
	SelectColumns map[string]interface{}
}

var quickListMap = make(map[string]func(ctx context.Context, columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error))

var quickFirstMap = make(map[string]func(ctx context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error))

var quickLoaderMap = make(map[string]func(ctx context.Context, key int64, field string) (map[string]interface{}, error))

func AddQuickList(name string, fn func(ctx context.Context, columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error)) {
	quickListMap[name] = fn
}
func AddQuickFirst(name string, fn func(ctx context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error)) {
	quickFirstMap[name] = fn
}
func AddQuickLoader(name string, fn func(ctx context.Context, key int64, field string) (map[string]interface{}, error)) {
	quickLoaderMap[name] = fn
}

func GetQuickFirst(name string) func(ctx context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
	fn, ok := quickFirstMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickList(name string) func(ctx context.Context, columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
	fn, ok := quickListMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickLoader(name string) func(ctx context.Context, key int64, field string) (map[string]interface{}, error) {
	fn, ok := quickLoaderMap[name]
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

func MapToStruct[T ModelInterface](data map[string]interface{}) (T, error) {
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

func FetchRelation(ctx context.Context, data map[string]interface{}, relation *SelectRelation) (interface{}, error) {
	switch relation.Relation.RelationType {
	case ast.RelationTypeBelongsTo:
		fieldValue, ok := data[relation.Relation.ForeignKey]
		if !ok {
			return nil, fmt.Errorf("field %s not found in function %s", relation.Relation.ForeignKey, "FetchRelation")
		}
		data, err := fetchBelongsTo(ctx, relation, fieldValue)
		if err != nil {
			return nil, err
		}
		data, err = GetQuickFirst(utils.UcFirst(relation.Relation.Name))(ctx, relation.SelectColumns, data)
		if err != nil {
			return nil, err
		}
		return data, nil
	case ast.RelationTypeHasMany:
		fieldValue, ok := data[relation.Relation.Reference]
		if !ok {
			return nil, fmt.Errorf("field %s not found in function %s", relation.Relation.Reference, "FetchRelation")
		}
		data, err := fetchHasMany(ctx, relation, fieldValue)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	return nil, nil
}

func fetchHasMany(ctx context.Context, relation *SelectRelation, fieldValue interface{}) ([]map[string]interface{}, error) {
	key, ok := fieldValue.(int64)
	if !ok {
		return nil, fmt.Errorf("relation %s field %s value is not int64", relation.Relation.Name, relation.Relation.ForeignKey)
	}
	data, err := GetQuickList(utils.UcFirst(relation.Relation.Name))(ctx, relation.SelectColumns, func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", relation.Relation.ForeignKey), key)
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

func fetchBelongsTo(ctx context.Context, relation *SelectRelation, fieldValue interface{}) (map[string]interface{}, error) {
	key, ok := fieldValue.(int64)
	if !ok {
		return nil, fmt.Errorf("relation %s field %s value is not int64", relation.Relation.Name, relation.Relation.ForeignKey)
	}
	data, err := GetQuickLoader(utils.UcFirst(relation.Relation.Name))(ctx, key, relation.Relation.Reference)
	if err != nil {
		return nil, err
	}
	return data, nil
}
