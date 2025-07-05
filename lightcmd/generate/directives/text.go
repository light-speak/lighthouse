package directives

import (
	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func init() {
	generate.AddDirective("text", text)
}

func text(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], `type:text`)
	return logic, nil
}
