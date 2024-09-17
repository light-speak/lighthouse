package merge

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed merge.gotpl
var mergeTemplate string

type MergeType struct {
	Model      string
	MergeField []*MergeField
}

type MergeField struct {
	Target string
	Source string
	Local  string
}

func GenMergeModels(mergeTypes []*MergeType) error {
	fileName := filepath.Join("graph", "merge.go")
	file, err := utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %v", err)
	}
	defer file.Close()

	data := struct {
		MergeTypes []*MergeType
	}{
		MergeTypes: mergeTypes,
	}

	tmpl := template.Must(template.New("graph").Funcs(template.FuncMap{
		"lcFirst":       utils.LcFirst,
		"ucFirst":       utils.UcFirst,
		"ucFirstWithID": utils.UcFirstWithID,
	}).Parse(mergeTemplate))

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	log.Info("已将合并函数追加到文件: %s", fileName)

	return nil
}
