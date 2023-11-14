//go:build windows
// +build windows

/*
MIT License
Copyright (c) 2023 API Testing Authors.
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package service

import (
	"fmt"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
	"golang.org/x/sys/windows/svc/mgr"
	"os/exec"
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
