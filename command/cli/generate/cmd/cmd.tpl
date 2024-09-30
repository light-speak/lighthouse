{{ $structName := (.Name | ucFirst) }}
type {{ $structName }} struct {}

func (c *{{ $structName }}) Name() string {
	{{ funcStart "Name" }}
	return "{{.Scope | lcFirst}}:{{.Name | lcFirst}}"
	{{ funcEnd "Name" }}
}

func (c *{{ $structName }}) Usage() string {
	{{ funcStart "Usage" }}
	return "This is a command generated by lighthouse cli"
	{{ funcEnd "Usage" }}
}

func (c *{{ $structName }}) Args() []*command.CommandArg {
	return []*command.CommandArg{
		{{ funcStart "Args" }}
		{{ funcEnd "Args" }}
	}
}

func (c *{{ $structName }}) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		{{ funcStart "Action" }}
		args, err := command.GetArgs(c.Args(), flagValues)
		if err != nil {
			return err
		}
		log.Info().Msgf("args: %v", args)
		{{ funcEnd "Action" }}
		return nil
	}
}


{{ section }}