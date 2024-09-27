package main

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/command"
)

func main() {
	lighthouseCommand := &command.LighthouseCommandList{}
	if err := command.Run(lighthouseCommand, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
