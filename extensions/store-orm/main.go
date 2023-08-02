package main

import (
	"os"

	"github.com/linuxsuren/api-testing/extensions/store-orm/cmd"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
