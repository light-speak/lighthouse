package directives

import (
	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func unique(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], `unique`)
	return logic, nil
}

func init() {
	generate.AddDirective("unique", unique)
}
