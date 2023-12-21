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
	"github.com/linuxsuren/api-testing/extensions/store-s3/pkg"
	ext "github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/spf13/cobra"
)

func NewRootCmd(s3Creator pkg.S3Creator) (c *cobra.Command) {
	opt := &option{
		s3Creator: s3Creator,
		Extension: ext.NewExtension("s3", "store", 7072),
	}
	c = &cobra.Command{
		Use:   opt.GetFullName(),
		Short: "S3 storage extension of api-testing",
		RunE:  opt.runE,
	}
	opt.AddFlags(c.Flags())
	return
}

func (o *option) runE(c *cobra.Command, _ []string) (err error) {
	remoteServer := pkg.NewRemoteServer(o.s3Creator)
	err = ext.CreateRunner(o.Extension, c, remoteServer)
	return
}

type option struct {
	// inner fields
	s3Creator pkg.S3Creator
	*ext.Extension
}
