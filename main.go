package main

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/command/cli/generate"
)

func main() {
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
	lighthouse := &generate.Lighthouse{}
	if err := command.Run(lighthouse, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
