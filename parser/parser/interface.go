package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseInterface parse an interface node
//
//	interface Node {
//	  id: ID!
//	}
func (p *Parser) parseInterface() *ast.InterfaceNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.LeftBrace)

	node := &ast.InterfaceNode{Name: name, Description: description}

	var fields []ast.FieldNode
	for p.currToken.Type != lexer.RightBrace {
		fields = append(fields, p.parseField(node))
	}

	node.Fields = fields

	p.expect(lexer.RightBrace)
	if p.InterfaceMap == nil {
		p.InterfaceMap = make(map[string]*ast.InterfaceNode)
	}
	p.InterfaceMap[name] = node
	return node
}
