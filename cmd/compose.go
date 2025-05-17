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
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type composeOptions struct {
	runOption
	projectName string
}

func createComposeRun() *cobra.Command {
	cmd := &cobra.Command{
		Use: "compose",
	}

	cmd.AddCommand(createComposeRunUp())
	return cmd
}

type composeMsgWriter struct {
}

func (w *composeMsgWriter) Write(p []byte) (n int, err error) {
	fmt.Println(`{ "type": "info", "message": "` + strings.TrimSpace(string(p)) + `" }`)
	return
}

func (o *composeOptions) preRunE(cmd *cobra.Command, args []string) error {
	return o.runOption.preRunE(cmd, nil)
}

func (o *composeOptions) runE(cmd *cobra.Command, args []string) error {
	return o.runOption.runE(cmd, args)
}

func createComposeRunUp() *cobra.Command {
	opt := &composeOptions{
		runOption: *newDefaultRunOption(&composeMsgWriter{}),
	}
	c := &cobra.Command{
		Use:     "up",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	c.SetOut(&composeMsgWriter{})

	c.Flags().StringVarP(&opt.projectName, "project-name", "", "", "Specify an alternate project name")
	opt.runOption.addFlags(c.Flags())
	return c
}
