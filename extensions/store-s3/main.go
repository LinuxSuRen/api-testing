package main

import (
	"os"

	"github.com/linuxsuren/api-testing/extensions/store-s3/cmd"
	"github.com/linuxsuren/api-testing/extensions/store-s3/pkg"
)

func main() {
	c := cmd.NewRootCmd(&pkg.DefaultS3Creator{})
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
