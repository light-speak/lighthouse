package main

import "search/service"

func main() {
	service.RunServer("8888", "http://localhost:9200")
}
