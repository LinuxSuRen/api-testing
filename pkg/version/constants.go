// Package version provides the version access of this app
package version

import "fmt"

// should be injected during the build process
var version string
var date string

// GetVersion returns the version
func GetVersion() string {
	return version
}

func GetDetailedVersion() string {
	return fmt.Sprintf(`Version: %s
Date: %s`, version, date)
}
