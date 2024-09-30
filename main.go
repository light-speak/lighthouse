package main

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/command/cli"
)

func main() {
	
	lighthouse := &cli.Lighthouse{}
	if err := command.Run(lighthouse, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
