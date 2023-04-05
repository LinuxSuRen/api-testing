// Package version provides the version access of this app
package version

// should be injected during the build process
var version string

// GetVersion returns the version
func GetVersion() string {
	return version
}
