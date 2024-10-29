package graphql

import (
	"os"
	"path/filepath"

	"github.com/light-speak/lighthouse/config"
	"github.com/light-speak/lighthouse/graphql/ast"
	_ "github.com/light-speak/lighthouse/graphql/ast/directive"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/validate"
	"github.com/light-speak/lighthouse/log"
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

func LoadSchema() error {
	schemaFiles := []string{}
	currPath, err := os.Getwd()
	if err != nil {
		return err
	}
	config, err := config.ReadConfig(currPath) // read yml config
	if err != nil {
		return err
	}
	projectSchemaFiles := []string{}
	for _, path := range config.Schema.Path {
		for _, ext := range config.Schema.Ext {
			projectSchemaFiles = append(projectSchemaFiles, filepath.Join(currPath, path, "*."+ext))
		}
	}
	for _, path := range projectSchemaFiles {
		matches, _ := filepath.Glob(path)
		schemaFiles = append(schemaFiles, matches...)
		for _, match := range matches {
			log.Debug().Msgf("Loading schema from %v", match)
		}
	}

	_, err = ParserSchema(schemaFiles)
	if err != nil {
		return err
	}

	return nil
}
