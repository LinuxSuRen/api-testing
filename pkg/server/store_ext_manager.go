/*
Copyright 2023-2024 API Testing Authors.

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
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/linuxsuren/api-testing/pkg/util/home"

	"github.com/linuxsuren/api-testing/pkg/downloader"
	"github.com/linuxsuren/api-testing/pkg/logging"

	fakeruntime "github.com/linuxsuren/go-fake-runtime"
)

var (
	serverLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("server")
)

type ExtManager interface {
	Start(name, socket string) (err error)
	StopAll() (err error)
	WithDownloader(downloader.PlatformAwareOCIDownloader)
}

type storeExtManager struct {
	execer               fakeruntime.Execer
	ociDownloader        downloader.PlatformAwareOCIDownloader
	socketPrefix         string
	filesNeedToBeRemoved []string
	extStatusMap         map[string]bool
	processs             []fakeruntime.Process
	processChan          chan fakeruntime.Process
	stopSingal           chan struct{}
}

var ss *storeExtManager

func NewStoreExtManager(execer fakeruntime.Execer) ExtManager {
	if ss == nil {
		ss = &storeExtManager{
			processChan: make(chan fakeruntime.Process),
			stopSingal:  make(chan struct{}, 1),
		}
		ss.execer = execer
		ss.socketPrefix = "unix://"
		ss.extStatusMap = map[string]bool{}
		ss.processCollect()
		ss.WithDownloader(&nonDownloader{})
	}
	return ss
}

func NewStoreExtManagerInstance(execer fakeruntime.Execer) ExtManager {
	ss = &storeExtManager{
		processChan: make(chan fakeruntime.Process),
		stopSingal:  make(chan struct{}, 1),
	}
	ss.execer = execer
	ss.socketPrefix = "unix://"
	ss.extStatusMap = map[string]bool{}
	ss.processCollect()
	ss.WithDownloader(&nonDownloader{})
	return ss
}

func (s *storeExtManager) Start(name, socket string) (err error) {
	if v, ok := s.extStatusMap[name]; ok && v {
		return
	}
	targetDir := home.GetUserBinDir()
	targetBinaryFile := filepath.Join(targetDir, name)

	var binaryPath string
	if _, err = os.Stat(targetBinaryFile); err == nil {
		binaryPath = targetBinaryFile
	} else {
		binaryPath, err = s.execer.LookPath(name)
		if err != nil {
			go func() {
				reader, dErr := s.ociDownloader.Download(name, "", "")
				if dErr != nil {
					serverLogger.Error(dErr, "failed to download extension", "name", name)
				} else {
					extFile := s.ociDownloader.GetTargetFile()

					targetFile := filepath.Base(extFile)
					if dErr = downloader.WriteTo(reader, targetDir, targetFile); dErr == nil {
						binaryPath = filepath.Join(targetDir, targetFile)
						s.startPlugin(socket, binaryPath, name)
					} else {
						serverLogger.Error(dErr, "failed to save extension", "targetFile", targetFile)
					}
				}
			}()
		}
	}

	if err == nil {
		go s.startPlugin(socket, binaryPath, name)
	}
	return
}

func (s *storeExtManager) startPlugin(socketURL, plugin, pluginName string) (err error) {
	socketFile := strings.TrimPrefix(socketURL, s.socketPrefix)
	_ = os.RemoveAll(socketFile) // always deleting the socket file to avoid start failing

	s.filesNeedToBeRemoved = append(s.filesNeedToBeRemoved, socketFile)
	s.extStatusMap[pluginName] = true
	if err = s.execer.RunCommandWithIO(plugin, "", os.Stdout, os.Stderr, s.processChan, "--socket", socketFile); err != nil {
		serverLogger.Info("failed to start: ", socketURL, "error: ", err.Error())
	}
	return
}

func (s *storeExtManager) StopAll() error {
	serverLogger.Info("stop", "extensions", len(s.processs))
	for _, p := range s.processs {
		if p != nil {
			p.Signal(syscall.SIGTERM)
		}
	}
	s.stopSingal <- struct{}{}
	return nil
}

func (s *storeExtManager) WithDownloader(ociDownloader downloader.PlatformAwareOCIDownloader) {
	s.ociDownloader = ociDownloader
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

var ErrDownloadNotSupport = errors.New("no support")

type nonDownloader struct{}

func (n *nonDownloader) WithBasicAuth(username string, password string) {
	// Do nothing because this is an empty implementation
}

func (n *nonDownloader) Download(image, tag, file string) (reader io.Reader, err error) {
	err = ErrDownloadNotSupport
	return
}

func (n *nonDownloader) WithOS(string) {
	// Do nothing because this is an empty implementation
}

func (n *nonDownloader) WithArch(string) {
	// Do nothing because this is an empty implementation
}

func (n *nonDownloader) WithRegistry(string) {
	// Do nothing because this is an empty implementation
}

func (d *nonDownloader) WithRoundTripper(rt http.RoundTripper) {
	// Do nothing because this is an empty implementation
}

func (d *nonDownloader) WithInsecure(bool) {
	// Do nothing because this is an empty implementation
}

func (d *nonDownloader) WithTimeout(time.Duration)   {}
func (d *nonDownloader) WithContext(context.Context) {}

func (n *nonDownloader) GetTargetFile() string {
	return ""
}
