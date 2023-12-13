//go:build linux
// +build linux

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
