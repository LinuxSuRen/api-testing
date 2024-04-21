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
package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/linuxsuren/api-testing/pkg/logging"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/linuxsuren/api-testing/pkg/util"
)

var (
	loaderLogger = logging.DefaultLogger(logging.LogLevelInfo).WithName("loader")
)

type fileLoader struct {
	paths  []string
	index  int
	parent string

	lock *sync.RWMutex
}

// NewFileLoader creates the instance of file loader
func NewFileLoader() Loader {
	return &fileLoader{index: -1, lock: &sync.RWMutex{}}
}

func NewFileWriter(parent string) Writer {
	return &fileLoader{index: -1, parent: parent, lock: &sync.RWMutex{}}
}

// HasMore returns if there are more test cases
func (l *fileLoader) HasMore() bool {
	l.index++
	return l.index < len(l.paths) && l.index >= 0
}

// Load returns the test case content
func (l *fileLoader) Load() (data []byte, err error) {
	targetFile := l.paths[l.index]
	data, err = loadData(targetFile)
	return
}

func loadData(targetFile string) (data []byte, err error) {
	if strings.HasPrefix(targetFile, "http://") || strings.HasPrefix(targetFile, "https://") {
		var ok bool
		data, ok, err = gRPCCompitableRequest(targetFile)
		if !ok && err == nil {
			var resp *http.Response
			if resp, err = http.Get(targetFile); err == nil {
				data, err = io.ReadAll(resp.Body)
			}
		}
	} else {
		data, err = os.ReadFile(targetFile)
	}
	return
}

func gRPCCompitableRequest(targetURLStr string) (data []byte, ok bool, err error) {
	if !strings.Contains(targetURLStr, "server.Runner/ConvertTestSuite") {
		return
	}

	var targetURL *url.URL
	if targetURL, err = url.Parse(targetURLStr); err != nil {
		return
	}

	suite := targetURL.Query().Get("suite")
	if suite == "" {
		err = fmt.Errorf("suite is required")
		return
	}

	payload := new(bytes.Buffer)
	payload.WriteString(fmt.Sprintf(`{"TestSuite":"%s", "Generator":"raw"}`, suite))

	var resp *http.Response
	if resp, err = http.Post(targetURLStr, "", payload); err == nil {
		if data, err = io.ReadAll(resp.Body); err != nil {
			return
		}

		var gRPCData map[string]interface{}
		if err = json.Unmarshal(data, &gRPCData); err == nil {
			var obj interface{}
			obj, ok = gRPCData["message"]
			data = []byte(fmt.Sprintf("%v", obj))
		}
	}
	return
}

// Put adds the test case path
func (l *fileLoader) Put(item string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.parent == "" {
		l.parent = path.Dir(item)
	}

	if strings.HasPrefix(item, "http://") || strings.HasPrefix(item, "https://") {
		l.paths = append(l.paths, item)
		return
	}

	for _, pattern := range util.Expand(item) {
		var files []string
		if files, err = filepath.Glob(pattern); err == nil {
			l.paths = append(l.paths, files...)
		}
		loaderLogger.Info(pattern, "pattern", len(files))
	}
	return
}

// GetContext returns the context of current test case
func (l *fileLoader) GetContext() string {
	return path.Dir(l.paths[l.index])
}

// GetCount returns the count of test cases
func (l *fileLoader) GetCount() int {
	return len(l.paths)
}

// Reset resets the index
func (l *fileLoader) Reset() {
	l.index = -1
}

func (l *fileLoader) ListTestSuite() (suites []TestSuite, err error) {
	l.lock.RLocker().Lock()
	defer l.lock.RUnlock()

	for _, target := range l.paths {
		var data []byte
		var loadErr error
		if data, loadErr = loadData(target); err != nil {
			loaderLogger.Info("failed to load data", loadErr)
			continue
		}

		var testSuite *TestSuite
		if testSuite, loadErr = Parse(data); loadErr != nil {
			loaderLogger.Info("failed to parse data", loadErr, "from", target)
			continue
		}
		suites = append(suites, *testSuite)
	}
	return
}
func (l *fileLoader) GetTestSuite(name string, full bool) (suite TestSuite, err error) {
	var items []TestSuite
	if items, err = l.ListTestSuite(); err == nil {
		for _, item := range items {
			if item.Name == name {
				suite = item
				break
			}
		}
	}
	return
}

