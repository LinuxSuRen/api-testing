//go:build windows
// +build windows

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
	"os/exec"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"golang.org/x/sys/windows/svc/mgr"
)

func NewService(execer fakeruntime.Execer, scriptPath string) Service {
	return &windowsService{}
}

type windowsService struct {
}

const name = "API Testing"

func (s *windowsService) Start() (output string, err error) {
	var svc *mgr.Mgr
	if svc, err = mgr.Connect(); err != nil {
		return
	}
	defer svc.Disconnect()

	var service *mgr.Service
	service, err = svc.OpenService(name)
	if err == nil {
		err = service.Start("server")
	}
	return
}

func (s *windowsService) Stop() (output string, err error) {
	return
}

func (s *windowsService) Restart() (output string, err error) {
	return
}

func (s *windowsService) Status() (output string, err error) {
	return
}

func (s *windowsService) Install() (output string, err error) {
	var svc *mgr.Mgr
	if svc, err = mgr.Connect(); err != nil {
		return
	}
	defer svc.Disconnect()

	var service *mgr.Service
	service, err = svc.OpenService(name)
	if err == nil {
		service.Close()
		err = fmt.Errorf("service %s already exists", name)
		return
	}

	var binaryPath string
	if binaryPath, err = exec.LookPath("atest.exe"); err != nil {
		return
	}

	service, err = svc.CreateService(name, binaryPath, mgr.Config{
		StartType:   mgr.StartAutomatic,
		DisplayName: name,
	}, "server")
	if err != nil {
		return
	}
	defer service.Close()
	return
}

func (s *windowsService) Uninstall() (output string, err error) {
	var svc *mgr.Mgr
	if svc, err = mgr.Connect(); err != nil {
		return
	}
	defer svc.Disconnect()

	var service *mgr.Service
	service, err = svc.OpenService(name)
	defer service.Close()
	if err == nil {
		if err = service.Delete(); err != nil {
			return
		}
	}
	return
}

func (s *windowsService) Available() bool {
	// TODO need a way to determine if it's available
	return true
}
