package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"
)

var quickExcuteMap = map[string]func(ctx *context.Context, field *ast.Field, scopes ...func(db *gorm.DB) *gorm.DB) (interface{}, errors.GraphqlErrorInterface){
	"find":     executeFind,
	"first":    executeFirst,
	"paginate": executePaginate,
}

func QuickExecute(ctx *context.Context, field *ast.Field) (interface{}, bool, errors.GraphqlErrorInterface) {
	scopes := make([]func(db *gorm.DB) *gorm.DB, 0)
	for _, arg := range field.DefinitionArgs {
		if len(arg.Directives) > 0 {
			var v interface{}
			argValue := field.Args[arg.Name]
			if argValue != nil {
				v = argValue.Value
			}
			scope, err := executeFilter(ctx, arg, v)
			if err != nil {
				return nil, false, err
			}
			scopes = append(scopes, scope)
		}
	}
	for _, directive := range field.DefinitionDirectives {
		if fn, ok := quickExcuteMap[directive.Name]; ok {
			if directive.GetArg("scopes") != nil {
				names := directive.GetArg("scopes").Value.([]interface{})
				for _, name := range names {
					n := fmt.Sprintf("%s%s", utils.UcFirst(utils.CamelCase(field.Type.GetRealType().Name)), utils.UcFirst(utils.CamelCase(name.(string))))
					fn := model.GetScopes(n)
					if fn == nil {
						log.Warn().Str("name", n).Msg("scope not found")
						continue
					}
					scopes = append(scopes, fn(ctx))
				}
			}
			res, err := fn(ctx, field, scopes...)
			if err != nil {
				return nil, true, err
			}
			if err := field.Type.ValidateValue(res, true); err != nil {
				return nil, true, err
			}
			return res, true, nil
		}
	}

	return nil, false, nil
}