func (l *fileLoader) CreateSuite(name, api string) (err error) {
	if name == "" {
		err = fmt.Errorf("name is required")
		return
	}

	var absPath string
	var suite *TestSuite
	if suite, absPath, err = l.GetSuite(name); err != nil {
		return
	}

	if suite != nil {
		err = fmt.Errorf("suite %s already exists", name)
	} else {
		if l.parent == "" {
			l.parent = path.Dir(absPath)
		}

		if err = os.MkdirAll(l.parent, 0755); err != nil {
			err = fmt.Errorf("failed to create %q", l.parent)
			return
		}

		newSuiteFile := path.Join(l.parent, fmt.Sprintf("%s.yaml", name))
		if newSuiteFile, err = filepath.Abs(newSuiteFile); err == nil {
			loaderLogger.Info("new suite file:", newSuiteFile)

			suite := &TestSuite{
				Name: name,
				API:  api,
			}
			if err = SaveTestSuiteToFile(suite, newSuiteFile); err == nil {
				l.Put(newSuiteFile)
			}
		}
	}
	return
}

func (l *fileLoader) GetSuite(name string) (suite *TestSuite, absPath string, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for i := range l.paths {
		suitePath := l.paths[i]
		if absPath, err = filepath.Abs(suitePath); err != nil {
			return
		}

		if suite, err = ParseTestSuiteFromFile(absPath); err != nil {
			suite = nil
			continue
		}

		if suite.Name == name {
			return
		} else {
			suite = nil
		}
	}
	return
}

// UpdateSuite updates the suite
func (l *fileLoader) UpdateSuite(suite TestSuite) (err error) {
	var absPath string
	var oldSuite *TestSuite
	if oldSuite, absPath, err = l.GetSuite(suite.Name); err == nil {
		suite.Items = oldSuite.Items // only update the suite info
		err = SaveTestSuiteToFile(&suite, absPath)
	}
	return
}

func (l *fileLoader) DeleteSuite(name string) (err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	found := false
	for i := range l.paths {
		suitePath := l.paths[i]
		var suite *TestSuite
		if suite, err = ParseTestSuiteFromFile(suitePath); err != nil {
			continue
		}

		if suite.Name == name {
			err = os.Remove(suitePath)
			l.paths = append(l.paths[:i], l.paths[i+1:]...)
			found = true
			return
		}
	}
	if !found {
		err = fmt.Errorf("suite %s not found", name)
	}
	return
}

func (l *fileLoader) ListTestCase(suite string) (testcases []TestCase, err error) {
	defer func() {
		l.Reset()
	}()

	for l.HasMore() {
		var data []byte
		if data, err = l.Load(); err != nil {
			continue
		}

		var testSuite *TestSuite
		if testSuite, err = Parse(data); err != nil {
			return
		}

		if testSuite.Name != suite {
			continue
		}

		testcases = testSuite.Items
		break
	}
	return
}

func (l *fileLoader) GetTestCase(suite, name string) (testcase TestCase, err error) {
	var items []TestCase
	if items, err = l.ListTestCase(suite); err == nil {
		found := false
		for _, item := range items {
			if item.Name == name {
				testcase = item
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("testcase %s not found", name)
		}
	}
	return
}

func (l *fileLoader) CreateTestCase(suiteName string, testcase TestCase) (err error) {
	var suite *TestSuite
	var suiteFilepath string
	for i := range l.paths {
		suitePath := l.paths[i]
		if suite, err = ParseTestSuiteFromFile(suitePath); err != nil {
			continue
		}

		if suite.Name == suiteName {
			suiteFilepath = suitePath
			break
		}
		suite = nil
	}

	if suite != nil {
		found := false
		for i := range suite.Items {
			if suite.Items[i].Name == testcase.Name {
				suite.Items[i] = testcase
				found = true
				break
			}
		}

		if !found {
			suite.Items = append(suite.Items, testcase)
		}

		err = SaveTestSuiteToFile(suite, suiteFilepath)
	}
	return
}

func (l *fileLoader) UpdateTestCase(suite string, testcase TestCase) (err error) {
	err = l.CreateTestCase(suite, testcase)
	return
}

func (l *fileLoader) DeleteTestCase(suiteName, testcase string) (err error) {
	var suite *TestSuite
	var suiteFilepath string
	for i := range l.paths {
		suitePath := l.paths[i]
		if suite, err = ParseTestSuiteFromFile(suitePath); err != nil {
			continue
		}

		if suite.Name == suiteName {
			suiteFilepath = suitePath
			break
		}
		suite = nil
	}

	if suite != nil {
		found := false
		for i := range suite.Items {
			if suite.Items[i].Name == testcase {
				suite.Items = append(suite.Items[:i], suite.Items[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			err = fmt.Errorf("testcase %s not found", testcase)
			return
		}

		err = SaveTestSuiteToFile(suite, suiteFilepath)
	}
	return
}

func (l *fileLoader) Verify() (readOnly bool, err error) {
	// always be okay
	return
}

func (l *fileLoader) PProf(string) []byte {
	// not support
	return nil
}

func (l *fileLoader) Close() {
	// not support
}
