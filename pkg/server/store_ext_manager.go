/**
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

package server

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
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
				log.Printf("failed to start %s, error: %v", socketURL, err)
			}
		}(socket, binaryPath)
	}
	return
}

func (s *storeExtManager) StopAll() error {
	log.Println("stop", len(s.processs), "extensions")
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
