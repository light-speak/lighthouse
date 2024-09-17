package main

import (
	"os"

	gqlGenerate "github.com/light-speak/lighthouse/graphql/generate"
	gqlInitialize "github.com/light-speak/lighthouse/graphql/initialize"
	"github.com/light-speak/lighthouse/initialize"
	"github.com/light-speak/lighthouse/log"
)

func main() {

	if len(os.Args) < 2 {
		log.Error("请输入一个命令， 如： init 或 gql:init")
		return
	}

	command := os.Args[1]
	switch command {
	case "init":
		if err := initialize.Run(); err != nil {
			log.Error("\nfailed to run initialize command: %v ", err)
		}
		break

	case "gql:init":
		if err := gqlInitialize.Run(); err != nil {
			log.Error("\nfailed to run gql:init command: %v ", err)
		}
		break

	case "gql:generate":
		if err := gqlGenerate.Run(); err != nil {
			log.Error("\nfailed to run gql:generate command: %v ", err)
		}
		break
	default:
		log.Error("命令 %s 不正确，请输入一个命令， 如： init 或 gql:init", command)
	}
}
