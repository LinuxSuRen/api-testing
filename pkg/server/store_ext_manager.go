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

// AIPluginInfo represents information about an AI plugin
type AIPluginInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Capabilities []string          `json:"capabilities"`
	SocketPath   string            `json:"socketPath"`
	Metadata     map[string]string `json:"metadata"`
}

// AIPluginHealth represents the health status of an AI plugin
type AIPluginHealth struct {
	Name         string            `json:"name"`
	Status       string            `json:"status"` // online, offline, error, processing
	LastCheckAt  time.Time         `json:"lastCheckAt"`
	ResponseTime time.Duration     `json:"responseTime"`
	ErrorMessage string            `json:"errorMessage,omitempty"`
	Metrics      map[string]string `json:"metrics,omitempty"`
}

// ExtManager handles general extension management (start, stop, download)
type ExtManager interface {
	Start(name, socket string) (err error)
	StopAll() (err error)
	WithDownloader(downloader.PlatformAwareOCIDownloader)
}

// AIPluginManager handles AI-specific plugin management
type AIPluginManager interface {
	DiscoverAIPlugins() ([]AIPluginInfo, error)
	CheckAIPluginHealth(name string) (*AIPluginHealth, error)
	GetAllAIPluginHealth() (map[string]*AIPluginHealth, error)
	RegisterAIPlugin(info AIPluginInfo) error
	UnregisterAIPlugin(name string) error
}

// CompositeManager combines both extension and AI plugin management
type CompositeManager interface {
	ExtManager
	AIPluginManager
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
	// AI Plugin Management
	aiPluginRegistry   map[string]AIPluginInfo    `json:"aiPluginRegistry"`
	aiPluginHealthMap  map[string]*AIPluginHealth `json:"aiPluginHealthMap"`
	healthCheckTicker  *time.Ticker               `json:"-"`
	healthCheckCtx     context.Context            `json:"-"`
	healthCheckCancel  context.CancelFunc         `json:"-"`
}

var ss *storeExtManager

func NewStoreExtManager(execer fakeruntime.Execer) CompositeManager {
	if ss == nil {
		ctx, cancel := context.WithCancel(context.Background())
		ss = &storeExtManager{
			processChan: make(chan fakeruntime.Process),
			stopSingal:  make(chan struct{}, 1),
			lock:        &sync.RWMutex{},
			// AI Plugin Management initialization
			aiPluginRegistry:  make(map[string]AIPluginInfo),
			aiPluginHealthMap: make(map[string]*AIPluginHealth),
			healthCheckCtx:    ctx,
			healthCheckCancel: cancel,
		}
		ss.execer = execer
		ss.socketPrefix = "unix://"
		ss.extStatusMap = map[string]bool{}
		ss.processCollect()
		ss.WithDownloader(&nonDownloader{})
		// Start AI plugin health monitoring
		ss.startAIHealthMonitoring()
	}
	return ss
}

func NewStoreExtManagerInstance(execer fakeruntime.Execer) CompositeManager {
	ctx, cancel := context.WithCancel(context.Background())
	ss = &storeExtManager{
		processChan: make(chan fakeruntime.Process),
		stopSingal:  make(chan struct{}, 1),
		lock:        &sync.RWMutex{},
		// AI Plugin Management initialization
		aiPluginRegistry:  make(map[string]AIPluginInfo),
		aiPluginHealthMap: make(map[string]*AIPluginHealth),
		healthCheckCtx:    ctx,
		healthCheckCancel: cancel,
	}
	ss.execer = execer
	ss.socketPrefix = "unix://"
	ss.extStatusMap = map[string]bool{}
	ss.processCollect()
	ss.WithDownloader(&nonDownloader{})
	// Start AI plugin health monitoring
	ss.startAIHealthMonitoring()
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
					extFile := ociDownloader.GetTargetFile()

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

func (n *nonDownloader) GetTargetFile() string {
	return ""
}

// AI Plugin Management Implementation

// startAIHealthMonitoring starts the periodic health check for AI plugins
func (s *storeExtManager) startAIHealthMonitoring() {
	s.healthCheckTicker = time.NewTicker(30 * time.Second) // Health check every 30 seconds
	
	go func() {
		for {
			select {
			case <-s.healthCheckCtx.Done():
				s.healthCheckTicker.Stop()
				return
			case <-s.healthCheckTicker.C:
				s.performHealthCheck()
			}
		}
	}()
}

// performHealthCheck performs health checks on all registered AI plugins
func (s *storeExtManager) performHealthCheck() {
	s.lock.RLock()
	plugins := make(map[string]AIPluginInfo)
	for name, info := range s.aiPluginRegistry {
		plugins[name] = info
	}
	s.lock.RUnlock()

	for name, info := range plugins {
		health, err := s.checkSingleAIPlugin(info)
		if err != nil {
			serverLogger.Error(err, "Failed to check AI plugin health", "plugin", name)
			health = &AIPluginHealth{
				Name:         name,
				Status:       "error",
				LastCheckAt:  time.Now(),
				ErrorMessage: err.Error(),
			}
		}

		s.lock.Lock()
		s.aiPluginHealthMap[name] = health
		s.lock.Unlock()
	}
}

// checkSingleAIPlugin performs health check on a single AI plugin
func (s *storeExtManager) checkSingleAIPlugin(info AIPluginInfo) (*AIPluginHealth, error) {
	startTime := time.Now()
	
	// For now, we'll simulate a health check by checking if the socket file exists
	// In a real implementation, this would make a gRPC health check call
	_, err := os.Stat(strings.TrimPrefix(info.SocketPath, "unix://"))
	
	responseTime := time.Since(startTime)
	
	health := &AIPluginHealth{
		Name:         info.Name,
		LastCheckAt:  time.Now(),
		ResponseTime: responseTime,
		Metrics: map[string]string{
			"version":     info.Version,
			"socket_path": info.SocketPath,
		},
	}

	if err != nil {
		if os.IsNotExist(err) {
			health.Status = "offline"
			health.ErrorMessage = "Plugin socket not found"
		} else {
			health.Status = "error"
			health.ErrorMessage = err.Error()
		}
	} else {
		health.Status = "online"
		health.ErrorMessage = ""
	}

	return health, nil
}

// DiscoverAIPlugins discovers AI-capable plugins in the system
func (s *storeExtManager) DiscoverAIPlugins() ([]AIPluginInfo, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	
	var plugins []AIPluginInfo
	for _, info := range s.aiPluginRegistry {
		plugins = append(plugins, info)
	}
	
	return plugins, nil
}

// CheckAIPluginHealth checks the health of a specific AI plugin
func (s *storeExtManager) CheckAIPluginHealth(name string) (*AIPluginHealth, error) {
	s.lock.RLock()
	info, exists := s.aiPluginRegistry[name]
	s.lock.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("AI plugin %s not found", name)
	}
	
	health, err := s.checkSingleAIPlugin(info)
	if err != nil {
		return nil, fmt.Errorf("failed to check health for AI plugin %s: %w", name, err)
	}
	
	// Update the health cache
	s.lock.Lock()
	s.aiPluginHealthMap[name] = health
	s.lock.Unlock()
	
	return health, nil
}

