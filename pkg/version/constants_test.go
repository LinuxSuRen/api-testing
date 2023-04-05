package version_test

import (
	"testing"

	"github.com/linuxsuren/api-testing/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver := version.GetVersion()
	assert.Empty(t, ver)
}
