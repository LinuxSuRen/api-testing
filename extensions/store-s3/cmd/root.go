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
