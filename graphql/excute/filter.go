package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/context"
	"gorm.io/gorm"

)

func executeFilter(ctx *context.Context, arg *ast.Argument, value interface{}) (func(db *gorm.DB) *gorm.DB, errors.GraphqlErrorInterface) {
	for _, directive := range arg.Directives {
		filter := filterMap[directive.Name]
		if filter == nil {
			continue
		}
		fieldName := arg.Name
		if fieldArg := directive.GetArg("field"); fieldArg != nil {
			fieldName = fieldArg.Value.(string)
		}
		return filter(ctx, fieldName, value), nil
	}
	return nil, nil
}

var filterMap = map[string]func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB{
	"eq": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s = ?", field), value)
		}
	},
	"neq": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s != ?", field), value)
		}
	},
	"gt": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s > ?", field), value)
		}
	},
	"gte": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s >= ?", field), value)
		}
	},
	"lt": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s < ?", field), value)
		}
	},
	"lte": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s <= ?", field), value)
		}
	},
	"in": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s IN (?)", field), value)
		}
	},
	"notIn": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s NOT IN (?)", field), value)
		}
	},
	"like": func(ctx *context.Context, field string, value interface{}) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where(fmt.Sprintf("%s LIKE ?", field), value)
		}
	},
}
