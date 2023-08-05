package main

import (
	"os"

	"github.com/linuxsuren/api-testing/extensions/store-s3/cmd"
	"github.com/linuxsuren/api-testing/extensions/store-s3/pkg"
)

func main() {
	cmd := cmd.NewRootCmd(&pkg.DefaultS3Creator{})
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
