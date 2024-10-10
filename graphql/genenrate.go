package graphql

import (
	"os"
	"path/filepath"

	"github.com/light-speak/lighthouse/config"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/validate"
	"github.com/light-speak/lighthouse/log"
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

	p := parser.NewParser(lexer)
	nodes := p.ParseSchema()

	for _, node := range nodes {
		err := validate.Validate(node, p)
		if err != nil {
			return err
		}
	}
	schema := generateSchema(nodes)
	log.Debug().Msgf("schema: \n\n%s", schema)

	return nil
}
