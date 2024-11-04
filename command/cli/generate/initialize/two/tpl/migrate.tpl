type Migrate struct{}

func (c *Migrate) Name() string {
	// Func:Name user code start. Do not remove this comment.
	return "app:migrate"
	// Func:Name user code end. Do not remove this comment.
}

func (c *Migrate) Usage() string {
	// Func:Usage user code start. Do not remove this comment.
	return "This is a command to migrate the database"
	// Func:Usage user code end. Do not remove this comment.
}

func (c *Migrate) Args() []*command.CommandArg {
	return []*command.CommandArg{
		// Func:Args user code start. Do not remove this comment.
		// Func:Args user code end. Do not remove this comment.
	}
}

func (c *Migrate) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		// Func:Action user code start. Do not remove this comment.
		err := models.Migrate()
		if err != nil {
			return err
		}
		// Func:Action user code end. Do not remove this comment.
		return nil
	}
}

func (c *Migrate) OnExit() func() {
	return func() {}
}
// Section: user code section start. Do not remove this comment.
// Section: user code section end. Do not remove this comment.
