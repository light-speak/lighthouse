package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseUnion parse a union node
// union User = Product | Order
func (p *Parser) parseUnion() *ast.UnionNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.Equal)

	var types []string
	for {
		p.nextToken()
		types = append(types, p.currToken.Value)
		if p.peekToken().Type != lexer.Pipe {
			break
		}
		p.nextToken() // consume the '|'
	}
	node := &ast.UnionNode{Name: name, Types: types, Description: description}
	if p.unionMap == nil {
		p.unionMap = make(map[string]*ast.UnionNode)
	}
	p.unionMap[name] = node
	return node
}