package main

import (
	"fmt"

	"github.com/light-speak/lighthouse/utils"
)

func main() {
	w, err := utils.GetCallerPath()
	if err != nil {
		panic(err)
	}
	fmt.Println(w)
}
