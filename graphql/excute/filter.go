package excute

import (
	"context"
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"gorm.io/gorm"
)

func executeFilter(ctx context.Context, arg *ast.Argument, value interface{}) (func(db *gorm.DB) *gorm.DB, errors.GraphqlErrorInterface) {
	for _, directive := range arg.Directives {
		filter := filterMap[directive.Name]
		if filter == nil {
			continue
		}
		fieldName := arg.Name
		if fieldArg := directive.GetArg("field"); fieldArg != nil {
			fieldName = fieldArg.Value.(string)
		}
		return filter(fieldName, value), nil
	}
	return nil, nil
}

var filterMap = map[string]func(field string, value interface{}) func(db *gorm.DB) *gorm.DB{
	"eq": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", field), value)
		}
	},
	"neq": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s != ?", field), value)
		}
	},
	"gt": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s > ?", field), value)
		}
	},
	"gte": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s >= ?", field), value)
		}
	},
	"lt": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s < ?", field), value)
		}
	},
	"lte": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s <= ?", field), value)
		}
	},
	"in": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s IN (?)", field), value)
		}
	},
	"notIn": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s NOT IN (?)", field), value)
		}
	},
	"like": func(field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s LIKE ?", field), value)
		}
	},
}
