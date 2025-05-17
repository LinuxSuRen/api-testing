package main

import (
	_ "embed"
	"os"

	"github.com/linuxsuren/api-testing/pkg/version"

	"github.com/linuxsuren/api-testing/cmd"
	"github.com/linuxsuren/api-testing/pkg/server"
	exec "github.com/linuxsuren/go-fake-runtime"
)

func main() {
	version.SetMod(goMod)
	c := cmd.NewRootCmd(exec.NewDefaultExecer(), server.NewDefaultHTTPServer())
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}

//go:embed go.mod
var goMod string
