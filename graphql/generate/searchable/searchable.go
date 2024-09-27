package searchable

import (
	_ "embed"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

type SearchableModelOption struct {
	ModelName string
	Fields    []*SearchableFieldOption
}

type SearchableFieldOption struct {
	FieldName      string
	SearchableType string
	IndexAnalyzer  string
	SearchAnalyzer string
}

//go:embed searchable.gotpl
var searchableTemplate string

func GenSearchableModels(options []*SearchableModelOption) error {
	fileName := filepath.Join("graph", "models", "searchable_gen.go")
	file, err := utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %v", err)
	}
	defer file.Close()

	data := struct {
		SearchableModels []*SearchableModelOption
	}{
		SearchableModels: options,
	}

	tmpl := template.Must(template.New("generate").Funcs(template.FuncMap{
		"lcFirst":       utils.LcFirst,
		"ucFirst":       utils.UcFirst,
		"ucFirstWithID": utils.UcFirstWithID,
		"pluralize":     utils.Pluralize,
		"toLower":       utils.ToLower,
	}).Parse(searchableTemplate))

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	log.Info("已将合并函数追加到文件: %s", fileName)

	return nil
}
