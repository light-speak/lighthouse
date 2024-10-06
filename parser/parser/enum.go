package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseEnum parse an enum node
//
//	enum Role {
//	  ADMIN
//	  USER
//	  GUEST
//	}
func (p *Parser) parseEnum() *ast.EnumNode {
	description := p.parseDescription()

	p.expect(lexer.Enum)
	name := p.currToken.Value
	p.expect(lexer.Letter)


	directives := p.parseDirectives()

	p.expect(lexer.LeftBrace)

	var values []string
	for p.currToken.Type != lexer.RightBrace {
		values = append(values, p.currToken.Value)
		p.expect(lexer.Letter)
	}

	p.expect(lexer.RightBrace)
	
	node := &ast.EnumNode{
		Name:        name,
		Values:      values,
		Description: description,
		Directives:  directives,
	}
	if p.enumMap == nil {
		p.enumMap = make(map[string]*ast.EnumNode)
	}
	p.enumMap[name] = node
	return node
}
