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
// Package version provides the version access of this app
package version

import "fmt"

const UnknownVersion = "unknown"

// should be injected during the build process
var version string
var date string
var commit string

// GetVersion returns the version
func GetVersion() string {
	if version == "" {
		return UnknownVersion
	}
	return version
}

func GetDate() string {
	return date
}

func GetCommit() string {
	return commit
}

func GetDetailedVersion() string {
	return fmt.Sprintf(`Version: %s
Date: %s`, version, date)
}
