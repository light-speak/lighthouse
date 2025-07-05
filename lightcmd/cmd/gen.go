package cmd

import (
	"embed"
	"path/filepath"

	"github.com/light-speak/lighthouse/templates"
)

//go:embed tpl
var tpl embed.FS

type GenCmd struct{}

func (c *GenCmd) Name() string {
	return "generate:command"
}

func (c *GenCmd) Usage() string {
	return "Generate a new command"
}

func (c *GenCmd) Args() []*CommandArg {
	return []*CommandArg{
		{
			Name:     "scope",
			Usage:    "The scope of the command",
			Required: false,
			Default:  "app",
			Type:     String,
		},
		{
			Name:     "name",
			Usage:    "The name of the command",
			Required: true,
			Type:     String,
		},
	}
}

func (c *GenCmd) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		args, err := GetArgs(c.Args(), flagValues)
		if err != nil {
			return err
		}

		var scope, name *string
		if scope, err = GetStringArg(args, "scope"); err != nil {
			return err
		}
		if name, err = GetStringArg(args, "name"); err != nil {
			return err
		}

		tpl, err := tpl.ReadFile("tpl/gen.tpl")
		if err != nil {
			return err
		}

		options := templates.Options{
			Path:     filepath.Join("commands"),
			Template: string(tpl),
			FileName: *name,
			Editable: true,
			FileExt:  "go",
			Data: map[string]interface{}{
				"Scope": *scope,
				"Name":  *name,
			},
		}
		templates.AddImportRegex(`logs\.`, "github.com/light-speak/lighthouse/logs", "")
		templates.AddImportRegex(`cmd\.`, "github.com/light-speak/lighthouse/lightcmd/cmd", "")

		if err := templates.Render(&options); err != nil {
			return err
		}

		return nil
	}
}

func (c *GenCmd) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&GenCmd{})
}
