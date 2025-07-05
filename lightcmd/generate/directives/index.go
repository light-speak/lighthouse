package directives

import (
	"fmt"

	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func init() {
	generate.AddDirective("index", index)
}

func index(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	if nameArg := directive.Arguments.ForName("name"); nameArg != nil && nameArg.Value != nil {
		logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], fmt.Sprintf(`index:%s`, nameArg.Value.Raw))
	} else {
		logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], `index`)
	}
	return logic, nil
}
