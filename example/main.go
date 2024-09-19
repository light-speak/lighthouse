package main

import (
	"example/graph"

	"github.com/light-speak/lighthouse/log"
)

func main() {
	log.Info("应用启动！")
	graph.StartServer()
}
