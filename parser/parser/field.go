package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseField parse a field node
// "It is ID"
// id: ID!
// name: String!
// age: Int
// email: String
// createdAt: DateTime
func (p *Parser) parseField(parent ast.Node) ast.FieldNode {
	description := p.parseDescription()

	name := p.currToken.Value
	p.expect(lexer.Letter)

	field := &ast.FieldNode{
		Name:        name,
		Description: description,
	}
	// parse arguments
	var args []ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments(field)
	}

	field.Args = args

	p.expect(lexer.Colon)

	fieldType := p.parseTypeReference()
	field.Type = fieldType

	directives := p.parseDirectives()
	field.Directives = directives
	field.Parent = parent

	return *field
}
