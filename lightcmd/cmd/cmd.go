package cmd

type Command struct{}

var commands []CommandInterface

func AddCommand(command CommandInterface) {
	commands = append(commands, command)
}

func (c *Command) GetCommands() []CommandInterface {
	return commands
}
