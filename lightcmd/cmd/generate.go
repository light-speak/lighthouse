package cmd

import (
	"github.com/light-speak/lighthouse/lightcmd/generate"
	_ "github.com/light-speak/lighthouse/lightcmd/generate/directives"
)

type Generate struct{}

func (c *Generate) Name() string {
	return "generate:schema"
}

func (c *Generate) Usage() string {
	return "This is a command to generate the service"
}

func (c *Generate) Args() []*CommandArg {
	return []*CommandArg{}
}

func (c *Generate) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		err := generate.GenerateSchema()
		if err != nil {
			return err
		}
		return nil
	}
}

func (c *Generate) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&Generate{})
}
