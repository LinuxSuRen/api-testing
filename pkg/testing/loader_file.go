package testing

import (
	"os"
	"path"
	"path/filepath"

	"github.com/linuxsuren/api-testing/pkg/util"
)

type fileLoader struct {
	paths []string
	index int
}

// NewFileLoader creates the instance of file loader
func NewFileLoader() Loader {
	return &fileLoader{index: -1}
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
