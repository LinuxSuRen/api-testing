/*
Copyright 2024-2025 API Testing Authors.

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

package home

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetUserConfigDir() string {
	return filepath.Join(Dir(), ".config/atest")
}

func GetUserBinDir() string {
	return filepath.Join(GetUserConfigDir(), "bin")
}

func GetUserDataDir() string {
	return filepath.Join(GetUserConfigDir(), "data")
}

func GetBindingDir() string {
	return filepath.Join(GetUserConfigDir(), "data/key-binding")
}

func GetExtensionSocketPath(name string) string {
	return filepath.Join(GetUserConfigDir(), fmt.Sprintf("%s.sock", name))
}

func getCommonHomeDir() string {
	return os.Getenv("HOME")
}

func getHomeDirViaShell() string {
	var stdout bytes.Buffer
	stdout.Reset()
	cmd := exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return ""
	}

	return strings.TrimSpace(stdout.String())
}
