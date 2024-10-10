package graphql

import (
	"testing"

	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/graphql/validate"
	"github.com/light-speak/lighthouse/log"
)

func TestReadGraphQLFile(t *testing.T) {
	l, err := parser.ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	for {
		token, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}
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
		log.Debug().Msgf("Type: %s", node.GetNodeType())
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
		err := validate.Validate(node, p)
		if err != nil {
			t.Fatal(err)
		}
	}

	schema := generateSchema(nodes)
	log.Debug().Msgf("schema: \n\n%s", schema)
}
