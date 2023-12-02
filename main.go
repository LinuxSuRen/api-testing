package main

import (
	"os"

	// _ "github.com/apache/skywalking-go"
	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	exec "github.com/linuxsuren/go-fake-runtime"
)

func main() {
	c := cmd.NewRootCmd(exec.NewDefaultExecer(), server.NewDefaultHTTPServer())
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
