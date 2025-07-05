package main

import (
	"os"

	"github.com/light-speak/lighthouse/lightcmd/cmd"
	"github.com/light-speak/lighthouse/logs"
)

func main() {
	command := &cmd.Command{}
	if err := cmd.Run(command, os.Args); err != nil {
		logs.Error().Err(err).Msg("")
		os.Exit(1)
	}
}
