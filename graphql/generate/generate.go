package generate

import (
	"fmt"
	"go/types"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/light-speak/lighthouse/db"
	"github.com/light-speak/lighthouse/graphql/generate/dataloader"
	"github.com/light-speak/lighthouse/graphql/generate/merge"
	"github.com/light-speak/lighthouse/graphql/generate/resolver"
	"github.com/light-speak/lighthouse/graphql/generate/scope"
	"github.com/light-speak/lighthouse/graphql/generate/searchable"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

var customModelName map[string]int

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	filteredModels := make([]*modelgen.Object, 0, len(b.Models))
	dataloaderModels := make([]*modelgen.Object, 0)

	for _, model := range b.Models {
		if shouldExcludeModel(model.Name) {
			continue
		}

		var ftype types.Type
		if len(model.Implements) > 0 {
			switch model.Implements[0] {
			case "BaseModel":
				ftype = buildNamedType("github.com/light-speak/lighthouse/db.Model")
			case "BaseModelSoftDelete":
				ftype = buildNamedType("github.com/light-speak/lighthouse/db.ModelSoftDelete")
			}
		}

		if ftype != nil {
			model.Fields = append(model.Fields, &modelgen.Field{
				Description: "Custom GORM Model",
				Type:        ftype,
				Tag:         `gorm:"embedded"`,
			})

			// 移除 "id" 字段
			for i, field := range model.Fields {
				if field.Name == "id" {
					model.Fields = append(model.Fields[:i], model.Fields[i+1:]...)
					break
				}
			}

			dataloaderModels = append(dataloaderModels, model)
			customModelName[model.Name] = 1
		}
		filteredModels = append(filteredModels, model)
	}

	b.Models = filteredModels

	if err := dataloader.GenAllDataloader(dataloaderModels); err != nil {
		log.Error("生成 dataloader 时出错: %v", err)
	}

	return b
}
func constraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	f, err := modelgen.DefaultFieldMutateHook(td, fd, f)
	if err != nil {
		return nil, err
	}

	if c := fd.Directives.ForName("constraint"); c != nil {
		if formatConstraint := c.Arguments.ForName("format"); formatConstraint != nil {
			f.Tag += " validate:" + formatConstraint.Value.String()
		}
	}

	return f, nil
}

// shouldExcludeModel 判断是否应该排除某个模型
func shouldExcludeModel(name string) bool {
	excludedTypes := map[string]bool{
		"Query":        true,
		"Mutation":     true,
		"Subscription": true,
	}

	return excludedTypes[name]
}
func buildNamedType(fullName string) types.Type {
	dotIndex := strings.LastIndex(fullName, ".")
	pkgPath, typeName := fullName[:dotIndex], fullName[dotIndex+1:]

	pkgName := pkgPath
	if slashIndex := strings.LastIndex(pkgPath, "/"); slashIndex != -1 {
		pkgName = pkgPath[slashIndex+1:]
	}

	pkg := types.NewPackage(pkgPath, pkgName)
	return types.NewNamed(types.NewTypeName(0, pkg, typeName, nil), nil, nil)
}

