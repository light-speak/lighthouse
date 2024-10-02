package parser

import (
	"testing"

	"github.com/light-speak/lighthouse/parser/lexer"
	"github.com/light-speak/lighthouse/log"
)

func TestReadGraphQLFile(t *testing.T) {
	l, err := ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	for {
		token := l.NextToken()
		log.Debug().Msgf("%+v", token.Value)
		if token.Type == lexer.EOF {
			break
		}
	}
}
