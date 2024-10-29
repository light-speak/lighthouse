package graphql

import (
	"testing"

	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
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
	var err error
	_, err = ParserSchema([]string{"demo.graphql"})
	if err != nil {
		t.Fatal(err)
	}
	p := GetParser()

	nl, err := parser.ReadGraphQLFile("query_example.graphql")
	if err != nil {
		t.Fatal(err)
	}

	qp := p.NewQueryParser(nl)
	qp.Variables = make(map[string]any)
	qp.Variables["$id"] = "1"
	qp.Variables["$test"] = "1"

	// log.Debug().Msgf("qp: %+v", qp.Fields["getUser"].Children)
	// log.Debug().Msgf("result: %+v", qp.Fields["getUser"].Children["result"].Children)
	// log.Debug().Msgf("user: %+v", qp.Fields["getUser"].Children["result"].Children["User"].Children)
	log.Info().Msgf("qp: %+v", qp.Parser.NodeStore.Objects["Query"].Fields["getUser"].Args["id"].Type.OfType)
	err = qp.Validate(p.NodeStore)
	if err != nil {
		t.Fatal(err)
	}
}
