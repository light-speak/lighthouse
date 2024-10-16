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

// func TestParseOperation(t *testing.T) {
// 	_, err := ParserSchema([]string{"demo.graphql"})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	p := GetParser()

// 	nl, err := parser.ReadGraphQLFile("query_example.graphql")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	qp := p.NewQueryParser(nl)
// 	qp.Parser.ParseSchema()
// 	for _, node := range qp.FragmentMap {
// 		err := validate.Validate(node, p)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 	}
// 	err = validate.Validate(qp.OperationNode, p)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }
