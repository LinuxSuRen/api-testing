package main

import (
	"os"

	"github.com/linuxsuren/api-testing/cmd"
	exec "github.com/linuxsuren/go-fake-runtime"
)

func main() {
	c := cmd.NewRootCmd(exec.DefaultExecer{})
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
