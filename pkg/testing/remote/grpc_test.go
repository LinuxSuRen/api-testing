package remote

import (
	"testing"

	atest "github.com/linuxsuren/api-testing/pkg/testing"
	"github.com/stretchr/testify/assert"
)

func TestNewGRPCLoader(t *testing.T) {
	factory := NewGRPCloaderFromStore()

	t.Run("invalid address", func(t *testing.T) {
		writer, err := factory.NewInstance(atest.Store{})
		assert.Error(t, err)
		assert.Nil(t, writer)
	})

	t.Run("valid address", func(t *testing.T) {
		writer, err := factory.NewInstance(atest.Store{
			Kind: atest.StoreKind{
				URL: "localhost:7070",
			},
		})
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

	t.Run("NewGRPCloaderFromStore", func(t *testing.T) {
		assert.NotNil(t, NewGRPCloaderFromStore())
	})
}
