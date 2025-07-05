type Command struct {}

var commands []cmd.CommandInterface

func AddCommand(command cmd.CommandInterface) {
	commands = append(commands, command)
}

func (c *Command) GetCommands() []cmd.CommandInterface {
	return commands
}
