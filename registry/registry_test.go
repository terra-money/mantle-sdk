package registry

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle/types"
)

func TestNewRegistry(t *testing.T) {
	// single entity
	func() {
		type Entity struct {
			Foo string
			Bar struct {
				Hello  uint64 `mantle:"index"`
				Mantle string `mantle:"index=custom"`
			}
		}
		indexer := func(q types.Query, c types.Commit) error {
			return nil
		}

		registry := NewRegistry([]types.IndexerRegisterer{
			func(register types.Register) {
				register(
					indexer,
					reflect.TypeOf((*Entity)(nil)),
				)
			},
		})

		assert.Equal(t, 1, len(registry.Indexers))
		assert.Equal(t, 1, len(registry.IndexerOutputs))
		assert.Equal(t, 1, len(registry.Models))
		assert.Equal(t, 2, len(registry.KVIndexMap)) // always includes BaseState

		kvi, ok := registry.KVIndexMap["Entity"]
		assert.True(t, ok)

		assert.NotNil(t, kvi.GetIndexEntry("Hello"))
		assert.NotNil(t, kvi.GetIndexEntry("custom"))
	}()

	// slice entity
	func() {
		type Entity struct {
			Foo string
			Bar struct {
				Hello  uint64 `mantle:"index"`
				Mantle string `mantle:"index=custom"`
			}
		}
		type Entities []Entity

		indexer := func(q types.Query, c types.Commit) error {
			return nil
		}

		registry := NewRegistry([]types.IndexerRegisterer{
			func(register types.Register) {
				register(
					indexer,
					reflect.TypeOf((*Entities)(nil)),
				)
			},
		})

		assert.Equal(t, 1, len(registry.Indexers))
		assert.Equal(t, 1, len(registry.IndexerOutputs))
		assert.Equal(t, 1, len(registry.Models))
		assert.Equal(t, 2, len(registry.KVIndexMap)) // always includes BaseState

		kvi, ok := registry.KVIndexMap["Entities"]
		assert.True(t, ok)

		assert.NotNil(t, kvi.GetIndexEntry("Hello"))
		assert.NotNil(t, kvi.GetIndexEntry("custom"))
	}()

	// map entity
	func() {
		type Entity struct {
			Foo string
			Bar struct {
				Hello  uint64 `mantle:"index"`
				Mantle string `mantle:"index=custom"`
			}
		}
		type Entities map[string]Entity

		indexer := func(q types.Query, c types.Commit) error {
			return nil
		}

		registry := NewRegistry([]types.IndexerRegisterer{
			func(register types.Register) {
				register(
					indexer,
					reflect.TypeOf((*Entities)(nil)),
				)
			},
		})

		assert.Equal(t, 1, len(registry.Indexers))
		assert.Equal(t, 1, len(registry.IndexerOutputs))
		assert.Equal(t, 1, len(registry.Models))
		assert.Equal(t, 2, len(registry.KVIndexMap)) // always includes BaseState

		kvi, ok := registry.KVIndexMap["Entities"]
		assert.True(t, ok)

		assert.NotNil(t, kvi.GetIndexEntry("Hello"))
		assert.NotNil(t, kvi.GetIndexEntry("custom"))
	}()

}
