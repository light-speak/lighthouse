package errors

import (
	"fmt"
	"github.com/light-speak/lighthouse/graphql/ast"
)

type ValidateError struct {
	Node    ast.Node
	Message string
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("validate error: %s, node: %s", e.Message, e.Node.GetName())
}

type LexerError struct {
	Path         *string
	Line         int
	LinePosition int
	Message      string
}

func (e *LexerError) Error() string {
	if e.Path != nil {
		return fmt.Sprintf("lexer error: %s, path: %s, line: %d, line position: %d", e.Message, *e.Path, e.Line, e.LinePosition)
	}
	return fmt.Sprintf("lexer error: %s, line: %d, line position: %d", e.Message, e.Line, e.LinePosition)
}
