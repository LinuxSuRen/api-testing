// Package cmd provides a service command
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	_ "embed"

	"github.com/linuxsuren/api-testing/cmd/service"
	"github.com/linuxsuren/api-testing/pkg/version"
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
		Example: `atest service install
atest service start`,
		Short: "Install atest as service",
		Long: `It could be a native or container service.
Try use sudo if you met any permission issues.

You could choose the alternative images:
Docker Hub: docker.io/linuxsuren/api-testing
GitHub Container Registry: ghcr.io/linuxsuren/api-testing
Scarf: linuxsuren.docker.scarf.sh/linuxsuren/api-testing
AliYun: registry.aliyuncs.com/linuxsuren/api-testing
DaoCloud: docker.m.daocloud.io/linuxsuren/api-testing`,
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
		Args:    cobra.MinimumNArgs(1),
	}
	flags := c.Flags()
	flags.StringVarP(&opt.scriptPath, "script-path", "", "", "The service script file path")
	flags.StringVarP(&opt.mode, "mode", "m", "",
		fmt.Sprintf("Availeble values: %v", service.ServiceModeOS.All()))
	flags.StringVarP(&opt.image, "image", "", service.DefaultImage, "The image of the service which as a container")
	flags.StringVarP(&opt.pull, "pull", "", "always", `Pull image before creating ("always"|"missing"|"never")`)
	flags.StringVarP(&opt.version, "version", "", version.GetVersion(), "The version of the service image")
	flags.StringVarP(&opt.LocalStorage, "local-storage", "", "/var/data/atest",
		"The local storage path which will be mounted into the container")
	flags.StringVarP(&opt.SecretServer, "secret-server", "", "", "The secret server URL")
	flags.StringVarP(&opt.SkyWalking, "skywalking", "", "", "Push the browser tracing data to the Apache SkyWalking URL")
	return
}

type serviceOption struct {
	action     string
	scriptPath string
	service    service.Service
	image      string
	version    string
	fakeruntime.Execer
	mode string
	pull string

	service.ServerFeatureOption
	stdOut io.Writer
}

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	o.stdOut = c.OutOrStdout()
	o.action = args[0]

	o.service = service.GetAvailableService(service.ServiceMode(o.mode), o.Execer,
		service.ContainerOption{
			Image:  o.action,
			Pull:   o.pull,
			Tag:    o.version,
			Writer: c.OutOrStdout(),
		}, o.ServerFeatureOption, o.scriptPath)

	if o.service == nil {
		err = fmt.Errorf("not supported service")
	} else if err == nil {
		local := os.ExpandEnv("$HOME/.config/atest")
		if err = o.Execer.MkdirAll(local, os.ModePerm); err == nil {
			err = o.Execer.MkdirAll(o.LocalStorage, os.ModePerm)
		}
	}
	return
}

func (o *serviceOption) runE(c *cobra.Command, args []string) (err error) {
	var output string
	switch Action(o.action) {
	case ActionInstall:
		output, err = o.service.Install()
	case ActionUninstall:
		output, err = o.service.Uninstall()
	case ActionStart:
		output, err = o.service.Start()
	case ActionStop:
		output, err = o.service.Stop()
	case ActionRestart:
		output, err = o.service.Restart()
	case ActionStatus:
		output, err = o.service.Status()
	default:
		err = fmt.Errorf("not support action: '%s'", o.action)
	}

	output = strings.TrimSpace(output)
	if output != "" {
		c.Println(output)
	}
	return
}

type Action string

const (
	ActionInstall   Action = "install"
	ActionUninstall Action = "uninstall"
	ActionStart     Action = "start"
	ActionStop      Action = "stop"
	ActionRestart   Action = "restart"
	ActionStatus    Action = "status"
)
