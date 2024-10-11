package graphql

import "github.com/light-speak/lighthouse/graphql/parser"

var Parser *parser.Parser

// GetSchema Get service schema
func GetSchema() string {
	return generateSchema(Parser.Nodes)
}
