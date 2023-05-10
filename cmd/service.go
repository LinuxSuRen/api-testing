// Package cmd provides a service command
package cmd

import (
	"fmt"
	"os"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"github.com/spf13/cobra"
)

func createServiceCommand(execer fakeruntime.Execer) (c *cobra.Command) {
	opt := &serviceOption{
		Execer: execer,
	}
	c = &cobra.Command{
		Use:     "service",
		Aliases: []string{"s"},
		Short:   "Install atest as a Linux service",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	flags := c.Flags()
	flags.StringVarP(&opt.action, "action", "a", "", "The action of service, support actions: install, start, stop, restart, status")
	flags.StringVarP(&opt.scriptPath, "script-path", "", "/lib/systemd/system/atest.service", "The service script file path")
	return
}

type serviceOption struct {
	action     string
	scriptPath string
	fakeruntime.Execer
}

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	if o.Execer.OS() != "linux" {
		err = fmt.Errorf("only support on Linux")
	}
	if o.action == "" && len(args) > 0 {
		o.action = args[0]
	}
	return
}

func (o *serviceOption) runE(c *cobra.Command, args []string) (err error) {
	var output string
	switch o.action {
	case "install", "i":
		if err = os.WriteFile(o.scriptPath, []byte(script), os.ModeAppend); err == nil {
			output, err = o.Execer.RunCommandAndReturn("systemctl", "", "enable", "atest")
		}
	case "start":
		output, err = o.Execer.RunCommandAndReturn("systemctl", "", "start", "atest")
	case "stop":
		output, err = o.Execer.RunCommandAndReturn("systemctl", "", "stop", "atest")
	case "restart":
		output, err = o.Execer.RunCommandAndReturn("systemctl", "", "restart", "atest")
	case "status":
		output, err = o.Execer.RunCommandAndReturn("systemctl", "", "status", "atest")
	default:
		err = fmt.Errorf("not support action: '%s'", o.action)
	}

	if output != "" {
		c.Println(output)
	}
	return
}

var script = `[Unit]
Description=API Testing

[Service]
ExecStart=/usr/bin/env atest server

[Install]
WantedBy=multi-user.target
`
