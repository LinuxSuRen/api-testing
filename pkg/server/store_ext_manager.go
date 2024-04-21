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
package server

import (
	"fmt"
	"github.com/linuxsuren/api-testing/pkg/logging"
	"os"
	"strings"
	"syscall"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

var (
	serverLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("server")
)

type ExtManager interface {
	Start(name, socket string) (err error)
	StopAll() (err error)
}

type storeExtManager struct {
	execer               fakeruntime.Execer
	socketPrefix         string
	filesNeedToBeRemoved []string
	extStatusMap         map[string]bool
	processs             []fakeruntime.Process
	processChan          chan fakeruntime.Process
	stopSingal           chan struct{}
}

var s *storeExtManager

func NewStoreExtManager(execer fakeruntime.Execer) ExtManager {
	if s == nil {
		s = &storeExtManager{
			processChan: make(chan fakeruntime.Process, 0),
			stopSingal:  make(chan struct{}, 1),
		}
		s.execer = execer
		s.socketPrefix = "unix://"
		s.extStatusMap = map[string]bool{}
		s.processCollect()
	}
	return s
}

func (s *storeExtManager) Start(name, socket string) (err error) {
	if v, ok := s.extStatusMap[name]; ok && v {
		return
	}

	binaryPath, lookErr := s.execer.LookPath(name)
	if lookErr != nil {
		err = fmt.Errorf("failed to find %s, error: %v", name, lookErr)
	} else {
		go func(socketURL, plugin string) {
			socketFile := strings.TrimPrefix(socketURL, s.socketPrefix)
			s.filesNeedToBeRemoved = append(s.filesNeedToBeRemoved, socketFile)
			s.extStatusMap[name] = true
			if err = s.execer.RunCommandWithIO(plugin, "", os.Stdout, os.Stderr, s.processChan, "--socket", socketFile); err != nil {
				serverLogger.Info("failed to start %s, error: %v", socketURL, err)
			}
		}(socket, binaryPath)
	}
	return
}

func (s *storeExtManager) StopAll() error {
	serverLogger.Info("stop", len(s.processs), "extensions")
	for _, p := range s.processs {
		if p != nil {
			p.Signal(syscall.SIGTERM)
		}
	}
	s.stopSingal <- struct{}{}
	return nil
}

func (s *storeExtManager) processCollect() {
	go func() {
		for {
			select {
			case p := <-s.processChan:
				s.processs = append(s.processs, p)
			case <-s.stopSingal:
				return
			}
		}
	}()
}
