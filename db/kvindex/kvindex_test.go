package kvindex

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndexMap(t *testing.T) {
	var testStruct interface{}

	// passing
	func() {
		type TestIndexStructChild struct {
			Hello string `model:"index"`
		}
		type TestIndexStruct struct {
			Foo string
			Bar TestIndexStructChild
		}

		indexEntryMap := map[string]IndexEntry{}
		testStruct = (*TestIndexStruct)(nil)
		assert.NotPanics(
			t,
			func() {
				m, indexEntryMapErr := createIndexMap(reflect.TypeOf(testStruct))
				if indexEntryMapErr != nil {
					panic(indexEntryMapErr)
				}
				indexEntryMap = m
			},
		)

		assert.Equal(t, len(indexEntryMap), 1)

		helloMap, exists := indexEntryMap["Hello"]
		if !exists {
			t.Fail()
		}
		assert.Equal(t, helloMap.Type().Kind(), reflect.String)
		assert.Equal(t, helloMap.Name(), "Hello")
		assert.Equal(t, helloMap.indexPath, []string{"Bar", "Hello"})
		assert.Equal(t, helloMap.isPrimaryKey, false)

		entity := TestIndexStruct{
			Foo: "foo",
			Bar: TestIndexStructChild{
				Hello: "helloindexed",
			},
		}

		indexKey, indexKeyErr := helloMap.ResolveIndexKeySingle(entity)
		if indexKeyErr != nil {
			t.Fail()
		}
		assert.Equal(t, "helloindexed", indexKey[0])
	}()

}
