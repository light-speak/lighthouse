package directives

import (
	"fmt"

	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func init() {
	generate.AddDirective("default", defaultDirective)
}

func defaultDirective(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	if valueArg := directive.Arguments.ForName("value"); valueArg != nil && valueArg.Value != nil {
		logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], fmt.Sprintf(`default:%s`, valueArg.Value.Raw))
	}
	return logic, nil
}
