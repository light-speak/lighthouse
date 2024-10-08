package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseType parse a type node
//
//	type User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseType() *ast.TypeNode {
	description := p.parseDescription()

	p.expect(lexer.Type)
	name := p.currToken.Value
	p.nextToken()

	implements := p.parseImplements()
	directives := p.parseDirectives()
	p.expect(lexer.LeftBrace)

	node := &ast.TypeNode{
		Name:        name,
		Description: description,
		Implements:  implements,
		Directives:  directives,
	}

	var fields []*ast.FieldNode

	for p.currToken.Type != lexer.RightBrace && p.currToken.Type == lexer.Letter {
		field := p.parseField(node)
		fields = append(fields, field)
	}
	p.expect(lexer.RightBrace)
	node.Fields = fields

	operationType := ast.OperationTypeEntity

	switch name {
	case "Query":
		operationType = ast.OperationTypeQuery
	case "Mutation":
		operationType = ast.OperationTypeMutation
	case "Subscription":
		operationType = ast.OperationTypeSubscription
	}

	node.OperationType = operationType

	p.AddType(name, node)
	return node
}
