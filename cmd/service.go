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
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/util/home"

	_ "embed"

	"github.com/linuxsuren/api-testing/pkg/version"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	service "github.com/linuxsuren/go-service"
	"github.com/spf13/cobra"
)

const modeDockerInSystemService = "docker-in-system"

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
		fmt.Sprintf("Availeble values: %v", append(service.ServiceModeOS.All(), modeDockerInSystemService)))
	flags.StringVarP(&opt.image, "image", "", defaultImage, "The image of the service which as a container")
	flags.StringVarP(&opt.pull, "pull", "", "always", `Pull image before creating ("always"|"missing"|"never")`)
	flags.StringVarP(&opt.version, "version", "", version.GetVersion(), "The version of the service image")
	flags.StringVarP(&opt.LocalStorage, "local-storage", "", "/var/data/atest",
		"The local storage path which will be mounted into the container")
	flags.IntVarP(&opt.port, "port", "", 8080, "The port of the service")
	flags.StringVarP(&opt.SecretServer, "secret-server", "", "", "The secret server URL")
	flags.StringVarP(&opt.SkyWalking, "skywalking", "", "", "Push the browser tracing data to the Apache SkyWalking URL")
	return
}

const defaultImage = "ghcr.io/linuxsuren/api-testing"

type serviceOption struct {
	action     string
	scriptPath string
	service    service.Service
	image      string
	version    string
	fakeruntime.Execer
	mode string
	pull string
	port int

	SecretServer string
	SkyWalking   string
	LocalStorage string
	stdOut       io.Writer
}

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	o.stdOut = c.OutOrStdout()
	o.action = args[0]

	serviceCommand := "atest"
	serviceArgs := []string{"server", "--extension-registry=ghcr.io"}
	if o.mode == modeDockerInSystemService {
		if o.Execer.OS() == "windows" {
			serviceCommand = "run"
			serviceArgs = []string{"atest", "server", "--extension-registry=ghcr.io"}
			o.mode = "docker"
		} else {
			serviceCommand = "docker"
			serviceArgs = []string{"run", "-v=atest:/root/.config/atest", "-v=atest-ssh:/root/.ssh",
				fmt.Sprintf("-p=%d:8080", o.port), "--pull", o.pull, fmt.Sprintf("%s:%s", o.image, o.version),
				"atest", "server", "--extension-registry=ghcr.io"}
			o.mode = string(service.ServiceModeOS)
		}
	}
	if strings.Contains(o.version, "unknown") {
		o.version = ""
	}

	if o.service, err = service.GetAvailableService(service.ServiceMode(o.mode),
		service.ContainerOption{
			Image:   o.image,
			Pull:    o.pull,
			Tag:     o.version,
			Writer:  c.OutOrStdout(),
			Volumes: map[string]string{"atest": "/root/.config/atest"},
			Ports:   map[int]int{o.port: 8080},
			Restart: "always",
			Name:    "atest",
		}, service.CommonService{
			ID:          "atest",
			Name:        "atest",
			Description: "API Testing Server",
			Command:     serviceCommand,
			Args:        serviceArgs,
			Execer:      o.Execer,
		}); err != nil {
		return
	}

	local := home.GetUserConfigDir()
	if err = o.Execer.MkdirAll(local, os.ModePerm); err == nil {
		err = o.Execer.MkdirAll(o.LocalStorage, os.ModePerm)
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
		err = fmt.Errorf("not support action: %q", o.action)
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
