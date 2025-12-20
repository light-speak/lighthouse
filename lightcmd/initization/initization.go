package initization

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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
	"resolver",
}

// initStep represents a single initialization step
type initStep struct {
	name string
	fn   func() error
}

func Run(module string, ms string) error {
	fmt.Println()
	logs.Info().Msg("🚀 Initializing new Lighthouse project...")
	logs.Info().Msgf("   Module: %s", module)
	logs.Info().Msgf("   Project: %s", filepath.Base(module))
	fmt.Println()

	projectModule = module
	projectName = filepath.Base(module)
	models = strings.Split(ms, ",")
	for _, dir := range dirs {
		templates.AddImportRegex(fmt.Sprintf(`(^|[^A-Za-z])%s\.`, dir), fmt.Sprintf("%s/%s", projectModule, dir), "")
	}

	curPath, err := os.Getwd()
	if err != nil {
		return err
	}
	currentDir = curPath

	// Define init steps with names
	initSteps := []initStep{
		{"Creating directories", initDir},
		{"Creating .env file", initEnv},
		{"Creating main.go", initMain},
		{"Creating commands", initCmd},
		{"Creating gqlgen.yml", initGqlgen},
		{"Creating GraphQL schema", initModels},
		{"Creating resolvers", initResolvers},
		{"Creating server", initServer},
		{"Creating go.mod", initMod},
		{"Creating config", initConfig},
		{"Creating app-start command", initAppStart},
		{"Creating .gitignore", initGitignore},
		{"Creating root schema", initGraphql},
		{"Creating migration command", initMigration},
		{"Creating schema command", initSchemaCmd},
		{"Creating Atlas loader", initLoader},
		{"Creating atlas.hcl", initAtlas},
	}

	// Register import patterns
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
	templates.AddImportRegex("messaging", "github.com/light-speak/lighthouse/messaging", "")
	templates.AddImportRegex("queue", "github.com/light-speak/lighthouse/queue", "")
	templates.AddImportRegex("redis", "github.com/light-speak/lighthouse/redis", "")
	templates.AddImportRegex("bytes", "bytes", "")

	// Execute init steps with progress
	totalSteps := len(initSteps)
	for i, step := range initSteps {
		progress := (i * 100) / totalSteps
		nextProgress := ((i + 1) * 100) / totalSteps
		utils.SmoothProgress(progress, nextProgress, step.name, 50*time.Millisecond, false)

		if err := step.fn(); err != nil {
			fmt.Println() // New line after progress bar
			logs.Error().Msgf("✗ Failed: %s - %v", step.name, err)
			return err
		}
	}
	utils.SmoothProgress(100, 100, "Files created", 50*time.Millisecond, true)
	fmt.Println() // New line after progress bar
	fmt.Println()

	// Change to project directory
	err = os.Chdir(filepath.Join(currentDir, projectName))
	if err != nil {
		return err
	}

	// Run go mod tidy
	logs.Info().Msg("📦 Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logs.Error().Msgf("Failed to run go mod tidy: %v", err)
		return err
	}

	// Generate schema
	logs.Info().Msg("⚡ Generating GraphQL schema...")
	cmd = exec.Command("lightcmd", "generate:schema")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logs.Error().Msgf("Failed to generate schema: %v", err)
		return err
	}

	// Success message
	fmt.Println()
	logs.Info().Msg("✅ Project initialized successfully!")
	fmt.Println()
	logs.Info().Msgf("   cd %s", projectName)
	logs.Info().Msg("   go run main.go app:start")
	fmt.Println()

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
	// Check if go.mod exists in parent directories (up to 3 levels)
	if hasGoModInParents(currentDir, 3) {
		logs.Info().Msg("Found go.mod in parent directory, skipping go.mod creation")
		return nil
	}

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

// hasGoModInParents checks if go.mod exists in the current directory or any parent directory up to maxLevels
func hasGoModInParents(dir string, maxLevels int) bool {
	for i := 0; i <= maxLevels; i++ {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}
	return false
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

func initMigration() error {
	migrationTpl, err := tpl.ReadFile("tpl/migration.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "commands"),
		Template:     string(migrationTpl),
		FileName:     "migration",
		Package:      "commands",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
	}
	templates.AddImportRegex("atlasexec", "ariga.io/atlas-go-sdk/atlasexec", "")
	templates.AddImportRegex("context", "context", "")
	templates.AddImportRegex("formatter", "github.com/vektah/gqlparser/v2/formatter", "")

	return templates.Render(options)
}

func initSchemaCmd() error {
	schemaCmdTpl, err := tpl.ReadFile("tpl/schema_cmd.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "commands"),
		Template:     string(schemaCmdTpl),
		FileName:     "schema",
		Package:      "commands",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
	}
	return templates.Render(options)
}

func initLoader() error {
	loaderTpl, err := tpl.ReadFile("tpl/loader.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName, "loader"),
		Template:     string(loaderTpl),
		FileName:     "main",
		FileExt:      "go",
		Editable:     true,
		SkipIfExists: true,
		SkipImport:   false,
	}
	templates.AddImportRegex("gormschema", "ariga.io/atlas-provider-gorm/gormschema", "")
	return templates.Render(options)
}

func initAtlas() error {
	atlasTpl, err := tpl.ReadFile("tpl/atlas.tpl")
	if err != nil {
		return err
	}
	options := &templates.Options{
		Path:         filepath.Join(projectName),
		Template:     string(atlasTpl),
		FileName:     "atlas",
		FileExt:      "hcl",
		Editable:     true,
		SkipIfExists: true,
		Data: map[string]string{
			"ProjectName": projectName,
		},
	}
	return templates.Render(options)
}
