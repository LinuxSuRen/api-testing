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
		fmt.Println(pattern, "pattern", files)
	}
	fmt.Println(l.paths, item, l.parent, err)
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
	defer func() {
		l.Reset()
	}()

	for l.HasMore() {
		var data []byte
		var loadErr error
		if data, loadErr = l.Load(); err != nil {
			fmt.Println("failed to load data", loadErr)
			continue
		}

		var testSuite *TestSuite
		if testSuite, loadErr = Parse(data); loadErr != nil {
			fmt.Println("failed to parse data", loadErr)
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
		newSuiteFile := path.Join(l.parent, fmt.Sprintf("%s.yaml", name))
		if newSuiteFile, err = filepath.Abs(newSuiteFile); err == nil {
			fmt.Println("new suite file:", newSuiteFile)

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

func (l *fileLoader) Verify() (err error) {
	// always be okay
	return
}
