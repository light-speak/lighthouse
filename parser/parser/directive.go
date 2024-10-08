package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

func (p *Parser) parseDirectives() []*ast.DirectiveNode {
	var directives []*ast.DirectiveNode
	if p.currToken.Type != lexer.At {
		return directives
	}
	for {
		directives = append(directives, p.parseDirective())
		if p.currToken.Type != lexer.At {
			break
		}
	}

	return directives
}

func (p *Parser) parseDirective() *ast.DirectiveNode {
	p.expect(lexer.At)
	name := p.currToken.Value
	p.nextToken() // skip name

	directive := &ast.DirectiveNode{
		Name: name,
	}

	var args []*ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments(directive)
	}
	directive.Args = args
	return directive
}
