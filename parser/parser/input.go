package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseInput parse an input node
//
//	input User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseInput() *ast.InputNode {
	// Input types are similar to regular types
	description := p.parseDescription()

	p.expect(lexer.Input)
	name := p.currToken.Value
	p.nextToken()

	p.expect(lexer.LeftBrace)

	node := &ast.InputNode{
		Name:        name,
		Description: description,
	}

	var fields []*ast.FieldNode

	for p.currToken.Type != lexer.RightBrace {
		field := p.parseField(node)
		fields = append(fields, field)
	}
	node.Fields = fields

	p.AddInput(node)
	return node
}