// 从字段返回类型中提取模型名称
func getModelNameFromField(field *ast.FieldDefinition) string {
	typeName := field.Type.Name()

	for field.Type.Elem != nil {
		field.Type = field.Type.Elem
		typeName = field.Type.Name()
	}

	return typeName
}
func generateDirectives(cfg *config.Config) {
	for _, d := range cfg.Schema.Query.Fields {
		modelName := getModelNameFromField(d)
		for _, directive := range d.Directives {
			for _, arg := range directive.Arguments {
				if arg.Name == "model" {
					modelName = arg.Value.Raw
					if err := scope.ModelScope(modelName); err != nil {
						log.Error("%+v", err)
					}
				}
				if arg.Name == "scopes" && arg.Value.Kind == ast.ListValue {
					for _, scopeValue := range arg.Value.Children {
						if err := scope.Generate(strings.ToLower(modelName), scopeValue.Value.Raw); err != nil {
							log.Error("%+v", err)
						}
					}
				}
			}
		}
	}

	var mergeTypes []*merge.MergeType

	for _, t := range cfg.Schema.Types {
		if t.Kind != ast.Object || strings.HasPrefix(t.Name, "_") || t.Name == "Entity" || t.Name == "Query" || t.Name == "Mutation" || t.Name == "Subscription" {
			continue
		}

		var mergeFields []*merge.MergeField
		var searchableModelOptions []*searchable.SearchableModelOption
		for _, f := range t.Fields {
			if requires := f.Directives.ForName("requires"); requires != nil {
				if field := requires.Arguments.ForName("fields"); field != nil {
					local := "True"
					if v, ok := customModelName[utils.UcFirst(f.Name)]; ok && v > 0 {
						local = "False"
					}
					mergeFields = append(mergeFields, &merge.MergeField{
						Target: f.Name,
						Source: field.Value.Raw,
						Local:  local,
					})
				}
			}
			var searchableFieldOption *searchable.SearchableFieldOption
			if searchableDirective := f.Directives.ForName("searchable"); searchableDirective != nil {
				if field := searchableDirective.Arguments.ForName("searchableType"); field != nil {
					searchableFieldOption = &searchable.SearchableFieldOption{
						FieldName:      f.Name,
						SearchableType: field.Value.Raw,
					}
					searchableFieldOption.IndexAnalyzer = "IK_MAX_WORD"
					searchableFieldOption.SearchAnalyzer = "IK_SMART"
					if field := searchableDirective.Arguments.ForName("indexAnalyzer"); field != nil {
						searchableFieldOption.IndexAnalyzer = field.Value.Raw
					}
					if field := searchableDirective.Arguments.ForName("searchAnalyzer"); field != nil {
						searchableFieldOption.SearchAnalyzer = field.Value.Raw
					}
				}

				if len(searchableModelOptions) == 0 || searchableModelOptions[len(searchableModelOptions)-1].ModelName != t.Name {
					searchableModelOptions = append(searchableModelOptions, &searchable.SearchableModelOption{
						ModelName: t.Name,
						Fields:    []*searchable.SearchableFieldOption{searchableFieldOption},
					})
				} else {
					searchableModelOptions[len(searchableModelOptions)-1].Fields = append(searchableModelOptions[len(searchableModelOptions)-1].Fields, searchableFieldOption)
				}
			}
		}

		if err := searchable.GenSearchableModels(searchableModelOptions); err != nil {
			log.Error("%+v", err)
		}

		mergeTypes = append(mergeTypes, &merge.MergeType{
			Model:      t.Name,
			MergeField: mergeFields,
		})
	}

	if err := merge.GenMergeModels(mergeTypes); err != nil {
		log.Error("%+v", err)
	}
}

// getLibraryPath 获取当前库的根目录路径
func getLibraryPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("获取当前文件路径失败")
	}

	// 返回 lighthouse 目录
	return filepath.Dir(filepath.Dir(currentFilePath)), nil
}

func Run() error {
	customModelName = make(map[string]int)
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	currentDir, err := getLibraryPath()
	if err != nil {
		return fmt.Errorf("获取当前目录失败: %w", err)
	}

	cfg.Model.ModelTemplate = filepath.Join(currentDir, "generate", "tpl", "models.gotpl")
	cfg.Resolver.ResolverTemplate = filepath.Join(currentDir, "generate", "tpl", "resolver.gotpl")

	if err := db.Init(); err != nil {
		return fmt.Errorf("初始化数据库失败: %w", err)
	}

	p := modelgen.Plugin{
		FieldHook:  constraintFieldHook,
		MutateHook: mutateHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))
	if err != nil {
		log.Warn("生成环节出现错误，已自动忽略，部分Func将在后续过程中生成，请检查生成文件: %v", err)
	}
	generateDirectives(cfg)
	resolver.GenerateResolver()

	return nil
}
