package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseUnion parse a union node
// union User = Product | Order
func (p *Parser) parseUnion() *ast.UnionNode {
	description := p.parseDescription()

	p.expect(lexer.Union)
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.Equal)

	var types []string
	for {
		types = append(types, p.currToken.Value)
		p.nextToken()
		if p.currToken.Type != lexer.Pipe {
			break
		}
		p.expect(lexer.Pipe)
	}
	node := &ast.UnionNode{Name: name, Types: types, Description: description}
	if p.UnionMap == nil {
		p.UnionMap = make(map[string]*ast.UnionNode)
	}
	p.UnionMap[name] = node
	return node
}
