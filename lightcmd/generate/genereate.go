package generate

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/templates"
	"github.com/vektah/gqlparser/v2/ast"
)

//go:embed tpl
var tpl embed.FS

// modelDepends maps model name to its dependencies
var modelDepends = make(map[string][]*ModelDependency)

var directives = map[string]func(*ast.Directive, *DirectiveLogic) (*DirectiveLogic, error){}

func AddDirective(name string, fn func(*ast.Directive, *DirectiveLogic) (*DirectiveLogic, error)) {
	directives[name] = fn
}

func GenerateSchema() error {
	logs.Info().Msg("Starting schema generation...")

	// Load config
	logs.Info().Msg("Loading gqlgen config...")
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		logs.Error().Msgf("Failed to load gqlgen config: %v", err)
		return fmt.Errorf("failed to load gqlgen config: %w", err)
	}
	logs.Info().Msgf("Config loaded from: %s", cfg.SchemaFilename)

	// Generate schema
	logs.Info().Msg("Generating GraphQL schema and models...")
	p := &modelgen.Plugin{
		FieldHook: fieldHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(p))
	if err != nil {
		logs.Error().Msgf("Failed to generate schema: %+v", err)
		return fmt.Errorf("failed to generate schema: %w", err)
	}

	// Print generated models info
	if len(loaderTypeToFieldsMap) > 0 {
		logs.Info().Msgf("Generated %d models with @loader directive:", len(loaderTypeToFieldsMap))
		for modelName, keys := range loaderTypeToFieldsMap {
			if len(keys) > 0 {
				logs.Info().Msgf("  - %s (keys: %s)", modelName, strings.Join(keys, ", "))
			} else {
				logs.Info().Msgf("  - %s", modelName)
			}
		}
	}

	// Generate dataloader
	logs.Info().Msg("Generating DataLoader...")
	err = generateLoader()
	if err != nil {
		logs.Error().Msgf("Failed to generate DataLoader: %v", err)
		return fmt.Errorf("failed to generate dataloader: %w", err)
	}
	logs.Info().Msg("DataLoader generated successfully")

	// Run go mod tidy
	logs.Info().Msg("Running go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logs.Error().Msgf("Failed to run go mod tidy: %v", err)
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	// Run gofmt
	logs.Info().Msg("Running gofmt...")
	cmd = exec.Command("gofmt", "-s", "-w", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		logs.Error().Msgf("Failed to run gofmt: %v", err)
		return fmt.Errorf("failed to run gofmt: %w", err)
	}

	logs.Info().Msg("Schema generation completed successfully!")
	return nil
}

func generateLoader() error {
	dataloaderTpl, err := tpl.ReadFile("tpl/dataloader.tpl")
	if err != nil {
		return fmt.Errorf("failed to read dataloader template: %w", err)
	}

	curPath, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	outputPath := filepath.Join(curPath, "models", "dataloader_gen.go")
	logs.Info().Msgf("Writing DataLoader to: %s", outputPath)

	// Log MorphTo loaders if any
	if len(loaderTypeMorphToMap) > 0 {
		logs.Info().Msgf("Generated %d MorphTo loaders:", len(loaderTypeMorphToMap))
		for modelName, fields := range loaderTypeMorphToMap {
			for _, field := range fields {
				logs.Info().Msgf("  - %s.%s -> [%s]", modelName, field.Field, strings.Join(field.Union, ", "))
			}
		}
	}

	// Log extra keys if any
	if len(loaderTypeExtraKeysMap) > 0 {
		for modelName, keys := range loaderTypeExtraKeysMap {
			if len(keys) > 0 {
				logs.Info().Msgf("  - %s extra keys: %s", modelName, strings.Join(keys, ", "))
			}
		}
	}

	options := &templates.Options{
		Path:         filepath.Join(curPath, "models"),
		Template:     string(dataloaderTpl),
		FileName:     "dataloader_gen",
		Package:      "models",
		FileExt:      "go",
		Editable:     false,
		SkipIfExists: false,
		Data: map[string]any{
			"Loader":  loaderTypeToFieldsMap,
			"MorphTo": loaderTypeMorphToMap,
			"Extra":   loaderTypeExtraKeysMap,
		},
	}
	templates.AddImportRegex("dataloadgen", "github.com/vikstrous/dataloadgen", "")
	templates.AddImportRegex("lighterr", "github.com/light-speak/lighthouse/lighterr", "")
	templates.AddImportRegex("context", "context", "")
	templates.AddImportRegex("databases", "github.com/light-speak/lighthouse/databases", "")
	templates.AddImportRegex("dataloader", "github.com/light-speak/lighthouse/routers/dataloader", "")
	templates.AddImportRegex("gorm", "gorm.io/gorm", "")

	err = templates.Render(options)
	if err != nil {
		return fmt.Errorf("failed to render dataloader template: %w", err)
	}
	return nil
}

type LoaderField struct {
	Field string
	Union []string
}

var loaderTypeMorphToMap = make(map[string][]*LoaderField)
var loaderTypeToFieldsMap = make(map[string][]string)
var loaderTypeExtraKeysMap = make(map[string][]string)

func fieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	if loaderDirective := td.Directives.ForName("loader"); loaderDirective != nil {
		if _, exists := loaderTypeToFieldsMap[td.Name]; !exists {
			loaderTypeToFieldsMap[td.Name] = make([]string, 0)
			keys := loaderDirective.Arguments.ForName("keys")
			if keys != nil {
				for _, key := range keys.Value.Children {
					loaderTypeToFieldsMap[td.Name] = append(loaderTypeToFieldsMap[td.Name], key.Value.Raw)
				}
			}
		}
		if _, exists := loaderTypeMorphToMap[td.Name]; !exists {
			loaderTypeMorphToMap[td.Name] = make([]*LoaderField, 0)
			morphKey := loaderDirective.Arguments.ForName("morphKey")
			unionTypes := loaderDirective.Arguments.ForName("unionTypes")
			if morphKey != nil && unionTypes != nil {
				for _, key := range unionTypes.Value.Children {
					loaderTypeMorphToMap[td.Name] = append(loaderTypeMorphToMap[td.Name], &LoaderField{
						Field: morphKey.Value.Raw,
						Union: []string{key.Value.Raw},
					})
				}
			}
		}
		if _, exists := loaderTypeExtraKeysMap[td.Name]; !exists {
			loaderTypeExtraKeysMap[td.Name] = make([]string, 0)
			extraKeys := loaderDirective.Arguments.ForName("extraKeys")
			if extraKeys != nil {
				for _, key := range extraKeys.Value.Children {
					loaderTypeExtraKeysMap[td.Name] = append(loaderTypeExtraKeysMap[td.Name], key.Value.Raw)
				}
			}
		}
	}
	tag, err := fieldLogic(td, fd, f)
	if err != nil {
		return nil, err
	}
	if tag != "" {
		f.Tag = tag
	}

	f, err = modelgen.DefaultFieldMutateHook(td, fd, f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type DirectiveLogic struct {
	TagKvs       map[string][]string
	ModelDepends []*DependField
	ModelWith    bool
}

// DependField represents a field dependency
type DependField struct {
	Field    string `json:"field"`    // Source field name
	Value    string `json:"value"`    // Value to be assigned
	Type     string `json:"type"`     // Type of dependency (FIELD or VALUE)
	Nullable bool   `json:"nullable"` // Whether the dependency field is nullable
}

// ModelDependency represents model field dependencies
type ModelDependency struct {
	Field       string         `json:"field"`        // Field name in current model
	DependField string         `json:"depend_field"` // Referenced field name
	Depends     []*DependField `json:"depends"`      // List of dependencies
}

func fieldLogic(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (string, error) {
	tag := f.Tag
	tagKvs := map[string][]string{}

	// Initialize modelDepends for this model
	modelName := td.Name
	if _, exists := modelDepends[modelName]; !exists {
		modelDepends[modelName] = make([]*ModelDependency, 0)
	}

	switch fd.Name {
	case "id":
		tagKvs["gorm"] = append(tagKvs["gorm"], `type:int unsigned`)
		tagKvs["gorm"] = append(tagKvs["gorm"], `primary_key`)
		tagKvs["gorm"] = append(tagKvs["gorm"], `auto_increment`)
	case "createdAt":
		tagKvs["gorm"] = append(tagKvs["gorm"], `type:datetime`)
	case "updatedAt":
		tagKvs["gorm"] = append(tagKvs["gorm"], `type:datetime`)
	case "deletedAt":
		tagKvs["gorm"] = append(tagKvs["gorm"], `type:datetime`)
		tagKvs["gorm"] = append(tagKvs["gorm"], `index`)
	}
	for dName, fn := range directives {
		directives := fd.Directives.ForNames(dName)
		for _, directive := range directives {
			logic, err := fn(directive, &DirectiveLogic{
				TagKvs: tagKvs,
			})
			if err != nil {
				return "", err
			}
			for k, v := range logic.TagKvs {
				tagKvs[k] = append(tagKvs[k], v...)
			}
			if len(logic.ModelDepends) > 0 {
				for _, dep := range logic.ModelDepends {
					dep.Nullable = !fd.Type.NonNull
				}
				// Check if dependency already exists
				exists := false
				for _, dep := range modelDepends[modelName] {
					if dep.Field == fd.Name {
						exists = true
						// Merge dependencies
						dep.Depends = append(dep.Depends, logic.ModelDepends...)
						break
					}
				}
				if !exists {
					modelDepends[modelName] = append(modelDepends[modelName], &ModelDependency{
						Field:       fd.Name,
						DependField: fd.Type.NamedType,
						Depends:     logic.ModelDepends,
					})
				}
			}
		}
	}

	hasTypeTag := false
	for _, v := range tagKvs["gorm"] {
		if strings.HasPrefix(v, "type:") {
			hasTypeTag = true
			break
		}
	}
	if isStringType(f) && !hasTypeTag {
		tagKvs["gorm"] = append(tagKvs["gorm"], "type:varchar(255)")
	}
	if fd.Type.NonNull {
		tagKvs["gorm"] = append(tagKvs["gorm"], "not null")
	}
	if fd.Description != "" {
		tagKvs["gorm"] = append(tagKvs["gorm"], fmt.Sprintf(`comment:%s`, fd.Description))
	}
	tagKvsStr := ""
	tagKvsStr += tag
	for k, v := range tagKvs {
		// Deduplicate values before joining
		uniqueVals := make([]string, 0)
		seen := make(map[string]bool)
		for _, val := range v {
			if !seen[val] {
				seen[val] = true
				uniqueVals = append(uniqueVals, val)
			}
		}
		tagKvsStr += fmt.Sprintf(" %s:\"%s\" ", k, strings.Join(uniqueVals, ";"))
	}
	return tagKvsStr, nil
}

func isStringType(field *modelgen.Field) bool {
	if field.Type.String() == "string" || field.Type.String() == "*string" {
		return true
	}
	return false
}
