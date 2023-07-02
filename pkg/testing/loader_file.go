package testing

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/linuxsuren/api-testing/pkg/util"
)

type fileLoader struct {
	paths  []string
	index  int
	parent string
}

// NewFileLoader creates the instance of file loader
func NewFileLoader() Loader {
	return &fileLoader{index: -1}
}

func NewFileWriter(parent string) Writer {
	return &fileLoader{index: -1, parent: parent}
}

// HasMore returns if there are more test cases
func (l *fileLoader) HasMore() bool {
	l.index++
	return l.index < len(l.paths)
}

// Load returns the test case content
func (l *fileLoader) Load() (data []byte, err error) {
	data, err = os.ReadFile(l.paths[l.index])
	return
}

// Put adds the test case path
func (l *fileLoader) Put(item string) (err error) {
	if l.parent == "" {
		l.parent = path.Dir(item)
	}

	for _, pattern := range util.Expand(item) {
		var files []string
		if files, err = filepath.Glob(pattern); err == nil {
			l.paths = append(l.paths, files...)
		}
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

func (l *fileLoader) CreateSuite(name, api string) (err error) {
	var found bool
	var parentDir string
	for i := range l.paths {
		suitePath := l.paths[i]

		var absPath string
		if absPath, err = filepath.Abs(suitePath); err != nil {
			return
		}

		parentDir = path.Dir(absPath)
		var suite *TestSuite
		if suite, err = ParseTestSuiteFromFile(absPath); err != nil {
			continue
		}

		if suite.Name == name {
			found = true
			break
		}
	}

	if found {
		err = fmt.Errorf("suite %s already exists", name)
	} else {
		if l.parent == "" {
			l.parent = parentDir
		}
		newSuiteFile := path.Join(l.parent, fmt.Sprintf("%s.yaml", name))
		fmt.Println("new suite file:", newSuiteFile)

		suite := &TestSuite{
			Name: name,
			API:  api,
		}
		if err = SaveTestSuiteToFile(suite, newSuiteFile); err == nil {
			l.Put(newSuiteFile)
		}
	}
	return
}

func (l *fileLoader) UpdateSuite(name, api string) (err error) {
	return
}

func (l *fileLoader) DeleteSuite(name string) (err error) {
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
