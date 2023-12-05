//go:build ignore

package main

import (
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/vektah/gqlparser/v2/ast"
	"go/types"
	"os"
	"strings"
)

func mutateHook(b *modelgen.ModelBuild) *modelgen.ModelBuild {
	for _, model := range b.Models {

		ftype := buildNamedType("github.com/light-speak/lighthouse/db.Model")
		model.Fields = append(model.Fields, &modelgen.Field{
			Description: "Custom Gorm Model",
			Name:        "",
			GoName:      "",
			Type:        ftype,
			Tag:         "",
			Omittable:   false,
		})

		for index, field := range model.Fields {
			if field.Name == "id" {
				model.Fields = append(model.Fields[:index], model.Fields[index+1:]...)
			}
		}
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

func main() {
	cfg, err := config.LoadConfigFromDefaultLocations()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load config", err.Error())
		os.Exit(2)
	}

	// Attaching the mutation function onto modelgen plugin
	p := modelgen.Plugin{
		FieldHook:  constraintFieldHook,
		MutateHook: mutateHook,
	}
	err = api.Generate(cfg, api.ReplacePlugin(&p))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
