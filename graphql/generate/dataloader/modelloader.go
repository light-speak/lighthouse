package dataloader

import (
	_ "embed"
	"fmt"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed modelloader.gotpl
var modelloaderTamplate string

func GenModelLoader(models []*modelgen.Object) error {
	var file *os.File
	var err error

	fileName := filepath.Join("graph", fmt.Sprintf("loaders.go"))
	file, err = utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return err
	}

	names := make([]string, 0)

	for _, model := range models {
		names = append(names, model.Name)
	}

	data := struct {
		Names []string
	}{
		Names: names,
	}

	tmpl := template.Must(template.New("graph").Funcs(template.FuncMap{
		"lcFirst": utils.LcFirst,
		"ucFirst": utils.UcFirst,
	}).Parse(modelloaderTamplate))
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	log.Info("Appended loaders to file: %s\n", fileName)

	return nil
}
