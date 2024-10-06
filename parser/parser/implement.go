package parser

import "github.com/light-speak/lighthouse/parser/lexer"

func (p *Parser) parseImplements() []string {
	var implements []string
	if p.currToken.Type != lexer.Implements {
		return implements
	}
	p.expect(lexer.Implements)
	for {
		implements = append(implements, p.currToken.Value)
		p.expect(lexer.Letter)
		if p.currToken.Type == lexer.LeftBrace {
			break
		}
		p.expect(lexer.And)
	}
	return implements
}
