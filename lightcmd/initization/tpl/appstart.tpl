type Start struct{}

func (c *Start) Name() string {
	return "app:start"
}

func (c *Start) Usage() string {
	return "This is a command to Start the service"
}

func (c *Start) Args() []*cmd.CommandArg {
	return []*cmd.CommandArg{}
}

func (c *Start) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		server.StartService()
		return nil
	}
}

func (c *Start) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&Start{})
}
