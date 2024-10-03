package parser

import (
	"testing"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/parser/lexer"
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

func TestParseSchema(t *testing.T) {
	l, err := ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(l)
	nodes := p.ParseSchema()
	for _, node := range nodes {
		log.Debug().Msgf("Type: %s", node.GetType())
	}
}
