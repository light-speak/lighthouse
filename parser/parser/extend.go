package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseExtend parse an extend node
//
//	extend type User {
//	  role: Role!
//	}
func (p *Parser) parseExtend() *ast.TypeNode {
	p.parseDescription() // Skip extend description

	p.nextToken()               // Skip 'extend'
	p.expect(lexer.Type, false) // Ensure the next token is 'type', but not move to next token, continue parsing

	// Parse the extended type using parseType
	return p.parseType()
}
