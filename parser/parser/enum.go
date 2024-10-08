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

	node := &ast.EnumNode{
		Name:        name,
		Description: description,
		Directives:  directives,
	}

	var values []*ast.EnumValueNode
	for p.currToken.Type != lexer.RightBrace {
		values = append(values, p.parseEnumValue(node))
	}

	node.Values = values
	p.expect(lexer.RightBrace)

	if p.EnumMap == nil {
		p.EnumMap = make(map[string]*ast.EnumNode)
	}
	p.EnumMap[name] = node
	return node
}

func (p *Parser) parseEnumValue(parent ast.Node) *ast.EnumValueNode {
	description := p.parseDescription()

	name := p.currToken.Value
	p.nextToken()

	directives := p.parseDirectives()

	return &ast.EnumValueNode{
		Name:        name,
		Description: description,
		Directives:  directives,
		Parent:      parent,
	}
}
