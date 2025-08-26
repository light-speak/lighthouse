package initization

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/templates"
	"github.com/light-speak/lighthouse/utils"
)

//go:embed tpl
var tpl embed.FS

var projectName string
var projectModule string
var currentDir string
var models []string

var dirs = []string{
	"commands",
	"configs",
	"server",
	"graph",
	"models",
	"schema",
	"resolver",
}

func Run(module string, ms string) error {
	logs.Info().Msgf("module: %s", module)
	projectModule = module
	projectName = filepath.Base(module)
	models = strings.Split(ms, ",")
	for _, dir := range dirs {
		templates.AddImportRegex(fmt.Sprintf(`%s\.`, dir), fmt.Sprintf("%s/%s", projectModule, dir), "")
	}

	curPath, err := os.Getwd()
	if err != nil {
		return err
	}
	currentDir = curPath
	initFunctions := []func() error{
		initDir,
		initEnv,
		initMain,
		initCmd,
		initGqlgen,
		initModels,
		initResolvers,
		initServer,
		initMod,
		initConfig,
		initAppStart,
		initGitignore,
		initGraphql,
	}
	templates.AddImportRegex("cmd", "github.com/light-speak/lighthouse/lightcmd/cmd", "")
	templates.AddImportRegex("handler", "github.com/99designs/gqlgen/graphql/handler", "")
	templates.AddImportRegex("playground", "github.com/99designs/gqlgen/graphql/playground", "")
	templates.AddImportRegex("logs", "github.com/light-speak/lighthouse/logs", "")
	templates.AddImportRegex("http", "net/http", "")
	templates.AddImportRegex("utils", "github.com/light-speak/lighthouse/utils", "")
	templates.AddImportRegex("godotenv", "github.com/joho/godotenv", "")
	templates.AddImportRegex("filepath", "path/filepath", "")
	templates.AddImportRegex("gorm", "gorm.io/gorm", "")
	templates.AddImportRegex("databases", "github.com/light-speak/lighthouse/databases", "")
	templates.AddImportRegex("lighterr", "github.com/light-speak/lighthouse/lighterr", "")
	templates.AddImportRegex("lru", "github.com/99designs/gqlgen/graphql/handler/lru", "")
	templates.AddImportRegex("extension", "github.com/99designs/gqlgen/graphql/handler/extension", "")
	templates.AddImportRegex("transport", "github.com/99designs/gqlgen/graphql/handler/transport", "")
	templates.AddImportRegex("ast", "github.com/vektah/gqlparser/v2/ast", "")
	templates.AddImportRegex("time", "time", "")
	templates.AddImportRegex("dataloader", "github.com/light-speak/lighthouse/routers/dataloader", "")

	for _, fn := range initFunctions {
		if err := fn(); err != nil {
			return err
		}
	}

	err = os.Chdir(filepath.Join(currentDir, projectName))
	if err != nil {
		return err
	}
	logs.Info().Msgf("project %s init success, changed to %s directory", projectName, currentDir)
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("lightcmd", "generate:schema")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func initDir() error {
	err := utils.MkdirAll(filepath.Join(currentDir, projectName))
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		err := utils.MkdirAll(filepath.Join(currentDir, projectName, dir))
		if err != nil {
			return err
		}
	}
	return nil
}

func initEnv() error {
	envTpl, err := tpl.ReadFile("tpl/env.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(envTpl),
		FileName:     ".env",
		Editable:     true,
		SkipIfExists: true,
	}
	return templates.Render(options)
}

func initMod() error {
	modTpl, err := tpl.ReadFile("tpl/mod.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(modTpl),
		FileName:     "go",
		FileExt:      "mod",
		Editable:     true,
		SkipIfExists: true,
		Data: map[string]string{
			"Module": projectModule,
		},
	}
	return templates.Render(options)
}

func initModels() error {
	schemaTpl, err := tpl.ReadFile("tpl/schema.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "schema"),
		Template:     string(schemaTpl),
		FileName:     "schema",
		FileExt:      "graphql",
		Editable:     true,
		SkipIfExists: true,
	}
	err = templates.Render(options)
	if err != nil {
		return err
	}

	for _, model := range models {
		modelTpl, err := tpl.ReadFile("tpl/model.tpl")
		if err != nil {
			return err
		}
		options := &templates.Options{
			Path:         filepath.Join(projectName, "schema"),
			Template:     string(modelTpl),
			FileName:     model,
			FileExt:      "graphql",
			Editable:     true,
			SkipIfExists: true,
			Data: map[string]any{
				"Model": model,
			},
		}
		err = templates.Render(options)
		if err != nil {
			return err
		}
	}

	return nil
}

func initMain() error {
	mainTpl, err := tpl.ReadFile("tpl/main.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(mainTpl),
		FileName:     "main",
		Package:      "main",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
		Data: map[string]any{
			"Module": projectModule,
		},
	}
	return templates.Render(options)
}

func initGraphql() error {
	graphqlTpl, err := tpl.ReadFile("tpl/graphql.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(graphqlTpl),
		FileName:     "",
		FileExt:      "graphql",
		Editable:     true,
		SkipIfExists: true,
	}
	return templates.Render(options)
}

func initCmd() error {
	cmdTpl, err := tpl.ReadFile("tpl/command.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "commands"),
		Template:     string(cmdTpl),
		FileName:     "command",
		Package:      "commands",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
	}
	return templates.Render(options)
}

func initGqlgen() error {
	gqlgenTpl, err := tpl.ReadFile("tpl/gqlgen.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(gqlgenTpl),
		FileName:     "gqlgen",
		FileExt:      "yml",
		Editable:     true,
		SkipIfExists: true,
	}
	err = templates.Render(options)
	if err != nil {
		return err
	}

	return nil
}

func initResolvers() error {
	resolverTpl, err := tpl.ReadFile("tpl/resolver.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "resolver"),
		Template:     string(resolverTpl),
		FileName:     "resolver",
		FileExt:      "go",
		Editable:     false,
		SkipIfExists: true,
	}
	err = templates.Render(options)
	if err != nil {
		return err
	}
	return nil
}

func initServer() error {
	serverTpl, err := tpl.ReadFile("tpl/server.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "server"),
		Template:     string(serverTpl),
		FileName:     "server",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
	}

	templates.AddImportRegex("context", "context", "")
	templates.AddImportRegex("routers", "github.com/light-speak/lighthouse/routers", "")
	templates.AddImportRegex("auth", "github.com/light-speak/lighthouse/routers/auth", "")
	templates.AddImportRegex("gqlerror", "github.com/vektah/gqlparser/v2/gqlerror", "")

	err = templates.Render(options)
	if err != nil {
		return err
	}
	return nil
}

func initConfig() error {
	configTpl, err := tpl.ReadFile("tpl/config.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "configs"),
		Template:     string(configTpl),
		FileName:     "config",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
	}

	return templates.Render(options)
}

func initAppStart() error {
	appstartTpl, err := tpl.ReadFile("tpl/appstart.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "commands"),
		Template:     string(appstartTpl),
		FileName:     "app-start",
		FileExt:      "go",
		Editable:     false,
		SkipIfExists: true,
	}
	return templates.Render(options)
}

func initGitignore() error {
	gitignoreTpl, err := tpl.ReadFile("tpl/gitignore.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(gitignoreTpl),
		FileName:     ".gitignore",
		FileExt:      "",
		Editable:     true,
		SkipIfExists: true,
	}
	return templates.Render(options)
}
