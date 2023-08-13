package compare

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNoEqualErr(t *testing.T) {
	err := newNoEqualErr("data", fmt.Errorf("this is msg"))
	err = newNoEqualErr("to", err)
	err = newNoEqualErr("path", err)
	assert.Equal(t, "compare: field path.to.data: this is msg", err.Error())
}
