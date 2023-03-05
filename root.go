package main

import (
	"os"

	c "github.com/linuxsuren/api-testing/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "atest",
	}
	cmd.AddCommand(c.CreateInitCommand(), c.CreateRunCommand())

	// run command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
