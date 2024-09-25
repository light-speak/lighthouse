package main

import (
	"user/graph/server"

	"github.com/light-speak/lighthouse/log"
)

func main() {
	log.Info("应用启动！")
	server.StartServer()
}
