package generate

import (
	"embed"
	"path/filepath"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/template"
	"github.com/light-speak/lighthouse/utils"
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

func GenOperationResolver(node *ast.ObjectNode, path string, name string) error {
	operationTemplate, err := modelFs.ReadFile("tpl/operation.tpl")
	if err != nil {
		return err
	}

	fields := []*ast.Field{}
	for _, field := range node.Fields {
		if len(field.Directives) == 0 && !utils.IsInternalType(field.Name) {
			fields = append(fields, field)
		}
	}

	options := &template.Options{
		Path:         filepath.Join(path, "resolver"),
		Template:     string(operationTemplate),
		FileName:     name,
		FileExt:      "go",
		Package:      "resolver",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Fields": fields,
		},
	}
	err = template.Render(options)
	if err != nil {
		return err
	}
	return nil
}

func GenOperationResolverGen(nodes []*ast.ObjectNode, path string) error {
	operationTemplate, err := modelFs.ReadFile("tpl/operation_gen.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "resolver"),
		Template:     string(operationTemplate),
		FileName:     "operation_gen",
		FileExt:      "go",
		Package:      "resolver",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Nodes": nodes,
		},
	}
	err = template.Render(options)
	if err != nil {
		return err
	}
	return nil
}

func GenEnum(nodes map[string]*ast.EnumNode, path string) error {
	enumTemplate, err := modelFs.ReadFile("tpl/enum.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "models"),
		Template:     string(enumTemplate),
		FileName:     "enum",
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
