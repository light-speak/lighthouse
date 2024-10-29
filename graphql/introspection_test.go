package graphql

import (
	"encoding/json"
	"testing"

	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/log"
)

func TestIntrospectionQuery(t *testing.T) {
	var err error
	_, err = ParserSchema([]string{"demo.graphql"})
	if err != nil {
		t.Fatal(err)
	}
	// Parser.NodeDetail(Parser.NodeStore.Nodes)
	nl, err := parser.ReadGraphQLFile("schema_query.graphql")
	if err != nil {
		t.Fatal(err)
	}
	qp := Parser.NewQueryParser(nl)
	qp = qp.ParseSchema()
	err = qp.Validate(Parser.NodeStore)
	if err != nil {
		t.Fatal(err)
	}
	// json, err := json.Marshal(qp.Fragments)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// log.Warn().Msgf("qp.Fragments: %+v", string(json))

	res := map[string]interface{}{}
	for _, field := range qp.Fields {
		res[field.Name], err = ResolveSchemaFields(qp, field)
		if err != nil {
			t.Fatal(err)
		}
	}
	json, err := json.Marshal(res)
	if err != nil {
		t.Fatal(err)
	}
	log.Info().Msgf("d: %+v", string(json))
}
