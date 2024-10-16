package validate

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
)

var p *parser.Parser
var store *ast.NodeStore

func ValidateNodes(nodes map[string]ast.Node, parser *parser.Parser) error {
	p = parser
	store = parser.NodeStore
	for _, node := range nodes {
		err := node.Validate(store)
		if err != nil {
			return err
		}
	}
	for _, directive := range store.Directives {
		if err := directive.Validate(store); err != nil {
			return err
		}
	}
	return nil
}
