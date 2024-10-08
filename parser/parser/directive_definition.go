package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseDirective parse a directive node
// directive @skip(if: true) on XXX | XXXX
func (p *Parser) parseDirectiveDefinition() *ast.DirectiveDefinitionNode {
	description := p.parseDescription()

	p.expect(lexer.Directive)
	p.expect(lexer.At)
	name := p.currToken.Value
	p.expect(lexer.Letter) // skip name

	node := &ast.DirectiveDefinitionNode{Name: name, Description: description}

	var args []*ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments(node)
	}
	node.Args = args
	p.expect(lexer.On) // skip ON
	locations := p.parseLocations()

	node.Locations = locations
	if p.DirectiveMap == nil {
		p.DirectiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	p.DirectiveMap[name] = node
	return node
}

// parseLocations parse locations
// ON FIELD_DEFINITION | ARGUMENT_DEFINITION | INTERFACE | UNION | ENUM | INPUT_OBJECT | SCALAR | OBJECT
func (p *Parser) parseLocations() []ast.Location {
	var locations []ast.Location

	for {
		location := ast.Location(p.currToken.Value)
		locations = append(locations, location)
		p.expect(lexer.Letter)
		if p.currToken.Type != lexer.Pipe {
			break
		}
		p.expect(lexer.Pipe)
	}

	return locations
}
