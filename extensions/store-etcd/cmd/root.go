/**
MIT License

Copyright (c) 2023 API Testing Authors.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
		Extension: ext.NewExtension("etcd", 7073),
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
