/*
Copyright 2025 API Testing Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
