package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseTypeReference parse a type reference
// ID, String, [Int], [[Int]], [User], [[User]], [User!]
func (p *Parser) parseTypeReference() *ast.FieldType {
	var fieldType *ast.FieldType
	if p.currToken.Type == lexer.LeftBracket {
		p.expect(lexer.LeftBracket) // skip [
		elemType := p.parseTypeReference()
		p.expect(lexer.RightBracket) // skip ]
		fieldType = &ast.FieldType{
			Name:     "List",
			IsList:   true,
			ElemType: elemType,
		}
	} else {
		fieldType = &ast.FieldType{
			Name: p.currToken.Value,
		}
		p.expect(lexer.Letter)
	}
	if p.currToken.Type == lexer.Exclamation {
		fieldType.IsNonNull = true
		p.expect(lexer.Exclamation)
	}
	return fieldType
}
