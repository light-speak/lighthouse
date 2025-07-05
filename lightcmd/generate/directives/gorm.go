package directives

import (
	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func init() {
	generate.AddDirective("gorm", gorm)
}

func gorm(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	if valueArg := directive.Arguments.ForName("value"); valueArg != nil && valueArg.Value != nil {
		logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], valueArg.Value.Raw)
	}
	return logic, nil
}
