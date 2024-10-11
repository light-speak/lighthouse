package graphql

import (
	"os"
	"path/filepath"

	"github.com/light-speak/lighthouse/config"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/validate"
	"github.com/light-speak/lighthouse/template"
)

func Generate() error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	config, err := config.ReadConfig(currentPath)
	if err != nil {
		return err
	}

	schemaFiles := []string{}
	for _, path := range config.Schema.Path {
		for _, ext := range config.Schema.Ext {
			schemaFiles = append(schemaFiles, filepath.Join(currentPath, path, "*."+ext))
		}
	}

	files := make([]string, 0)
	for _, path := range schemaFiles {
		matches, _ := filepath.Glob(path)
		files = append(files, matches...)
	}

	lexer, err := parser.ReadGraphQLFiles(files)
	if err != nil {
		return err
	}

	Parser = parser.NewParser(lexer)
	nodes := Parser.ParseSchema()

	for _, node := range nodes {
		err := validate.Validate(node, Parser)
		if err != nil {
			return err
		}
	}

	typeNodes := []*ast.TypeNode{}
	responseNodes := []*ast.TypeNode{}

	for _, node := range nodes {
		switch node.GetNodeType() {
		case ast.NodeTypeType:
			typeNode, _ := node.(*ast.TypeNode)
			if typeNode.IsResponse {
				responseNodes = append(responseNodes, typeNode)
			} else {
				typeNodes = append(typeNodes, typeNode)
			}
		}
	}

	if err := model.GenType(typeNodes, currentPath); err != nil {
		return err
	}
	if err := model.GenResponse(responseNodes, currentPath); err != nil {
		return err
	}

	schema := generateSchema(nodes)
	options := &template.Options{
		Path:         currentPath,
		Template:     schema,
		FileName:     "schema",
		FileExt:      "graphql",
		Editable:     false,
		SkipIfExists: false,
	}
	template.Render(options)

	return nil
}
