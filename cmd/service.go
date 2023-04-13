package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func createServiceCommand() (c *cobra.Command) {
	opt := &serviceOption{}
	c = &cobra.Command{
		Use:     "service",
		Aliases: []string{"s"},
		Short:   "Install atest as a Linux service",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	flags := c.Flags()
	flags.StringVarP(&opt.action, "action", "a", "", "The action of service, support actions: install")
	return
}

type serviceOption struct {
	action string
}

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	if runtime.GOOS != "linux" {
		err = fmt.Errorf("only support on Linux")
	}
	return
}

func (o *serviceOption) runE(c *cobra.Command, args []string) (err error) {
	switch o.action {
	case "install", "i":
		err = os.WriteFile("/lib/systemd/system/atest.service", []byte(script), os.ModeAppend)
	default:
		err = fmt.Errorf("not support action: '%s'", o.action)
	}
	return
}

var script = `[Unit]
Description=API Testing

[Service]
ExecStart=atest server

[Install]
WantedBy=multi-user.target
`
