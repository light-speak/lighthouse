// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
type SubscribeStart struct{}

func (s *SubscribeStart) Name() string {
	// Func:Name user code start. Do not remove this comment.
	return "subscribe:start"
	// Func:Name user code end. Do not remove this comment.
}

func (s *SubscribeStart) Usage() string {
	// Func:Usage user code start. Do not remove this comment.
	return "This is a command to start subscribe kafka messages"
	// Func:Usage user code end. Do not remove this comment.
}

func (s *SubscribeStart) Args() []*command.CommandArg {
	return []*command.CommandArg{}
}

func (s *SubscribeStart) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		// Func:Action user code start. Do not remove this comment.
		subscriber.Start()
		// Func:Action user code end. Do not remove this comment.
		return nil
	}
}

func (s *SubscribeStart) OnExit() func() {
	return func() {}
}

// Section: user code section start. Do not remove this comment.
// Section: user code section end. Do not remove this comment.
