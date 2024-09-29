package command

import (
	"flag"
	"fmt"
)

func Run(c CommandList, args []string) error {
	if len(args) < 2 {
		printLogo()
		return fmt.Errorf("please specify a command")
	}

	cmdName := args[1]
	for _, cmd := range c.GetCommands() {
		if cmd.Name() == cmdName {
			flagSet := flag.NewFlagSet(cmdName, flag.ContinueOnError)
			flagValues := make(map[string]interface{})
			help := flagSet.Bool("help", false, "show help")

			for _, arg := range cmd.Args() {
				switch arg.Type {
				case String:
					defaultValue := ""
					if arg.Default != nil {
						defaultValue = arg.Default.(string)
					}
					flagValues[arg.Name] = flagSet.String(arg.Name, defaultValue, fmt.Sprintf("%s (required: %t)", arg.Usage, arg.Required))
				case Int:
					flagValues[arg.Name] = flagSet.Int(arg.Name, 0, fmt.Sprintf("%s (required: %t) (default: %d)", arg.Usage, arg.Required, arg.Default))
				case Bool:
					flagValues[arg.Name] = flagSet.Bool(arg.Name, false, fmt.Sprintf("%s (required: %t) (default: %t)", arg.Usage, arg.Required, arg.Default))
				}
			}

			if err := flagSet.Parse(args[2:]); err != nil {
				return fmt.Errorf("failed to parse flags: %w", err)
			}

			// If no flags are provided or help flag is set, show help
			if len(args) == 2 || *help {
				fmt.Print("\033[33m")
				printLogo()
				fmt.Print("\033[1mCommand: ")
				fmt.Printf("\033[32m%s [flags]\n", cmdName)
				fmt.Printf("\033[0m%s\n", cmd.Usage())

				flagSet.PrintDefaults()
				return nil
			}

			// Only check for required parameters if help flag is not set
			for _, arg := range cmd.Args() {
				if arg.Required {
					switch arg.Type {
					case String:
						if *flagValues[arg.Name].(*string) == "" {
							return fmt.Errorf("missing required parameter: --%s", arg.Name)
						}
					case Int:
						if *flagValues[arg.Name].(*int) == 0 {
							return fmt.Errorf("missing required parameter: --%s", arg.Name)
						}
					case Bool:
					}
				}
			}

			return cmd.Action()(flagValues)
		}
	}

	return fmt.Errorf("unknown command: %s", cmdName)
}


func printLogo() {
	fmt.Print("\033[33m")
	fmt.Print(`
  _     _        _	 _    _
   |     )        |       |    |
   |    _    __   |__  _  |_   |__     _   _   _  ___   __
   |     |  '   \     \     |      \     \  |   |   __|  _ \
   |___  | |    | |   |   |_   |   | |    | |_  |  __ \ '__  
       | |  ' - | |   |     /  |   |  '- /      /     /    |
                |                 
            \__ /           v0.0.1           by @light-speak
	`)
	fmt.Print("\033[0m\n")
}
