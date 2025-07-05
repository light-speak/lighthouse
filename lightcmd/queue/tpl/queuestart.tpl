type Queue struct{}

func (c *Queue) Name() string {
	return "queue:start"
}

func (c *Queue) Usage() string {
	return "This is a command to start the queue service"
}

func (c *Queue) Args() []*cmd.CommandArg {
	return []*cmd.CommandArg{}
}

func (c *Queue) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		lightqueue.StartQueue()
		return nil
	}
}

func (c *Queue) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&Queue{})
}
