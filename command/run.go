package command

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/light-speak/lighthouse/version"
)

func Run(c CommandList, args []string) error {
	if len(args) == 2 {
		var command *Command
		for _, cmd := range c.GetCommands() {
			if cmd.Name() == args[1] {
				command = &cmd
			}
		}
		if command == nil {
			return fmt.Errorf("unknown command: %s", args[1])
		}
		return runREPL(*command)
	}

	return runCommand(c, args)
}

func runREPL(cmd Command) error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		<-sigChan
		fmt.Print("\r")
		fmt.Println("\033[33mThank you for using lighthouse cli. Bye!\033[0m")
		cmd.OnExit()
		os.Exit(0)
	}()

	printLogo()
	fmt.Println("Welcome to use lighthouse cli. Type 'exit' to quit.")
	fmt.Printf("\033[36m%s\033[0m\n", cmd.Usage())

	scanner := bufio.NewScanner(os.Stdin)
	args := make(map[string]interface{})
	for _, arg := range cmd.Args() {
		for {
			fmt.Printf("\033[32mPlease input \033[35m%s\033[0m [%s] ", arg.Usage, GetTypeName(arg.Type))
			if arg.Required {
				fmt.Print("\033[31mrequired\033[0m")
			} else {
				fmt.Printf("\033[33mdefault: %s\033[0m", arg.Default)
			}
			fmt.Print(" : ")
			if !scanner.Scan() {
				return fmt.Errorf("failed to read input")
			}
			input := scanner.Text()
			if input == "exit" {
				return nil
			}
			if input != "" {
				args[arg.Name] = &input
				break
			}
			if input == "" && !arg.Required && arg.Default != nil {
				break
			}
			fmt.Println("Input cannot be empty. Please try again.")
		}
	}
	return cmd.Action()(args)
}

func runCommand(c CommandList, args []string) error {
	cmdName := args[1]
	cmd := findCommand(c, cmdName)
	if cmd == nil {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	flagSet, flagValues := setupFlags(cmd)
	help := flagSet.Bool("help", false, "show help")

	if err := flagSet.Parse(args[2:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}

	if len(args) == 2 {
		cmd.Action()(flagValues)
		return nil
	}
	if *help {
		printHelp(cmd, cmdName, flagSet)
		return nil
	}

	if err := validateRequiredFlags(cmd, flagValues); err != nil {
		return err
	}

	return cmd.Action()(flagValues)
}

func findCommand(c CommandList, name string) Command {
	for _, cmd := range c.GetCommands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func setupFlags(cmd Command) (*flag.FlagSet, map[string]interface{}) {
	flagSet := flag.NewFlagSet(cmd.Name(), flag.ContinueOnError)
	flagValues := make(map[string]interface{})

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

	return flagSet, flagValues
}

func printHelp(cmd Command, cmdName string, flagSet *flag.FlagSet) {
	fmt.Print("\033[33m")
	printLogo()
	fmt.Print("\033[1mCommand: ")
	fmt.Printf("\033[32m%s [flags]\n", cmdName)
	fmt.Printf("\033[0m%s\n", cmd.Usage())
	flagSet.PrintDefaults()
}

func validateRequiredFlags(cmd Command, flagValues map[string]interface{}) error {
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
	return nil
}

func printLogo() {
	fmt.Print("\033[33m")
	fmt.Printf(`
  _     _        _	 _    _
   |     )        |       |    |
   |    _    __   |__  _  |_   |__     _   _   _  ___   __
   |     |  '   \     \     |      \     \  |   |   __|  _ \
   |___  | |    | |   |   |_   |   | |    | |_  |  __ \ '__  
       | |  ' - | |   |     /  |   |  '- /      /     /    |
                |                 
            \__ /           %s           by @light-speak
	`, version.Version)
	fmt.Print("\033[0m\n")
}
