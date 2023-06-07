package main

import (
	"os"

	"github.com/linuxsuren/api-testing/extensions/collector/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
