package cmd

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
	"github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/spf13/cobra"
)

func createJSONSchemaCmd() (c *cobra.Command) {
	c = &cobra.Command{
		Use:   "json",
		Short: "Print the JSON schema of the test suites struct",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var data []byte
			schema := jsonschema.Reflect(&testing.TestSuite{})
			if data, err = json.MarshalIndent(schema, "", "  "); err == nil {
				cmd.Println(string(data))
			}
			return
		},
	}
	return
}
