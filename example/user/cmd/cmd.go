// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package cmd

import (
	"user/cmd/start"

	"github.com/light-speak/lighthouse/command"
)

type Command struct {
	// Section: user code section start. Do not remove this comment.

	// Section: user code section end. Do not remove this comment.

}

func (c *Command) GetCommands() []command.Command {
	return []command.Command{
		// Func:GetCommands user code start. Do not remove this comment.
		&start.Start{},
		// Func:GetCommands user code end. Do not remove this comment.
	}
}
