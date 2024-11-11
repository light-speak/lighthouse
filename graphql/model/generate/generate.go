package generate

import (
	"embed"
	"fmt"
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

type SearchableField struct {
	Field          *ast.Field
	Type           string
	IndexAnalyzer  *string
	SearchAnalyzer *string
}

func GenInterface(nodes []ast.Node, path string) error {
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
	if err := template.Render(options); err != nil {
		return err
	}

	// Generate scopes for each model
	scopeTemplate, err := modelFs.ReadFile("tpl/scope.tpl")
	if err != nil {
		return err
	}

	for _, node := range filteredNodes {
		if len(node.Scopes) > 0 {
			scopeOptions := &template.Options{
				Path:         filepath.Join(path, "models"),
				Template:     string(scopeTemplate),
				FileName:     fmt.Sprintf("%s_scopes", utils.LcFirst(utils.SnakeCase(node.GetName()))),
				FileExt:      "go",
				Package:      "models",
				Editable:     false,
				SkipIfExists: false,
				Data: map[string]interface{}{
					"Name":   node.GetName(),
					"Scopes": node.Scopes,
				},
			}
			if err := template.Render(scopeOptions); err != nil {
				return err
			}
		}
	}

	return nil
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

func GenSearchable(fields map[*ast.ObjectNode][]*SearchableField, path string) error {
	searchTemplate, err := modelFs.ReadFile("tpl/searchable.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "models"),
		Template:     string(searchTemplate),
		FileName:     "searchable",
		FileExt:      "go",
		Package:      "models",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Fields": fields,
		},
	}
	return template.Render(options)
}

func GenJob(name string, path string) error {
	jobTemplate, err := modelFs.ReadFile("tpl/job.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "queue"),
		Template:     string(jobTemplate),
		FileName:     utils.SnakeCase(name),
		FileExt:      "go",
		Package:      "queue",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Name": name,
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

	for _, field := range node.Fields {
		hasQuickDirective := false
		for _, directive := range field.Directives {
			if ast.IsQuickDirective(directive) {
				hasQuickDirective = true
				break
			}
		}
		if !hasQuickDirective && !utils.IsInternalType(field.Name) {
			options := &template.Options{
				Path:         filepath.Join(path, "resolver"),
				Template:     string(operationTemplate),
				FileName:     utils.SnakeCase(field.Name),
				FileExt:      "go",
				Package:      "resolver",
				Editable:     true,
				SkipIfExists: false,
				Data: map[string]interface{}{
					"Fields": []*ast.Field{field},
				},
			}
			err = template.Render(options)
			if err != nil {
				return err
			}
		}
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

func GenAttr(nodes map[string][]*ast.Field, path string) error {
	// 生成每个类型的 attr resolver 文件
	for typeName, fields := range nodes {
		for _, field := range fields {
			attrTemplate, err := modelFs.ReadFile("tpl/attr.tpl")
			if err != nil {
				return err
			}
			options := &template.Options{
				Path:         filepath.Join(path, "resolver"),
				Template:     string(attrTemplate),
				FileName:     fmt.Sprintf("%s_attr", utils.SnakeCase(field.Name)),
				FileExt:      "go",
				Package:      "resolver",
				Editable:     true,
				SkipIfExists: false,
				Data: map[string]interface{}{
					"Fields": map[string][]*ast.Field{
						typeName: {field},
					},
				},
			}
			if err := template.Render(options); err != nil {
				return err
			}
		}
	}

	attrGenTemplate, err := modelFs.ReadFile("tpl/attr_gen.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(path, "resolver"),
		Template:     string(attrGenTemplate),
		FileName:     "attr_gen",
		FileExt:      "go",
		Package:      "resolver",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]interface{}{
			"Fields": nodes,
		},
	}
	return template.Render(options)
}
