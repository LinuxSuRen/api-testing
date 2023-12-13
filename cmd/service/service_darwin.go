//go:build darwin
// +build darwin

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

	"github.com/linuxsuren/api-testing/pkg/util"
	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

type macOSService struct {
	commonService
	cli string
	id  string
}

func NewService(execer fakeruntime.Execer, scriptPath string) Service {
	return &macOSService{
		commonService: commonService{
			Execer:     execer,
			scriptPath: util.EmptyThenDefault(scriptPath, "/Library/LaunchDaemons/com.github.linuxsuren.atest.plist"),
			script:     macOSServiceScript,
		},
	}
}

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

func (s *macOSService) Available() bool {
	// TODO need a way to determine if it's available
	return true
}

//go:embed data/macos_service.xml
var macOSServiceScript string
