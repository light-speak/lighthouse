package model

import (
	"encoding/json"
	"fmt"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

type SelectRelation struct {
	Relation      *ast.Relation
	selectColumns map[string]interface{}
}

var quickListMap = make(map[string]func() ([]map[string]interface{}, error))

var quickFirstMap = make(map[string]func(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error))

func AddQuickList(name string, fn func() ([]map[string]interface{}, error)) {
	quickListMap[name] = fn
}
func AddQuickFirst(name string, fn func(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error)) {
	quickFirstMap[name] = fn
}
func GetQuickFirst(name string) func(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
	fn, ok := quickFirstMap[name]
	if !ok {
		return nil
	}
	return fn
}

func GetQuickList(name string) func() ([]map[string]interface{}, error) {
	fn, ok := quickListMap[name]
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
			selectRelations[key] = &SelectRelation{Relation: relation, selectColumns: value.(map[string]interface{})}
			selectColumns = append(selectColumns, relation.ForeignKey)
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

func FetchRelation(res map[string]interface{}, relation *SelectRelation, fieldValue interface{}) (map[string]interface{}, error) {
	switch relation.Relation.RelationType {
	case ast.RelationTypeBelongsTo:
		data, err := fetchBelongsTo(relation, fieldValue)
		if err != nil {
			return nil, err
		}
		res[relation.Relation.Name] = data
	}
	return res, nil
}

func fetchBelongsTo(relation *SelectRelation, fieldValue interface{}) (map[string]interface{}, error) {
	data, err := GetQuickFirst(utils.UcFirst(relation.Relation.Name))(relation.selectColumns, func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s = ?", relation.Relation.Reference), fieldValue)
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}
