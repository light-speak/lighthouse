package command

import "fmt"

type LighthouseCommand struct{}

func (c *LighthouseCommand) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		name := *flagValues["name"].(*string)
		if name == "" {
			name = "world"
		}
		fmt.Printf("hello, %s!\n", name)
		return nil
	}
}

func (c *LighthouseCommand) Name() string {
	return "hello"
}

func (c *LighthouseCommand) Usage() string {
	return "say hello"
}

func (c *LighthouseCommand) Args() []*CommandArg {
	return []*CommandArg{
		{
			Name:     "name",
			Type:     String,
			Usage:    "name to say hello to",
			Required: false,
		},
	}
}

type LighthouseCommandList struct{}

func (c *LighthouseCommandList) GetCommands() []Command {
	return []Command{
		&LighthouseCommand{},
	}
}
