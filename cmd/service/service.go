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

import fakeruntime "github.com/linuxsuren/go-fake-runtime"

// Service is the interface of service
type Service interface {
	Start() (string, error)     // start the service
	Stop() (string, error)      // stop the service gracefully
	Restart() (string, error)   // restart the service gracefully
	Status() (string, error)    // status of the service
	Install() (string, error)   // install the service
	Uninstall() (string, error) // uninstall the service
}

type commonService struct {
	fakeruntime.Execer
	scriptPath string
	script     string
}

func emptyThenDefault(value, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}
