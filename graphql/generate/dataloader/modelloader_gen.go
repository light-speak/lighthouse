package dataloader

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	_ "embed"

	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed modelloader_gen.gotpl
var modelLoaderGenTemplate string

type templateData struct {
	Package      string
	ModelName    string
	ModelPackage string
	Name         string
}

func GenModelLoaderGen(models []*modelgen.Object) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}
	mod, err := utils.GetModulePath()
	if err != nil {
		return fmt.Errorf("获取模块路径失败: %w", err)
	}
	outputDir := filepath.Join(wd, "graph", "generate")
	for _, model := range models {
		gen(model, outputDir, mod)
	}

	return nil
}

func gen(model *modelgen.Object, outputDir string, mod string) error {
	data := templateData{
		Package:      "generate",
		ModelName:    model.Name,
		ModelPackage: fmt.Sprintf("%s/graph/models", mod),
		Name:         fmt.Sprintf("%sLoader", model.Name),
	}

	filePath := filepath.Join(outputDir, fmt.Sprintf("%s_loader_gen.go", utils.LcFirst(data.Name)))
	file, err := utils.CreateOrTruncateFile(filePath)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %w", err)
	}
	defer file.Close()

	tmpl := template.Must(template.New("modelloader_gen").
		Funcs(template.FuncMap{
			"ucFirst":   utils.UcFirst,
			"lcFirst":   utils.LcFirst,
			"pluralize": utils.Pluralize,
			"toLower":   utils.ToLower,
		}).
		Parse(modelLoaderGenTemplate))
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	return nil
}
