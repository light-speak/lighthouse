package resolver

import (
	_ "embed"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed resolver.gotpl
var resolverTemplate string



func GenerateResolver() error {
	fileName := filepath.Join("resolver", "config.go")
	file, err := utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %v", err)
	}
	defer file.Close()

	packageName, err := utils.GetModulePath()
	if err != nil {
		return fmt.Errorf("获取模块路径失败: %v", err)
	}

	data := struct {
		Package string
	}{
		Package: packageName,
	}

	tmpl, err := template.New("generate").Parse(resolverTemplate)
	if err != nil {
		return fmt.Errorf("解析模板失败: %v", err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	log.Info("已将resolver追加到文件: %s", fileName)
	return nil
}
