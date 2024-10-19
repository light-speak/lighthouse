package graphql

import (
	"testing"

	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/log"
)

func TestReadGraphQLFile(t *testing.T) {
	l, err := parser.ReadGraphQLFile("base.graphql")
	if err != nil {
		t.Fatal(err)
	}
	for {
		token, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		log.Debug().Str("type", string(token.Type)).Str("value", token.Value).Msg("")
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
	p.NodeDetail(nodes)
}

func TestValidate(t *testing.T) {
	nodes, err := ParserSchema([]string{"demo.graphql"})
	if err != nil {
		t.Fatal(err)
	}
	p := GetParser()
	p.NodeDetail(nodes)
}

func TestParseOperation(t *testing.T) {
	_, err := ParserSchema([]string{"demo.graphql"})
	if err != nil {
		t.Fatal(err)
	}
	p := GetParser()

	nl, err := parser.ReadGraphQLFile("query_example.graphql")
	if err != nil {
		t.Fatal(err)
	}

	qp := p.NewQueryParser(nl)
	err = qp.Validate(p.NodeStore)
	if err != nil {
		t.Fatal(err)
	}
	log.Debug().Msgf("qp: %+v", qp.Fields["getUser"].Children["name"])
	// log.Debug().Msgf("qp: %+v", qp.Fields["getUser"].Children["result"].Children["User"])
}
