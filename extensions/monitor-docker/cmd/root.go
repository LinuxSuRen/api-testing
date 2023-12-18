/*
Copyright 2023 API Testing Authors.

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
	"github.com/docker/cli/cli/command"
	"github.com/linuxsuren/api-testing/extensions/monitor-docker/pkg"
	ext "github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/spf13/cobra"
)

func NewRootCommand(dockerCli command.Cli) (c *cobra.Command) {
	opt := options{
		dockerCli: dockerCli,
		Extension: ext.NewExtension("docker", "monitor", 7074),
	}
	c = &cobra.Command{
		Use:  "server",
		RunE: opt.runE,
	}
	opt.AddFlags(c.Flags())
	return
}

type options struct {
	*ext.Extension
	dockerCli command.Cli
}

func (o *options) runE(c *cobra.Command, _ []string) (err error) {
	remoteServer := pkg.NewRemoteServer(o.dockerCli)
	err = ext.CreateMonitor(o.Extension, c, remoteServer)
	return
}
