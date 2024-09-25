package dataloader

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"text/template"

	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed modelloader.gotpl
var modelloaderTamplate string

func GenModelLoader(models []*modelgen.Object) error {
	fileName := filepath.Join("graph", "generate", "loaders.go")
	file, err := utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %v", err)
	}
	defer file.Close()

	names := make([]string, len(models))
	for i, model := range models {
		names[i] = model.Name
	}

	packageName, err := utils.GetModulePath()
	if err != nil {
		return fmt.Errorf("获取模块路径失败: %v", err)
	}

	data := struct {
		Names   []string
		Package string
	}{
		Names:   names,
		Package: packageName,
	}

	tmpl := template.Must(template.New("generate").Funcs(template.FuncMap{
		"lcFirst": utils.LcFirst,
		"ucFirst": utils.UcFirst,
	}).Parse(modelloaderTamplate))

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	log.Info("已将加载器追加到文件: %s", fileName)

	return nil
}
