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
	"github.com/light-speak/lighthouse/graphql/generate/scope"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
	"github.com/vektah/gqlparser/v2/ast"
)

var customModelName map[string]int

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	var filteredModels []*modelgen.Object   // 用于存储过滤后的模型
	var dataloaderModels []*modelgen.Object // 用于存储过滤后的模型

	for _, model := range b.Models {
		// 检查并过滤掉不需要附加 GORM Model 的类型
		if shouldExcludeModel(model.Name) {
			// 从 b.Models 中移除该类型，不添加到 filteredModels
			continue
		}
		var ftype types.Type = nil
		if len(model.Implements) > 0 {
			switch model.Implements[0] {
			case "BaseModel":
				ftype = buildNamedType("gitlab.staticoft.com/lighthouse/db.Model")
				break
			case "BaseModelSoftDelete":
				ftype = buildNamedType("gitlab.staticoft.com/lighthouse/db.ModelSoftDelete")
				break
			}
		}

		if ftype != nil {
			model.Fields = append(model.Fields, &modelgen.Field{
				Description: "Custom GORM Model",
				Name:        "",
				GoName:      "",
				Type:        ftype,
				Tag:         `gorm:"embedded"`, // GORM 中常用的嵌入式模型
				Omittable:   false,
			})

			// 移除字段名为 "id" 的字段
			for index, field := range model.Fields {
				if field.Name == "id" {
					model.Fields = append(model.Fields[:index], model.Fields[index+1:]...)
				}
			}
			// 将处理后的模型添加到 filteredModels 列表中
			dataloaderModels = append(dataloaderModels, model)
			customModelName[model.Name] = 1
		}
		filteredModels = append(filteredModels, model)
	}

	// 更新 b.Models 列表为过滤后的模型列表
	b.Models = filteredModels

	err := dataloader.GenModelLoader(dataloaderModels)
	if err != nil {
		log.Error("%s", err)
	}

	return b
}

func constraintFieldHook(td *ast.Definition, fd *ast.FieldDefinition, f *modelgen.Field) (*modelgen.Field, error) {
	if f, err := modelgen.DefaultFieldMutateHook(td, fd, f); err != nil {
		return f, err
	}

	c := fd.Directives.ForName("constraint")
	if c != nil {
		formatConstraint := c.Arguments.ForName("format")
		if formatConstraint != nil {
			f.Tag += " validate:" + formatConstraint.Value.String()
		}
	}

	return f, nil
}

// shouldExcludeModel 确定是否排除模型
func shouldExcludeModel(name string) bool {
	// 根据名称排除特定类型
	excludedTypes := map[string]bool{
		"Query":        true,
		"Mutation":     true,
		"Subscription": true,
	}

	// 例如，这里可以扩展更多排除规则
	if _, exists := excludedTypes[name]; exists {
		return true
	}

	return false
}

func buildNamedType(fullName string) types.Type {
	dotIndex := strings.LastIndex(fullName, ".")
	// type is pkg.Name
	pkgPath := fullName[:dotIndex]
	typeName := fullName[dotIndex+1:]

	pkgName := pkgPath
	slashIndex := strings.LastIndex(pkgPath, "/")
	if slashIndex != -1 {
		pkgName = pkgPath[slashIndex+1:]
	}

	pkg := types.NewPackage(pkgPath, pkgName)
	return types.NewNamed(types.NewTypeName(0, pkg, typeName, nil), nil, nil)
}

// 从字段返回类型中提取模型名称
func getModelNameFromField(field *ast.FieldDefinition) string {
	// 获取字段返回类型的名称
	typeName := field.Type.Name()

	// 如果类型是非空或列表类型，继续获取基础类型名称
	for field.Type.Elem != nil {
		typeName = field.Type.Elem.Name()
		field.Type = field.Type.Elem
	}

	return typeName
}

func generateDirectives(cfg *config.Config) {
	for _, d := range cfg.Schema.Query.Fields {
		modelName := getModelNameFromField(d)
		// 遍历每个字段上的所有 Directives
		for _, directive := range d.Directives {
			// 解析 directive 的 Arguments
			for _, arg := range directive.Arguments {
				// 如果参数有Model，就用这个当ModelName
				if arg.Name == "model" {
					modelName = arg.Value.Raw
					err := scope.ModelScope(modelName)
					if err != nil {
						log.Error("%+v", err)
					}
				}
				if arg.Name == "scopes" && arg.Value.Kind == ast.ListValue {
					// 如果是 scopes 参数，并且是一个列表
					for _, scopeValue := range arg.Value.Children {
						s := scopeValue.Value // 提取每个 scope 的值
						//fmt.Printf("Found scope: %s for field: %s\n", s, d.Name)
						err := scope.Generate(strings.ToLower(modelName), s.Raw)
						if err != nil {
							log.Error("%+v", err)
						}
					}
				}
			}

		}
	}

	var mergeTypes []*merge.MergeType

	for _, t := range cfg.Schema.Types {
		var mergeFields []*merge.MergeField

		for _, f := range t.Fields {
			requires := f.Directives.ForName("requires")
			if requires != nil {
				field := requires.Arguments.ForName("fields")
				target := f.Name
				source := field.Value.Raw
				local := "True"
				if v, ok := customModelName[utils.UcFirst(target)]; ok && v > 0 {
					local = "False"
				}
				mergeFields = append(mergeFields, &merge.MergeField{
					Target: target,
					Source: source,
					Local:  local,
				})
			}
		}

		if t.Kind == ast.Object && !strings.HasPrefix(t.Name, "_") && t.Name != "Entity" && t.Name != "Query" && t.Name != "Mutation" && t.Name != "Subscription" {
			// 判断是否是本地模型
			mergeTypes = append(mergeTypes, &merge.MergeType{
				Model:      t.Name,
				MergeField: mergeFields,
			})
		}
	}

	err := merge.GenMergeModels(mergeTypes)
	if err != nil {
		log.Error("%+v", err)
	}
}

// getLibraryPath 获取当前库的根目录路径
func getLibraryPath() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}

	// b 是当前文件的路径, 返回 lighthouse 目录
	return filepath.Dir(filepath.Dir(currentFilePath)), nil
}

func Run() error {
	customModelName = make(map[string]int)
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		return err
	}

	currentDir, err := getLibraryPath()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	cfg.Model.ModelTemplate = filepath.Join(currentDir, "generate", "tpl", "models.gotpl")
	cfg.Resolver.ResolverTemplate = filepath.Join(currentDir, "generate", "tpl", "resolver.gotpl")

	err = db.Init()
	if err != nil {
		return err
	}

	p := modelgen.Plugin{
		FieldHook:  constraintFieldHook,
		MutateHook: mutateHook,
	}

	err = api.Generate(cfg, api.ReplacePlugin(&p))
	generateDirectives(cfg)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
