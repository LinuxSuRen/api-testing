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
package util

import (
	"fmt"
	"net/http"
)

// OKOrErrorMessage returns OK or error message
func OKOrErrorMessage(err error) string {
	return OrErrorMessage(err, "OK")
}

// OrErrorMessage returns error message or message
func OrErrorMessage(err error, message string) string {
	if err != nil {
		return err.Error()
	}
	return message
}

func ErrorWrap(err error, msg string, args ...interface{}) error {
	if err != nil {
		err = fmt.Errorf(msg, args...)
	}
	return err
}

// IgnoreErrServerClosed ignores ErrServerClosed
func IgnoreErrServerClosed(err error) error {
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
