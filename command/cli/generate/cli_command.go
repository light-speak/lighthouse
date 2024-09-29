package generate

import (
	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/command/cli/generate/cmd"
)

type Lighthouse struct{}

func (c *Lighthouse) GetCommands() []command.Command {
	return []command.Command{
		&cmd.GenCmd{},
	}
}
