package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseScalar parse a scalar node
// scalar DateTime
func (p *Parser) parseScalar() *ast.ScalarNode {
	description := p.parseDescription()

	p.expect(lexer.Scalar)
	name := p.currToken.Value
	p.nextToken()

	directives := p.parseDirectives()
	node := &ast.ScalarNode{Name: name, Description: description, Directives: directives}

	p.AddScalar(node)
	return node
}
