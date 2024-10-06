package parser

import "github.com/light-speak/lighthouse/parser/lexer"

// parseDescription parses a description if present
func (p *Parser) parseDescription() string {
	if p.peekToken().Type == lexer.Message {
		p.expect(lexer.Message)
		description := p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
		return description
	}
	return ""
}
