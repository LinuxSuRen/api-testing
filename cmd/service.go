// Package cmd provides a service command
package cmd

import (
	"fmt"
	"os"

	_ "embed"

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
	flags.StringVarP(&opt.scriptPath, "script-path", "", "", "The service script file path")
	return
}

type serviceOption struct {
	action     string
	scriptPath string
	service    Service
	fakeruntime.Execer
}

func (o *serviceOption) preRunE(c *cobra.Command, args []string) (err error) {
	if o.Execer.OS() != fakeruntime.OSLinux && o.Execer.OS() != fakeruntime.OSDarwin {
		err = fmt.Errorf("only support on Linux/Darwin instead of %s", o.Execer.OS())
	} else {
		if o.action == "" && len(args) > 0 {
			o.action = args[0]
		}
		o.service = newService(o.Execer, o.scriptPath)
	}
	return
}

func (o *serviceOption) runE(c *cobra.Command, args []string) (err error) {
	var output string
	switch o.action {
	case "install", "i":
		output, err = o.service.Install()
	case "start":
		output, err = o.service.Start()
	case "stop":
		output, err = o.service.Stop()
	case "restart":
		output, err = o.service.Restart()
	case "status":
		output, err = o.service.Status()
	default:
		err = fmt.Errorf("not support action: '%s'", o.action)
	}

	if output != "" {
		c.Println(output)
	}
	return
}

// Service is the interface of service
type Service interface {
	Start() (string, error)   // start the service
	Stop() (string, error)    // stop the service gracefully
	Restart() (string, error) // restart the service gracefully
	Status() (string, error)  // status of the service
	Install() (string, error) // install the service
}

func emptyThenDefault(value, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
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

type commonService struct {
	fakeruntime.Execer
	scriptPath string
	script     string
}

type macOSService struct {
	commonService
}

var (
	//go:embed data/macos_service.xml
	macOSServiceScript string
	//go:embed data/linux_service.txt
	linuxServiceScript string
)

func (s *macOSService) Start() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", "launchctl", "start", "com.github.linuxsuren.atest")
	return
}

func (s *macOSService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", "launchctl", "stop", "com.github.linuxsuren.atest")
	return
}

func (s *macOSService) Restart() (output string, err error) {
	if output, err = s.Stop(); err == nil {
		output, err = s.Start()
	}
	return
}

func (s *macOSService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn("sudo", "", "launchctl", "runstats", "system/com.github.linuxsuren.atest")
	return
}

func (s *macOSService) Install() (output string, err error) {
	if err = os.WriteFile(s.scriptPath, []byte(s.script), os.ModeAppend); err == nil {
		output, err = s.Execer.RunCommandAndReturn("sudo", "", "launchctl", "enable", "system/com.github.linuxsuren.atest")
	}
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
