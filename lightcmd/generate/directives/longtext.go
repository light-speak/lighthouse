package directives

import (
	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func init() {
	generate.AddDirective("longtext", longtext)
}

func longtext(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], `type:longtext`)
	return logic, nil
}
