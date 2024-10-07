package parser

import (
	"testing"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/parser/lexer"
	"github.com/light-speak/lighthouse/parser/parser"
	"github.com/light-speak/lighthouse/parser/validate"
)

func TestReadGraphQLFile(t *testing.T) {
	l, err := parser.ReadGraphQLFile("demo.graphql")
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
	l, err := parser.ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser(l)
	nodes := p.ParseSchema()
	for _, node := range nodes {
		log.Debug().Msgf("Type: %s", node.GetType())
	}
}

func TestValidate(t *testing.T) {
	l, err := parser.ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser(l)
	

	nodes := p.ParseSchema()
	for _, node := range nodes {
		validate.Validate(node, p)
	}
}
