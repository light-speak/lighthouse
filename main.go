package main

import (
	"os"

	"github.com/light-speak/lighthouse/command"
	"github.com/light-speak/lighthouse/command/cli"
	"github.com/light-speak/lighthouse/log"
)

func main() {
	lighthouse := &cli.Lighthouse{}
	if err := command.Run(lighthouse, os.Args); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}
