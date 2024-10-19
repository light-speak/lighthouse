package graphql

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/validate"
)

var (
	Parser      *parser.Parser
	QueryParser *parser.Parser
	Schema      *schema
)

type SchemaType struct {
	Name string `json:"name"`
}

type Introspection struct {
	Schema *schema `json:"__schema"`
}

type schema struct {
	Query        *SchemaType                `json:"queryType"`
	Mutation     *SchemaType                `json:"mutationType"`
	Subscription *SchemaType                `json:"subscriptionType"`
	Directives   []*ast.DirectiveDefinition `json:"directives"`
	Types        []ast.Node                 `json:"types"`
}

// // GetSchema Get service schema
// func GetSchema() string {
// 	return generateSchema(Parser.Nodes)
// }

func GetParser() *parser.Parser {
	if Parser == nil {
		panic("Parser is not initialized")
	}
	return Parser
}

func ParserSchema(files []string) (map[string]ast.Node, error) {
	lexer, err := parser.ReadGraphQLFiles(files)
	if err != nil {
		return nil, err
	}

	Parser = parser.NewParser(lexer)
	nodes := Parser.ParseSchema()
	err = validate.ValidateNodes(nodes, Parser)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func ExecuteQuery(qp *parser.QueryParser) (interface{}, error) {
	var err error
	res := map[string]interface{}{}
	for _, field := range qp.Fields {
		res[field.Name], err = ResolveIntrospectionSchema(qp, field)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
