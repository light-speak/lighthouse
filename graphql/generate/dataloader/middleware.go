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

//go:embed middleware.gotpl
var middlewareTemplate string

func genDataLoaders(models []*modelgen.Object) error {
	fileName := filepath.Join("graph", "generate", "middleware.go")
	file, err := utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return fmt.Errorf("创建或截断文件失败: %v", err)
	}
	defer file.Close()

	tmpl := template.Must(template.New("generate").Parse(middlewareTemplate))
	err = tmpl.Execute(file, struct{ Models []*modelgen.Object }{Models: models})
	if err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	log.Info("已将中间件追加到文件: %s", fileName)
	return nil
}

func GenAllDataloader(models []*modelgen.Object) error {
	if err := genDataLoaders(models); err != nil {
		return fmt.Errorf("生成数据加载器失败: %w", err)
	}

	if err := genDataloaderGen(models); err != nil {
		return fmt.Errorf("生成数据加载器生成器失败: %w", err)
	}

	// wd, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("获取当前工作目录失败: %w", err)
	// }

	// mod, err := utils.GetModulePath()
	// if err != nil {
	// 	return fmt.Errorf("获取模块路径失败: %w", err)
	// }
	// for _, model := range models {
	// 	loaderName := fmt.Sprintf("%sLoader", model.Name)
	// 	typePath := fmt.Sprintf("*%s/graph/models.%s", mod, model.Name)
	// 	outputDir := filepath.Join(wd, "graph", "generate")
	// 	log.Warn("loaderName: %s, typePath: %s, outputDir: %s", loaderName, typePath, outputDir)

	// 	// if err := generator.Generate(loaderName, "int64", typePath, outputDir); err != nil {
	// 	// 	log.Error("生成 %s 失败: %v", loaderName, err)
	// 	// 	continue
	// 	// }

	// }
	

	if err := GenModelLoaderGen(models); err != nil {
		return fmt.Errorf("生成模型加载器生成器失败: %w", err)
	}

	if err := GenModelLoader(models); err != nil {
		return fmt.Errorf("生成模型加载器失败: %w", err)
	}

	return nil
}
