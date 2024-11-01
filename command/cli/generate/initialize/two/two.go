package two

import (
	"embed"
	"path/filepath"

	"github.com/light-speak/lighthouse/template"
)

//go:embed tpl
var twoFs embed.FS

var projectName string
var projectModule string

func Run(module string) error {
	projectModule = module
	projectName = filepath.Base(module)
	err := initRepo()
	if err != nil {
		return err
	}
	err = initService()
	if err != nil {
		return err
	}
	err = initModel()
	if err != nil {
		return err
	}
	err = initStart()
	if err != nil {
		return err
	}
	err = initMigrate()
	if err != nil {
		return err
	}
	return nil
}

func initRepo() error {
	repoTemplate, err := twoFs.ReadFile("tpl/repo.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "repo"),
		Template:     string(repoTemplate),
		FileName:     "repo",
		FileExt:      "go",
		Package:      "repo",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func initService() error {
	serviceTemplate, err := twoFs.ReadFile("tpl/service.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "service"),
		Template:     string(serviceTemplate),
		FileName:     "service",
		FileExt:      "go",
		Package:      "service",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func initModel() error {
	modelTemplate, err := twoFs.ReadFile("tpl/model.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "models"),
		Template:     string(modelTemplate),
		FileName:     "model",
		FileExt:      "go",
		Package:      "models",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func initStart() error {
	startTemplate, err := twoFs.ReadFile("tpl/start.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "cmd"),
		Template:     string(startTemplate),
		FileName:     "start",
		FileExt:      "go",
		Package:      "cmd",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func initMigrate() error {
	migrateTemplate, err := twoFs.ReadFile("tpl/migrate.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "cmd"),
		Template:     string(migrateTemplate),
		FileName:     "migrate",
		FileExt:      "go",
		Package:      "cmd",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}
