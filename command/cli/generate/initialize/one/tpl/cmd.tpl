type Command struct{
  {{ section }}
}

func (c *Command) GetCommands() []command.Command {
	return []command.Command{
    {{ funcStart "GetCommands" }}
    {{ funcEnd "GetCommands" }}
	}
}
