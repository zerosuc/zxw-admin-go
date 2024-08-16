package main

import (
	"fmt"
	"os"
	"server/core"
)

func main() {
	cmd := core.NewServerCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
