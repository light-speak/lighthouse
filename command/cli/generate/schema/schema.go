// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package schema

import (
	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/graphql"
)

type Schema struct{}

func (c *Schema) Name() string {
	// Func:Name user code start. Do not remove this comment.
	return "generate:schema"
	// Func:Name user code end. Do not remove this comment.
}

func (c *Schema) Usage() string {
	// Func:Usage user code start. Do not remove this comment.
	return "This is a command generated by lighthouse cli"
	// Func:Usage user code end. Do not remove this comment.
}

func (c *Schema) Args() []*command.CommandArg {
	return []*command.CommandArg{
		// Func:Args user code start. Do not remove this comment.
		// Func:Args user code end. Do not remove this comment.
	}
}

func (c *Schema) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		// Func:Action user code start. Do not remove this comment.
		err := graphql.Generate()
		if err != nil {
			return err
		}
		// Func:Action user code end. Do not remove this comment.
		return nil
	}
}

func (c *Schema) OnExit() func() {
	return func() {}
}

// Section: user code section start. Do not remove this comment.
// Section: user code section end. Do not remove this comment.
