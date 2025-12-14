package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(c CommandListInterface, args []string) error {
	// Check if user is asking for help
	if len(args) == 2 && (args[1] == "help" || args[1] == "--help") {
		printAllCommands(c)
		return nil
	}

	if len(args) == 2 {
		var command *CommandInterface
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

func printAllCommands(c CommandListInterface) {
	fmt.Println("\033[1m\033[33mAvailable Commands:\033[0m")
	for _, cmd := range c.GetCommands() {
		fmt.Printf("\033[32m%-15s\033[0m %s\n", cmd.Name(), cmd.Usage())
	}
	fmt.Println("\nUse 'command --help' for more information about a command.")
}

// setupSignalHandler 设置信号处理，返回一个 done channel 用于通知主程序退出
func setupSignalHandler(cmd CommandInterface) chan struct{} {
	sigChan := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Print("\r")
		fmt.Println("\n\033[33mReceived shutdown signal, cleaning up...\033[0m")
		if onExit := cmd.OnExit(); onExit != nil {
			onExit()
		}
		fmt.Println("\033[33mThank you for using light-cmd. Bye!\033[0m")
		close(done)
	}()

	return done
}

func runREPL(cmd CommandInterface) error {
	done := setupSignalHandler(cmd)

	// 如果命令没有参数，直接执行 Action
	if len(cmd.Args()) == 0 {
		errChan := make(chan error, 1)
		go func() {
			errChan <- cmd.Action()(make(map[string]interface{}))
		}()

		select {
		case <-done:
			return nil
		case err := <-errChan:
			return err
		}
	}

	fmt.Println("Welcome to use light-cmd. Type 'exit' to quit.")
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

	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Action()(args)
	}()

	select {
	case <-done:
		return nil
	case err := <-errChan:
		return err
	}
}

func runCommand(c CommandListInterface, args []string) error {
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

	if *help {
		printHelp(cmd, cmdName, flagSet)
		return nil
	}

	if err := validateRequiredFlags(cmd, flagValues); err != nil {
		return err
	}

	// 设置信号处理
	done := setupSignalHandler(cmd)

	errChan := make(chan error, 1)
	go func() {
		errChan <- cmd.Action()(flagValues)
	}()

	select {
	case <-done:
		return nil
	case err := <-errChan:
		return err
	}
}

func findCommand(c CommandListInterface, name string) CommandInterface {
	for _, cmd := range c.GetCommands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func setupFlags(cmd CommandInterface) (*flag.FlagSet, map[string]interface{}) {
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

func printHelp(cmd CommandInterface, cmdName string, flagSet *flag.FlagSet) {
	fmt.Print("\033[33m")
	fmt.Print("\033[1mCommand: ")
	fmt.Printf("\033[32m%s [flags]\n", cmdName)
	fmt.Printf("\033[0m%s\n", cmd.Usage())
	flagSet.PrintDefaults()
}

func validateRequiredFlags(cmd CommandInterface, flagValues map[string]interface{}) error {
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
