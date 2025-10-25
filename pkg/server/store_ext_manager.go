/*
Copyright 2023-2025 API Testing Authors.

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
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	sync "sync"
	"syscall"
	"time"

	"github.com/linuxsuren/api-testing/pkg/testing"
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
	lock                 *sync.RWMutex
}

var ss *storeExtManager

func NewStoreExtManager(execer fakeruntime.Execer) ExtManager {
	if ss == nil {
		ss = &storeExtManager{
			processChan: make(chan fakeruntime.Process),
			stopSingal:  make(chan struct{}, 1),
			lock:        &sync.RWMutex{},
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
		lock:        &sync.RWMutex{},
	}
	ss.execer = execer
	ss.socketPrefix = "unix://"
	ss.extStatusMap = map[string]bool{}
	ss.processCollect()
	ss.WithDownloader(&nonDownloader{})
	return ss
}

func (s *storeExtManager) StartPlugin(storeKind testing.StoreKind) {
	for _, plugin := range storeKind.Dependencies {
		s.Start(plugin.Name, fmt.Sprintf("unix://%s", home.GetExtensionSocketPath(plugin.Name)))
	}
}

func (s *storeExtManager) Start(name, socket string) (err error) {
	if name == "" {
	}

	serverLogger.Info("start", "extension", name, "socket", socket)
	if v, ok := s.extStatusMap[name]; ok && v {
		return
	}

	platformBasedName := name
	if s.execer.OS() == "windows" {
		platformBasedName += ".exe"
	} else {
		socket = fmt.Sprintf("unix://%s", home.GetExtensionSocketPath(name))
	}

	targetDir := home.GetUserBinDir()
	targetBinaryFile := filepath.Join(targetDir, platformBasedName)

	var binaryPath string
	if _, err = os.Stat(targetBinaryFile); err == nil {
		binaryPath = targetBinaryFile
	} else {
		serverLogger.Info("failed to find extension", "error", err.Error())

		binaryPath, err = s.execer.LookPath(platformBasedName)
		if err != nil {
			err = fmt.Errorf("not found extension, try to download it, error: %v", err)
			go func() {
				ociDownloader := downloader.NewStoreDownloader()
				ociDownloader.WithKind("store")
				ociDownloader.WithOS(s.execer.OS())
				reader, dErr := ociDownloader.Download(name, "", "")
				if dErr != nil {
					serverLogger.Error(dErr, "failed to download extension", "name", name)
				} else {
					extFile := ociDownloader.GetTargetFile(name)

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
	if strings.Contains(socketURL, ":") && !strings.HasPrefix(socketURL, s.socketPrefix) {
		err = s.startPluginViaHTTP(socketURL, plugin, pluginName)
		return
	}
	socketFile := strings.TrimPrefix(socketURL, s.socketPrefix)
	_ = os.RemoveAll(socketFile) // always deleting the socket file to avoid start failing

	s.lock.Lock()
	s.filesNeedToBeRemoved = append(s.filesNeedToBeRemoved, socketFile)
	s.extStatusMap[pluginName] = true
	s.lock.Unlock()

	if err = s.execer.RunCommandWithIO(plugin, "", os.Stdout, os.Stderr, s.processChan, "--socket", socketFile); err != nil {
		serverLogger.Info("failed to start ext manager", "socket", socketURL, "error: ", err.Error())
	}
	return
}

func (s *storeExtManager) startPluginViaHTTP(httpURL, plugin, pluginName string) (err error) {
	port := strings.Split(httpURL, ":")[1]
	if err = s.execer.RunCommandWithIO(plugin, "", os.Stdout, os.Stderr, s.processChan, "--port", port); err != nil {
		serverLogger.Info("failed to start ext manager", "port", port, "error: ", err.Error())
	}
	return
}

func (s *storeExtManager) StopAll() error {
	serverLogger.Info("stop", "extensions", len(s.processs))
	for _, p := range s.processs {
		if p != nil {
			// Use Kill on Windows, Signal on other platforms
			if isWindows() {
				p.Kill()
			} else {
				p.Signal(syscall.SIGTERM)
			}
		}
	}
	s.stopSingal <- struct{}{}
	return nil
}

// isWindows returns true if the program is running on Windows OS.
func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") ||
		(strings.Contains(strings.ToLower(os.Getenv("GOOS")), "windows"))
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

var _ downloader.PlatformAwareOCIDownloader = &nonDownloader{}

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

func (n *nonDownloader) WithKind(string) {
	// Do nothing because this is an empty implementation
}

func (n *nonDownloader) WithImagePrefix(imagePrefix string) {
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

func (n *nonDownloader) GetTargetFile(string) string {
	return ""
}

func (n *nonDownloader) WithOptions(opts ...downloader.OICDownloaderOption) {
	// Do nothing because this is an empty implementation
}
