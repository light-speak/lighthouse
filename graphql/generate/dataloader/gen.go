package dataloader

import (
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"path/filepath"
	"text/template"
)

//go:embed gen.gotpl
var genTemplate string

func genDataloaderGen(models []*modelgen.Object) error {
	fileName := filepath.Join("graph", "gen.go")
	file, err := utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %v", err)
	}
	defer file.Close()

	tmpl, err := template.New("graph").Parse(genTemplate)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	data := struct {
		Models []*modelgen.Object
	}{
		Models: models,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	log.Info("已将中间件追加到文件: %s", fileName)
	return nil
}
