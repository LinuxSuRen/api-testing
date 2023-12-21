//go:build linux
// +build linux

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
	_ "embed"
	"os"
	"strings"

	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

func NewService(execer fakeruntime.Execer, scriptPath string) Service {
	return &linuxService{
		commonService: commonService{
			Execer:     execer,
			scriptPath: util.EmptyThenDefault(scriptPath, "/lib/systemd/system/atest.service"),
			script:     linuxServiceScript,
		},
	}
}

type linuxService struct {
	commonService
}

func (s *linuxService) Start() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "start", ServiceName)
	return
}

func (s *linuxService) Stop() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "stop", ServiceName)
	return
}

func (s *linuxService) Restart() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "restart", ServiceName)
	return
}

func (s *linuxService) Status() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "status", ServiceName)
	if err != nil && err.Error() == "exit status 3" {
		// this is normal case
		err = nil
	}
	return
}

func (s *linuxService) Install() (output string, err error) {
	if err = os.WriteFile(s.scriptPath, []byte(s.script), os.ModeAppend); err == nil {
		output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "enable", ServiceName)
	}
	return
}

func (s *linuxService) Uninstall() (output string, err error) {
	output, err = s.Execer.RunCommandAndReturn(SystemCtl, "", "disable", ServiceName)
	return
}

func (s *linuxService) Available() bool {
	output, err := s.Execer.RunCommandAndReturn("systemctl", "", "is-system-running")
	output = strings.TrimSpace(output)
	return err == nil && output != "offline"
}

//go:embed data/linux_service.txt
var linuxServiceScript string
