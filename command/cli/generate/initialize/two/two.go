package two

import (
	"embed"
	"fmt"
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
	err = initResolver()
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
	err = initSearcher()
	if err != nil {
		return err
	}
	err = initQueue()
	if err != nil {
		return err
	}
	err = initQueueStart()
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

func initResolver() error {
	resolverTemplate, err := twoFs.ReadFile("tpl/resolver.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "resolver"),
		Template:     string(resolverTemplate),
		FileName:     "resolver",
		FileExt:      "go",
		Package:      "resolver",
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
		Imports: []*template.Import{
			{
				Path:  fmt.Sprintf("%s/resolver", projectModule),
				Alias: "_",
			},
			{
				Path:  fmt.Sprintf("%s/repo", projectModule),
				Alias: "_",
			},
			{
				Path:  fmt.Sprintf("%s/models", projectModule),
				Alias: "_",
			},
		},
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

func initSearcher() error {
	searcherTemplate, err := twoFs.ReadFile("tpl/searcher.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "cmd"),
		Template:     string(searcherTemplate),
		FileName:     "searcher",
		FileExt:      "go",
		Package:      "cmd",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func initQueue() error {
	queueTemplate, err := twoFs.ReadFile("tpl/queue.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "queue"),
		Template:     string(queueTemplate),
		FileName:     "queue",
		FileExt:      "go",
		Package:      "queue",
		Editable:     false,
		SkipIfExists: false,
		Data:         map[string]interface{}{},
	}
	return template.Render(options)
}

func initQueueStart() error {
	t, err := twoFs.ReadFile("tpl/queue_start.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "cmd"),
		Template:     string(t),
		FileName:     "queue",
		FileExt:      "go",
		Package:      "cmd",
		Editable:     false,
		SkipIfExists: false,
		Data:         map[string]interface{}{},
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
