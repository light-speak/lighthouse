package one

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/light-speak/lighthouse/template"
	"github.com/light-speak/lighthouse/version"
)

var projectName string
var projectModule string
var dirs = []string{"cmd", "schema", "service", "models", "resolver", "repo"}

//go:embed tpl
var oneFs embed.FS

func Run(module string) error {
	projectModule = module
	projectName = filepath.Base(module)
	for _, dir := range dirs {
		template.AddImportRegex(fmt.Sprintf(`%s\.`, dir), fmt.Sprintf("%s/%s", projectModule, dir), "")
	}

	initFunctions := []func() error{
		InitDir,
		InitEnv,
		InitMain,
		InitCmd,
		InitMod,
		InitIdeHelper,
		InitLighthouseYml,
		InitResolver,
	}

	for _, initFunc := range initFunctions {
		if err := initFunc(); err != nil {
			return err
		}
	}

	return nil
}

func InitDir() error {
	var err error
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Join(currentDir, projectName), 0755)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		err = os.MkdirAll(filepath.Join(currentDir, projectName, dir), 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitEnv() error {
	envTemplate, err := oneFs.ReadFile("tpl/env.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName),
		Template:     string(envTemplate),
		FileName:     ".env",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func InitMain() error {
	mainTemplate, err := oneFs.ReadFile("tpl/main.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName),
		Template:     string(mainTemplate),
		FileName:     "main",
		FileExt:      "go",
		Package:      "main",
		Editable:     true,
		SkipIfExists: true,
		Data: map[string]string{
			"Module": projectModule,
		},
	}
	return template.Render(options)
}

func InitCmd() error {
	cmdTemplate, err := oneFs.ReadFile("tpl/cmd.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName, "cmd"),
		Template:     string(cmdTemplate),
		FileName:     "cmd",
		FileExt:      "go",
		Package:      "cmd",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func InitMod() error {
	modTemplate, err := oneFs.ReadFile("tpl/mod.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName),
		Template:     string(modTemplate),
		FileName:     "go",
		FileExt:      "mod",
		Editable:     true,
		SkipIfExists: true,
		Data: map[string]string{
			"Module":  projectModule,
			"Version": version.Version,
		},
	}
	return template.Render(options)
}

func InitIdeHelper() error {
	ideHelperTemplate, err := oneFs.ReadFile("tpl/ide-helper.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName),
		Template:     string(ideHelperTemplate),
		FileName:     "ide-helper",
		FileExt:      "graphql",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func InitLighthouseYml() error {
	lighthouseTemplate, err := oneFs.ReadFile("tpl/lighthouse.tpl")
	if err != nil {
		return err
	}
	options := &template.Options{
		Path:         filepath.Join(projectName),
		Template:     string(lighthouseTemplate),
		FileName:     "lighthouse",
		FileExt:      "yml",
		Editable:     true,
		SkipIfExists: true,
	}
	return template.Render(options)
}

func InitResolver() error {
	resolverTemplate, err := oneFs.ReadFile("tpl/resolver.tpl")
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
