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
		log.Error("请输入一个命令，如：init、gql:init 或 gql:generate")
		return
	}

	command := os.Args[1]
	var err error

	switch command {
	case "init":
		err = initialize.Run()
	case "gql:init":
		err = gqlInitialize.Run()
	case "gql:generate":
		err = gqlInitialize.Run()
		if err != nil {
			log.Error("执行 gqlInitialize.Run() 失败: %v", err)
			return
		}
		err = gqlGenerate.Run()
	default:
		log.Error("命令 %s 不正确，请输入一个有效命令：init、gql:init 或 gql:generate", command)
		return
	}

	if err != nil {
		log.Error("执行命令 %s 失败: %v", command, err)
	}
}
