package merge

import (
	_ "embed"
	"fmt"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"os"
	"path/filepath"
	"text/template"
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
	var file *os.File
	var err error

	fileName := filepath.Join("graph", fmt.Sprintf("merge.go"))
	file, err = utils.CreateOrTruncateFile(fileName)
	if err != nil {
		return err
	}

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
	err = tmpl.Execute(file, data)

	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	log.Info("Appended loaders to file: %s\n", fileName)

	return nil
}
