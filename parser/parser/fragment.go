package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

func (p *Parser) parseFragment() *ast.FragmentNode {
	description := p.parseDescription()
	p.expect(lexer.Fragment)
	name := p.currToken.Value
	p.nextToken()

	p.expect(lexer.On)
	onType := p.currToken.Value
	p.nextToken()

	directives := p.parseDirectives()

	p.expect(lexer.LeftBrace)

	node := &ast.FragmentNode{
		Name:        name,
		Description: description,
		OnType:      onType,
		Directives:  directives,
	}

	var fields []ast.FieldNode
	for p.currToken.Type != lexer.RightBrace && p.currToken.Type == lexer.Letter {
		field := p.parseField(node)
		fields = append(fields, field)
	}

	p.expect(lexer.RightBrace)
	node.Fields = fields

	return node
}
