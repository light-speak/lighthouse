package generate

import (
	"embed"
	"path/filepath"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/template"
)

//go:embed tpl
var modelFs embed.FS

var excludeType = map[string]struct{}{
	"PaginateInfo": {},
	"Query":        {},
	"Mutation":     {},
	"Subscription": {},
}

func GenInterface(nodes map[string]*ast.InterfaceNode, path string) error {
	interfaceTemplate, err := modelFs.ReadFile("tpl/interface.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "models"),
		Template:     string(interfaceTemplate),
		FileName:     "interface",
		FileExt:      "go",
		Package:      "models",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Nodes": nodes,
		},
	}
	return template.Render(options)
}

func GenObject(nodes []*ast.ObjectNode, path string) error {
	filteredNodes := []*ast.ObjectNode{}
	for _, node := range nodes {
		if _, ok := excludeType[node.GetName()]; ok {
			continue
		}
		filteredNodes = append(filteredNodes, node)
	}
	modelTemplate, err := modelFs.ReadFile("tpl/model.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "models"),
		Template:     string(modelTemplate),
		FileName:     "model",
		FileExt:      "go",
		Package:      "models",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Nodes": filteredNodes,
		},
	}
	return template.Render(options)
}

func GenResponse(nodes []*ast.ObjectNode, path string) error {
	filteredNodes := []*ast.ObjectNode{}
	for _, node := range nodes {
		if _, ok := excludeType[node.GetName()]; ok {
			continue
		}
		filteredNodes = append(filteredNodes, node)
	}
	responseTemplate, err := modelFs.ReadFile("tpl/response.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "models"),
		Template:     string(responseTemplate),
		FileName:     "response",
		FileExt:      "go",
		Package:      "models",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Nodes": filteredNodes,
		},
	}
	return template.Render(options)
}

func GenInput(nodes map[string]*ast.InputObjectNode, path string) error {
	inputTemplate, err := modelFs.ReadFile("tpl/input.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "models"),
		Template:     string(inputTemplate),
		FileName:     "input",
		FileExt:      "go",
		Package:      "models",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Nodes": nodes,
		},
	}
	return template.Render(options)
}

func GenRepo(nodes []*ast.ObjectNode, path string) error {
	repoTemplate, err := modelFs.ReadFile("tpl/repo.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "repo"),
		Template:     string(repoTemplate),
		FileName:     "repo",
		FileExt:      "go",
		Package:      "repo",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Nodes": nodes,
		},
	}
	return template.Render(options)
}