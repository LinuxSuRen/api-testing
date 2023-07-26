package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGRPCLoader(t *testing.T) {
	t.Run("invalid address", func(t *testing.T) {
		writer, err := NewGRPCLoader("")
		assert.Error(t, err)
		assert.Nil(t, writer)
	})

	t.Run("valid address", func(t *testing.T) {
		writer, err := NewGRPCLoader("localhost:8907")
		assert.NoError(t, err)
		assert.NotNil(t, writer)

		assert.False(t, writer.HasMore())
		var data []byte
		data, err = writer.Load()
		assert.NoError(t, err)
		assert.Nil(t, data)

		assert.NoError(t, writer.Put("foo"))

		assert.Empty(t, writer.GetContext())

		assert.Equal(t, 0, writer.GetCount())
		writer.Reset()
	})
}
