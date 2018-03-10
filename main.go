package main

import (
	"os"

	"github.com/uphy/commandbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
