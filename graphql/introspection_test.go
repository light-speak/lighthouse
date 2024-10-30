package graphql

import (
	"testing"

	"github.com/light-speak/lighthouse/graphql/parser"
)

func TestIntrospectionQuery(t *testing.T) {
	var err error
	_, err = ParserSchema([]string{"demo.graphql"})
	if err != nil {
		t.Fatal(err)
	}

	// Parse schema query
	nl, err := parser.ReadGraphQLFile("schema_query.graphql")
	if err != nil {
		t.Fatal(err)
	}

	// Create query parser and validate
	qp := Parser.NewQueryParser(nl)
	qp = qp.ParseSchema()
	err = qp.Validate(Parser.NodeStore)
	if err != nil {
		t.Fatal(err)
	}

	// Resolve fields
	res := map[string]interface{}{}
	for _, field := range qp.Fields {
		res[field.Name], err = ResolveSchemaFields(qp, field)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Verify results
	if len(res) == 0 {
		t.Error("Expected non-empty result map")
	}
}
