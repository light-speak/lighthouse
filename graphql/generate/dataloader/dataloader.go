package dataloader

import (
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/dataloaden/pkg/generator"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed middleware.gotpl
var middlewareTemplate string

func genDataLoaders(models []*modelgen.Object) error {
	var file *os.File
	var err error

	fileName := filepath.Join("graph", fmt.Sprintf("middleware.go"))
	file, err = utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Models []*modelgen.Object
	}{
		Models: models,
	}

	tmpl := template.Must(template.New("graph").Parse(middlewareTemplate))
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	log.Info("Appended middleware to file: %s\n", fileName)
	return nil
}

func genAllDataloader(models []*modelgen.Object) error {
	err := genDataLoaders(models)
	if err != nil {
		return err
	}
	err = genDataloaderGen(models)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	mod, err := utils.GetModulePathFromGoMod()
	if err != nil {
		return err
	}

	for _, model := range models {
		if err := generator.Generate(fmt.Sprintf("%sLoader", model.Name), "int64", fmt.Sprintf("*%s/graph.%s", mod, model.Name), filepath.Join(wd, "graph")); err != nil {
			log.Error("%s", err)
		}
	}

	err = GenModelLoader(models)
	if err != nil {
		return err
	}

	return nil
}
