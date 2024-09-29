package cmd

import (
	"fmt"

	"github.com/light-speak/lighthouse/command"
)

type GenCmd struct {
}

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
		},
		{
			Name:     "name",
			Usage:    "The name of the command",
			Required: true,
		},
		{
			Name:     "path",
			Usage:    "The path of the command",
			Required: false,
			Default:  "cmd",
		},
	}
}

func (c *GenCmd) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		args, err := command.GetArgs(c.Args(), flagValues)
		if err != nil {
			return err
		}

		var scope, name string
		if scope, err = command.GetStringArg(args); err != nil {
			return err
		}
		if name, err = command.GetStringArg(args); err != nil {
			return err
		}
	

		fmt.Printf("Generating command: %s in scope: %s\n", name, scope)
		return nil
	}
}
