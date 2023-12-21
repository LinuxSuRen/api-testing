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
	"github.com/linuxsuren/api-testing/extensions/store-etcd/pkg"
	ext "github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root Command
func NewRootCommand() (c *cobra.Command) {
	opt := &options{
		Extension: ext.NewExtension("etcd", "store", 7073),
	}
	c = &cobra.Command{
		Use:   opt.GetFullName(),
		Short: "A store extension for etcd",
		RunE:  opt.runE,
	}
	opt.AddFlags(c.Flags())
	return
}

type options struct {
	*ext.Extension
}

func (o *options) runE(c *cobra.Command, _ []string) (err error) {
	remoteServer := pkg.NewRemoteServer(pkg.NewRealEtcd())
	err = ext.CreateRunner(o.Extension, c, remoteServer)
	return
}
