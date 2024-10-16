package graphql

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/validate"
)

var Parser *parser.Parser

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