// GetAllAIPluginHealth returns the health status of all AI plugins
func (s *storeExtManager) GetAllAIPluginHealth() (map[string]*AIPluginHealth, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	
	// Return a copy to avoid concurrent access issues
	healthMap := make(map[string]*AIPluginHealth)
	for name, health := range s.aiPluginHealthMap {
		// Create a copy of the health struct
		healthCopy := &AIPluginHealth{
			Name:         health.Name,
			Status:       health.Status,
			LastCheckAt:  health.LastCheckAt,
			ResponseTime: health.ResponseTime,
			ErrorMessage: health.ErrorMessage,
			Metrics:      make(map[string]string),
		}
		
		// Copy metrics map
		for k, v := range health.Metrics {
			healthCopy.Metrics[k] = v
		}
		
		healthMap[name] = healthCopy
	}
	
	return healthMap, nil
}

// RegisterAIPlugin registers a new AI plugin with the system
func (s *storeExtManager) RegisterAIPlugin(info AIPluginInfo) error {
	if info.Name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}
	
	if info.SocketPath == "" {
		return fmt.Errorf("plugin socket path cannot be empty")
	}
	
	s.lock.Lock()
	defer s.lock.Unlock()
	
	// Check if plugin is already registered
	if _, exists := s.aiPluginRegistry[info.Name]; exists {
		serverLogger.Info("AI plugin already registered, updating info", "plugin", info.Name)
	}
	
	s.aiPluginRegistry[info.Name] = info
	
	// Initialize health status
	s.aiPluginHealthMap[info.Name] = &AIPluginHealth{
		Name:        info.Name,
		Status:      "unknown",
		LastCheckAt: time.Now(),
		Metrics: map[string]string{
			"version":     info.Version,
			"socket_path": info.SocketPath,
		},
	}
	
	serverLogger.Info("AI plugin registered successfully", "plugin", info.Name, "version", info.Version)
	
	return nil
}

// UnregisterAIPlugin removes an AI plugin from the system
func (s *storeExtManager) UnregisterAIPlugin(name string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	
	if _, exists := s.aiPluginRegistry[name]; !exists {
		return fmt.Errorf("AI plugin %s not found", name)
	}
	
	delete(s.aiPluginRegistry, name)
	delete(s.aiPluginHealthMap, name)
	
	serverLogger.Info("AI plugin unregistered successfully", "plugin", name)
	
	return nil
}

// StopAll enhanced to also clean up AI plugin monitoring
func (s *storeExtManager) StopAll() error {
	// Stop AI health monitoring
	if s.healthCheckCancel != nil {
		s.healthCheckCancel()
	}
	
	// Original StopAll implementation
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

	for _, fileToRemove := range s.filesNeedToBeRemoved {
		if err := os.RemoveAll(fileToRemove); err != nil {
			serverLogger.Info("failed to remove", "file", fileToRemove, "error", err)
		}
	}
	
	// Send stop signal
	s.stopSingal <- struct{}{}
	
	return nil
}
