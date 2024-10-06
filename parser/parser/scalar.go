package parser

import "github.com/light-speak/lighthouse/parser/ast"

// parseScalar parse a scalar node
// scalar DateTime
func (p *Parser) parseScalar() *ast.ScalarNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	node := &ast.ScalarNode{Name: name, Description: description}
	if p.scalarMap == nil {
		p.scalarMap = make(map[string]*ast.ScalarNode)
	}
	p.scalarMap[name] = node
	return node
}
