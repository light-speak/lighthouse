package initialize

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed directive.gotpl
var directiveTemplate string

func generateDirective(currentDir string) error {
	directiveDir := filepath.Join(currentDir, "graph", "server")
	directiveFile := filepath.Join(directiveDir, "directive.go")

	if _, err := os.Stat(directiveFile); err == nil {
		log.Warn("文件 %s 已存在，跳过创建", directiveFile)
		return nil
	}

	if _, err := os.Stat(directiveDir); os.IsNotExist(err) {
		if err := os.MkdirAll(directiveDir, os.ModePerm); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}

	file, err := os.Create(directiveFile)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	packageName, err := utils.GetModulePath()
	if err != nil {
		return fmt.Errorf("获取包名失败: %w", err)
	}
	log.Info(packageName)

	data := struct {
		Package string
	}{
		Package: packageName,
	}

	tmpl := template.Must(template.New("directive").Parse(directiveTemplate))
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	log.Info("已创建文件: %s", directiveFile)
	return nil
}
