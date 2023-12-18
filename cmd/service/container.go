/*
Copyright 2023 API Testing Authors.

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

package service

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/version"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

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

type ContainerOption struct {
	Image, Tag, Pull string
	Writer           io.Writer
}

const defaultImage = "linuxsuren.docker.scarf.sh/linuxsuren/api-testing"

func NewContainerService(execer fakeruntime.Execer, client string,
	featureOption ServerFeatureOption, containerOption ContainerOption) (svc Service) {
	tag := containerOption.Tag
	image := containerOption.Image
	pull := containerOption.Pull
	writer := containerOption.Writer

	if tag == "" {
		tag = "latest"
	} else if tag == version.UnknownVersion {
		tag = "master"
	} else if !strings.HasPrefix(tag, "v") {
		tag = fmt.Sprintf("v%s", tag)
	}
	if image == "" {
		image = defaultImage
	}

	containerServer := &containerService{
		Execer:       execer,
		client:       client,
		name:         ServiceName,
		image:        image,
		tag:          tag,
		pull:         pull,
		localStorage: featureOption.LocalStorage,
		secretServer: featureOption.SecretServer,
		skyWalking:   featureOption.SkyWalking,
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

func (s *containerService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(s.client, "", "stats", s.name)
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

func (s *containerService) Available() bool {
	return s.isAvailable("docker")
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
	return args
}

func (s *containerService) isAvailable(client string) bool {
	clientPath, err := s.Execer.LookPath(client)
	return err == nil && clientPath != ""
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
		err = s.Execer.RunCommandWithIO(s.client, "", s.stdOut, s.errOut, nil, "start", s.name)
	} else {
		err = s.Execer.RunCommandWithIO(s.client, "", s.stdOut, s.errOut, nil, s.getStartArgs()...)
		if err == nil {
			output, err = s.installService()
		}
	}
	return
}

func (s *podmanService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "stop", PodmanServiceName)
	return
}

func (s *podmanService) installService() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(s.client, "", "generate", "systemd", "--new", "--files", "--name", s.name)
	if err == nil {
		var result string
		result, err = s.Execer.RunCommandAndReturn("mv", "", PodmanServiceName, "/etc/systemd/system")
		if err == nil {
			output = fmt.Sprintf("%s\n%s", output, result)
			if result, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "enable", PodmanServiceName); err == nil {
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
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "disable", PodmanServiceName)
	return
}

func (s *podmanService) Available() bool {
	return s.isAvailable("podman")
}

const DefaultImage = "linuxsuren.docker.scarf.sh/linuxsuren/api-testing"
