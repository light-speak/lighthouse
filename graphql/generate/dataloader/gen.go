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

//go:embed gen.gotpl
var genTemplate string

func genDataloaderGen(models []*modelgen.Object) error {
	var file *os.File
	var err error

	fileName := filepath.Join("graph", fmt.Sprintf("gen.go"))
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

	tmpl := template.Must(template.New("graph").Parse(genTemplate))
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	log.Info("Appended middleware to file: %s\n", fileName)
	return nil
}
