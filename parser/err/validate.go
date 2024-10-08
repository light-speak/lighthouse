package err

import (
	"fmt"

	"github.com/light-speak/lighthouse/parser/ast"
)

type ValidateError struct {
	Node    ast.Node
	Message string
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("validate error: %s, node: %s", e.Message, e.Node.GetName())
}