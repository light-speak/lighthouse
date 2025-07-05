package directives

import (
	"fmt"

	"github.com/light-speak/lighthouse/lightcmd/generate"
	"github.com/vektah/gqlparser/v2/ast"
)

func init() {
	generate.AddDirective("varchar", varchar)
}

func varchar(directive *ast.Directive, logic *generate.DirectiveLogic) (*generate.DirectiveLogic, error) {
	if nameArg := directive.Arguments.ForName("length"); nameArg != nil && nameArg.Value != nil {
		logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], fmt.Sprintf(`type:varchar(%s)`, nameArg.Value.String()))
	} else {
		logic.TagKvs["gorm"] = append(logic.TagKvs["gorm"], `type:varchar(255)`)
	}
	return logic, nil
}
