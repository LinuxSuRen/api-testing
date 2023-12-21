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
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

// Service is the interface of service
type Service interface {
	Start() (string, error)     // start the service
	Stop() (string, error)      // stop the service gracefully
	Restart() (string, error)   // restart the service gracefully
	Status() (string, error)    // status of the service
	Install() (string, error)   // install the service
	Uninstall() (string, error) // uninstall the service
	Available() bool
}

type commonService struct {
	fakeruntime.Execer
	scriptPath string
	script     string
}

type ServiceMode string

const (
	ServiceModeOS        ServiceMode = "os"
	ServiceModeContainer ServiceMode = "container"
	ServiceModePodman    ServiceMode = "podman"
	ServiceModeDocker    ServiceMode = "docker"
)

func (s ServiceMode) All() []ServiceMode {
	return []ServiceMode{ServiceModeOS, ServiceModeContainer,
		ServiceModePodman, ServiceModeDocker}
}

func (s ServiceMode) String() string {
	return string(s)
}

type ServerFeatureOption struct {
	SecretServer string
	SkyWalking   string
	LocalStorage string
}

func GetAvailableService(mode ServiceMode, execer fakeruntime.Execer,
	containerOption ContainerOption, featureOption ServerFeatureOption,
	scriptPath string) (svc Service) {
	osService := NewService(execer, scriptPath)
	dockerService := NewContainerService(execer, "docker", featureOption, containerOption)
	podmanService := NewContainerService(execer, "podman", featureOption, containerOption)

	switch mode {
	case ServiceModeOS:
		svc = osService
	case ServiceModeDocker, ServiceModeContainer:
		svc = dockerService
	case ServiceModePodman:
		svc = podmanService
	default:
		if osService.Available() {
			svc = osService
		} else if dockerService.Available() {
			svc = dockerService
		} else if podmanService.Available() {
			svc = podmanService
		}
	}

	if svc != nil && !svc.Available() {
		svc = nil
	}
	return
}
