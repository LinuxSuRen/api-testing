package cmd

import (
	"github.com/linuxsuren/api-testing/extensions/store-s3/pkg"
	ext "github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/spf13/cobra"
)

func NewRootCmd(s3Creator pkg.S3Creator) (c *cobra.Command) {
	opt := &option{
		s3Creator: s3Creator,
		Extension: ext.NewExtension("s3", 7072),
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
