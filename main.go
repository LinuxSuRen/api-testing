package main

import (
	"os"

	"github.com/linuxsuren/api-testing/cmd"
)

func main() {
	c := cmd.NewRootCmd()
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
