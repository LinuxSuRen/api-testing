package main

import (
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := &cobra.Command{
		Use: "atest",
	}
	cmd.AddCommand(createInitCommand(), createRunCommand())

	// run command
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
