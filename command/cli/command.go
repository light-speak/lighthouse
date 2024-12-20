package cli

import (
	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/command/cli/generate/cmd"
	"github.com/light-speak/lighthouse/command/cli/generate/initialize"
	"github.com/light-speak/lighthouse/command/cli/generate/schema"
	"github.com/light-speak/lighthouse/command/cli/version"
)

type Lighthouse struct{}

func (c *Lighthouse) GetCommands() []command.Command {
	return []command.Command{
		&cmd.GenCmd{},
		&version.Version{},
		&initialize.Initialize{},
		&schema.Schema{},
	}
}
