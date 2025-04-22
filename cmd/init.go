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
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/spf13/cobra"
)

type initOption struct {
	execer        fakeruntime.Execer
	kustomization string
	waitNamespace string
	waitResource  string
}

// createInitCommand returns the init command
func createInitCommand(execer fakeruntime.Execer) (cmd *cobra.Command) {
	opt := &initOption{execer: execer}
	cmd = &cobra.Command{
		Use:    "init",
		Long:   "Support to init Kubernetes cluster with kustomization, and wait it with command: kubectl wait",
		Hidden: true,
		RunE:   opt.runE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.kustomization, "kustomization", "k", "", "The kustomization file path")
	flags.StringVarP(&opt.waitNamespace, "wait-namespace", "", "", "")
	flags.StringVarP(&opt.waitResource, "wait-resource", "", "", "")
	return
}

func (o *initOption) runE(cmd *cobra.Command, args []string) (err error) {
	if o.kustomization != "" {
		if err = o.execer.RunCommand("kubectl", "apply", "-k", o.kustomization, "--wait=true"); err != nil {
			return
		}
	}

	if o.waitNamespace != "" && o.waitResource != "" {
		if err = o.execer.RunCommand("kubectl", "wait", "-n", o.waitNamespace, o.waitResource, "--for", "condition=Available=True", "--timeout=900s"); err != nil {
			return
		}
	}
	return
}
