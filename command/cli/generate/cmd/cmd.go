package cmd

import (
	_ "embed"
	"path/filepath"

	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/template"
)

//go:embed cmd.tpl
var temp string

type GenCmd struct{}

func (c *GenCmd) Name() string {
	return "generate:command"
}

func (c *GenCmd) Usage() string {
	return "Generate a new command"
}

func (c *GenCmd) Args() []*command.CommandArg {
	return []*command.CommandArg{
		{
			Name:     "scope",
			Usage:    "The scope of the command",
			Required: false,
			Default:  "app",
			Type:     command.String,
		},
		{
			Name:     "name",
			Usage:    "The name of the command",
			Required: true,
			Type:     command.String,
		},
		{
			Name:     "path",
			Usage:    "The path of the command",
			Required: false,
			Default:  "cmd",
			Type:     command.String,
		},
	}
}

func (c *GenCmd) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		args, err := command.GetArgs(c.Args(), flagValues)
		if err != nil {
			return err
		}

		var scope, name, path *string
		if scope, err = command.GetStringArg(args, "scope"); err != nil {
			return err
		}
		if name, err = command.GetStringArg(args, "name"); err != nil {
			return err
		}
		if path, err = command.GetStringArg(args, "path"); err != nil {
			return err
		}

		options := template.Options{
			Path:     filepath.Clean(filepath.Join(*path, *name)),
			Template: temp,
			FileName: *name,
			Editable: true,
			FileExt:  "go",
			Data: map[string]interface{}{
				"Scope": *scope,
				"Name":  *name,
			},
		}

		if err := template.Render(&options); err != nil {
			return err
		}

		return nil
	}
}

func (c *GenCmd) OnExit() func() {
	return func() {}
}
