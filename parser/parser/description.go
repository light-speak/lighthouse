package parser

import (
	"strings"

	"github.com/light-speak/lighthouse/parser/lexer"
)

// parseDescription parses a description if present
func (p *Parser) parseDescription() string {
	if p.PreviousToken().Type == lexer.Message {
		description := strings.Split(p.PreviousToken().Value, "\"")[1]
		return description
	}
	return ""
}
