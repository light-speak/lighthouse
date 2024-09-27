package command

import (
	"flag"
	"fmt"
)

func Run(c CommandList, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("please specify a command")
	}

	cmdName := args[1]
	for _, cmd := range c.GetCommands() {
		if cmd.Name() == cmdName {
			flagSet := flag.NewFlagSet(cmdName, flag.ContinueOnError)
			flagValues := make(map[string]interface{})

			for _, arg := range cmd.Args() {
				switch arg.Type {
				case String:
					flagValues[arg.Name] = flagSet.String(arg.Name, "", arg.Usage)
				case Int:
					flagValues[arg.Name] = flagSet.Int(arg.Name, 0, arg.Usage)
				case Bool:
					flagValues[arg.Name] = flagSet.Bool(arg.Name, false, arg.Usage)
				}
			}

			if err := flagSet.Parse(args[2:]); err != nil {
				return fmt.Errorf("failed to parse flags: %w", err)
			}

			for _, arg := range cmd.Args() {
				if arg.Required {
					switch arg.Type {
					case String:
						if *flagValues[arg.Name].(*string) == "" {
							return fmt.Errorf("missing required parameters: --%s", arg.Name)
						}
					case Int:
						if *flagValues[arg.Name].(*int) == 0 {
							return fmt.Errorf("missing required parameters: --%s", arg.Name)
						}
					}
				}
			}

			return cmd.Action()(flagValues)
		}
	}

	return fmt.Errorf("unknown command: %s", cmdName)
}
