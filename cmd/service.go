// Package cmd provides a service command
package cmd

import (
	"fmt"
	"os"

	_ "embed"

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
		Short:   "Install atest as service",
		Long:    `It could be a native or container service`,
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
	flags.StringVarP(&opt.version, "version", "", version.GetVersion(), "The version of the service image")
	return
}

type serviceOption struct {
	action     string
	scriptPath string
	service    Service
	image      string
	version    string
	fakeruntime.Execer
	mode string
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

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	if o.action == "" && len(args) > 0 {
		o.action = args[0]
	}

	switch serviceMode(o.mode) {
	case ServiceModeOS:
		o.service, err = o.getOSService()
	default:
		o.service, err = o.getContainerService()
	}

	if o.service == nil {
		err = fmt.Errorf("not supported service")
		return
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

// Service is the interface of service
type Service interface {
	Start() (string, error)   // start the service
	Stop() (string, error)    // stop the service gracefully
	Restart() (string, error) // restart the service gracefully
	Status() (string, error)  // status of the service
	Install() (string, error) // install the service
	Uninstall() (string, error)
}

func emptyThenDefault(value, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func (o *serviceOption) getOSService() (service Service, err error) {
	if o.Execer.OS() != fakeruntime.OSLinux && o.Execer.OS() != fakeruntime.OSDarwin {
		err = fmt.Errorf("only support on Linux/Darwin instead of %s", o.Execer.OS())
	} else {
		service = newService(o.Execer, o.scriptPath)
	}
	return
}

func newService(execer fakeruntime.Execer, scriptPath string) (service Service) {
	switch execer.OS() {
	case fakeruntime.OSDarwin:
		service = &macOSService{
			commonService: commonService{
				Execer:     execer,
				scriptPath: emptyThenDefault(scriptPath, "/Library/LaunchDaemons/com.github.linuxsuren.atest.plist"),
				script:     macOSServiceScript,
			},
			cli: "launchctl",
			id:  "com.github.linuxsuren.atest",
		}
	case fakeruntime.OSLinux:
		service = &linuxService{
			commonService: commonService{
				Execer:     execer,
				scriptPath: emptyThenDefault(scriptPath, "/lib/systemd/system/atest.service"),
				script:     linuxServiceScript,
			},
		}
	}
	return
}

func (o *serviceOption) getContainerService() (service Service, err error) {
	var client string
	switch serviceMode(o.mode) {
	case ServiceModeDocker:
		client = "docker"
	case ServiceModePodman, ServiceModeContainer:
		client = "podman"
	default:
		err = fmt.Errorf("not support mode: '%s'", o.mode)
		return
	}

	if client, err = o.LookPath(client); err == nil {
		service = newContainerService(o.Execer, client, o.image, o.version)
	}
	return
}

type commonService struct {
	fakeruntime.Execer
	scriptPath string
	script     string
}

type macOSService struct {
	commonService
	cli string
	id  string
}

var (
	//go:embed data/macos_service.xml
	macOSServiceScript string
	//go:embed data/linux_service.txt
	linuxServiceScript string
)

func (s *macOSService) Start() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "start", s.id)
	return
}

func (s *macOSService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "stop", s.id)
	return
}

func (s *macOSService) Restart() (output string, err error) {
	if output, err = s.Stop(); err == nil {
		output, err = s.Start()
	}
	return
}

func (s *macOSService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "runstats", s.id)
	return
}

func (s *macOSService) Install() (output string, err error) {
	if err = os.WriteFile(s.scriptPath, []byte(s.script), os.ModeAppend); err == nil {
		output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "enable", s.id)
	}
	return
}

func (s *macOSService) Uninstall() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", s.cli, "disable", s.id)
	return
}

type linuxService struct {
	commonService
}

func (s *linuxService) Start() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("systemctl", "", "start", "atest")
	return
}

func (s *linuxService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("systemctl", "", "stop", "atest")
	return
}

func (s *linuxService) Restart() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("systemctl", "", "restart", "atest")
	return
}

func (s *linuxService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("systemctl", "", "status", "atest")
	return
}

func (s *linuxService) Install() (output string, err error) {
	if err = os.WriteFile(s.scriptPath, []byte(s.script), os.ModeAppend); err == nil {
		output, err = s.Execer.RunCommandAndReturn("systemctl", "", "enable", "atest")
	}
	return
}

func (s *linuxService) Uninstall() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("systemctl", "", "disable", "atest")
	return
}

type containerService struct {
	Execer fakeruntime.Execer
	name   string
	client string
	image  string
	tag    string
}

const defaultImage = "ghcr.io/linuxsuren/api-testing"

func newContainerService(execer fakeruntime.Execer, client, image, tag string) (service Service) {
	if tag == "" {
		tag = "latest"
	}
	if image == "" {
		image = defaultImage
	}
	service = &containerService{
		Execer: execer,
		client: client,
		name:   "atest",
		image:  image,
		tag:    tag,
	}
	return
}

func (s *containerService) Start() (output string, err error) {
	err = s.Execer.SystemCall(s.client, []string{s.client, "run", "--name=" + s.name,
		"--restart=always", "-d", "--pull=always", "--network=host",
		s.image + ":" + s.tag}, os.Environ())
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

func (s *containerService) Status() (output string, err error) {
	return
}

func (s *containerService) Install() (output string, err error) {
	output, err = s.Start()
	return
}

func (s *containerService) Uninstall() (output string, err error) {
	output, err = s.Stop()
	return
}
