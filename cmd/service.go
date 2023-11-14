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
DaoCloud: docker.m.daocloud.io/linuxsuren/api-testing`,
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	flags := c.Flags()
	flags.StringVarP(&opt.action, "action", "a", "",
		fmt.Sprintf("The action of the service, available values: %v", Action("").All()))
	flags.StringVarP(&opt.scriptPath, "script-path", "", "", "The service script file path")
	flags.StringVarP(&opt.mode, "mode", "m", string(ServiceModeOS),
		fmt.Sprintf("Availeble values: %v", ServiceModeOS.All()))
	flags.StringVarP(&opt.image, "image", "", defaultImage, "The image of the service which as a container")
	flags.StringVarP(&opt.pull, "pull", "", "always", `Pull image before creating ("always"|"missing"|"never")`)
	flags.StringVarP(&opt.version, "version", "", version.GetVersion(), "The version of the service image")
	flags.StringVarP(&opt.localStorage, "local-storage", "", "/var/data/atest",
		"The local storage path which will be mounted into the container")
	flags.StringVarP(&opt.secretServer, "secret-server", "", "", "The secret server URL")
	flags.StringVarP(&opt.skyWalking, "skywalking", "", "", "Push the browser tracing data to the Apache SkyWalking URL")
	return
}

type serverFeatureOption struct {
	secretServer string
	skyWalking   string
	localStorage string
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

	serverFeatureOption
	stdOut io.Writer
}

type serviceMode string

const (
	ServiceModeOS        serviceMode = "os"
	ServiceModeContainer serviceMode = "container"
	ServiceModePodman    serviceMode = "podman"
	ServiceModeDocker    serviceMode = "docker"
)

func (s serviceMode) All() []serviceMode {
	return []serviceMode{ServiceModeOS, ServiceModeContainer,
		ServiceModePodman, ServiceModeDocker}
}

func (s serviceMode) String() string {
	return string(s)
}

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	o.stdOut = c.OutOrStdout()
	if o.action == "" && len(args) > 0 {
		o.action = args[0]
	}

	switch serviceMode(o.mode) {
	case ServiceModeOS:
		o.service = service.NewService(o.Execer, o.scriptPath)
	default:
		o.service, err = o.getContainerService()
	}

	if o.service == nil {
		err = fmt.Errorf("not supported service")
		return
	} else if err == nil {
		err = o.Execer.MkdirAll(o.localStorage, os.ModePerm)
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

func (a Action) All() []Action {
	return []Action{ActionInstall, ActionUninstall,
		ActionStart, ActionStop,
		ActionRestart, ActionStatus}
}

func emptyThenDefault(value, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func (o *serviceOption) getOSService() (svc service.Service, err error) {
	//if o.Execer.OS() != fakeruntime.OSLinux && o.Execer.OS() != fakeruntime.OSDarwin {
	//	err = fmt.Errorf("only support on Linux/Darwin instead of %s", o.Execer.OS())
	//} else {
	//	svc = service.NewService(o.Execer, o.scriptPath)
	//}
	svc = service.NewService(o.Execer, o.scriptPath)
	return
}

func (o *serviceOption) getContainerService() (service service.Service, err error) {
	var client string
	switch serviceMode(o.mode) {
	case ServiceModeDocker:
		client = ServiceModeDocker.String()
	case ServiceModePodman, ServiceModeContainer:
		client = ServiceModePodman.String()
	default:
		err = fmt.Errorf("not support mode: '%s'", o.mode)
		return
	}

	var clientPath string
	if clientPath, err = o.LookPath(client); err == nil {
		if clientPath == "" {
			clientPath = client
		}
		service = newContainerService(o.Execer, clientPath,
			o.image, o.version, o.pull, o.serverFeatureOption, o.stdOut)
	}
	return
}

type commonService struct {
	fakeruntime.Execer
	scriptPath string
	script     string
}

type containerService struct {
	Execer       fakeruntime.Execer
	name         string
	client       string
	image        string
	tag          string
	pull         string
	localStorage string
	secretServer string
	skyWalking   string
	stdOut       io.Writer
	errOut       io.Writer
}

const defaultImage = "linuxsuren.docker.scarf.sh/linuxsuren/api-testing"

func newContainerService(execer fakeruntime.Execer, client, image, tag, pull string,
	featureOption serverFeatureOption, writer io.Writer) (svc service.Service) {
	if tag == "" {
		tag = "latest"
	}
	if image == "" {
		image = defaultImage
	}

	containerServer := &containerService{
		Execer:       execer,
		client:       client,
		name:         service.ServiceName,
		image:        image,
		tag:          tag,
		pull:         pull,
		localStorage: featureOption.localStorage,
		secretServer: featureOption.secretServer,
		skyWalking:   featureOption.skyWalking,
		stdOut:       writer,
		errOut:       writer,
	}

	if strings.HasSuffix(client, ServiceModePodman.String()) {
		svc = &podmanService{
			containerService: containerServer,
		}
	} else {
		svc = containerServer
	}
	return
}

func (s *containerService) Start() (output string, err error) {
	if s.exist() {
		output, err = s.Execer.RunCommandAndReturn(s.client, "", "start", s.name)
	} else {
		err = s.Execer.SystemCall(s.client, append([]string{s.client}, s.getStartArgs()...), os.Environ())
	}
	return
}

func (s *containerService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(s.client, "", "stop", s.name)
	return
}

func (s *containerService) Restart() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(s.client, "", "restart", s.name)
	return
}

func (s *containerService) Status() (_ string, err error) {
	err = s.Execer.SystemCall(s.client, []string{s.client, "stats", s.name}, os.Environ())
	return
}

func (s *containerService) Install() (output string, err error) {
	output, err = s.Start()
	return
}

func (s *containerService) Uninstall() (output string, err error) {
	output, err = s.Stop()
	if err == nil {
		output, err = s.Execer.RunCommandAndReturn(s.client, "", "rm", s.name)
	}
	return
}

func (s *containerService) exist() bool {
	output, err := s.Execer.RunCommandAndReturn(s.client, "", "ps", "--all", "--filter", fmt.Sprintf("name=%s", s.name))
	return err == nil && strings.Contains(output, s.name)
}

func (s *containerService) getStartArgs() []string {
	args := []string{"run", "--name=" + s.name,
		"--restart=always",
		"-d",
		fmt.Sprintf("--pull=%s", s.pull),
		"--network=host",
		"-v", s.localStorage + ":/var/www/data",
		"-v", os.ExpandEnv("$HOME/.config/atest:/root/.config/atest"),
		s.image + ":" + s.tag,
		"atest", "server"}
	if s.secretServer != "" {
		args = append(args, "--secret-server="+s.secretServer)
	}
	if s.skyWalking != "" {
		args = append(args, "--skywalking="+s.skyWalking)
	}
	args = append(args, "--console-path=/var/www/html")
	return args
}

type podmanService struct {
	*containerService
}

func (s *podmanService) Install() (output string, err error) {
	output, err = s.Start()
	return
}

func (s *podmanService) Start() (output string, err error) {
	if s.exist() {
		err = s.Execer.RunCommandWithIO(s.client, "", s.stdOut, s.errOut, "start", s.name)
	} else {
		err = s.Execer.RunCommandWithIO(s.client, "", s.stdOut, s.errOut, s.getStartArgs()...)
		if err == nil {
			output, err = s.installService()
		}
	}
	return
}

func (s *podmanService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(service.SystemCtl, "", "stop", service.PodmanServiceName)
	return
}

func (s *podmanService) installService() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(s.client, "", "generate", "systemd", "--new", "--files", "--name", s.name)
	if err == nil {
		var result string
		result, err = s.Execer.RunCommandAndReturn("mv", "", service.PodmanServiceName, "/etc/systemd/system")
		if err == nil {
			output = fmt.Sprintf("%s\n%s", output, result)
			if result, err = s.Execer.RunCommandAndReturn(service.SystemCtl, "", "enable", service.PodmanServiceName); err == nil {
				output = fmt.Sprintf("%s\n%s", output, result)
			}
		}
	}
	return
}

func (s *podmanService) Uninstall() (output string, err error) {
	output, err = s.containerService.Uninstall()
	if err == nil {
		var result string
		if result, err = s.uninstallService(); err == nil {
			output = fmt.Sprintf("%s\n%s", output, result)
		}
	}
	return
}

func (s *podmanService) uninstallService() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(service.SystemCtl, "", "disable", service.PodmanServiceName)
	return
}
