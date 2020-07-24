package registry

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/types"
)

type TestModel1 struct {
	Test int
}

type TestModel2 struct {
	Test int
}

type TestModel3 struct {
	Test struct {
		Whatever string `mantle:"index=whatever"`
	}
}

func TestNewRegistry(t *testing.T) {
	indexer1 := func(q types.Query, c types.Commit) {}
	indexer2 := func(q types.Query, c types.Commit) {}

	models := []interface{}{
		(*TestModel1)(nil),
		(*TestModel2)(nil),
		(*TestModel3)(nil),
	}

	testIndexers := []types.IndexerRegisterer{
		func(register types.Register) {
			register(
				indexer1,
				reflect.TypeOf(models[0]),
			)
		},
		func(register types.Register) {
			register(
				indexer2,
				reflect.TypeOf(models[1]),
				reflect.TypeOf(models[2]),
			)
		},
	}

	registry := NewRegistry(testIndexers)

	assert.Equal(t, len(registry.Indexers), 2)
	assert.Equal(t, len(registry.Models), 3)
	assert.Equal(t, len(registry.KVIndexMap), 3)

	// check models
	for i, modelType := range registry.Models {
		assert.Equal(t, modelType, reflect.TypeOf(models[i]))
	}

	// check kvindexMap
	// testmodel1,2 don't have any indexes so internals are empty
	assert.Equal(t, registry.KVIndexMap["TestModel1"], kvindex.NewKVIndex(reflect.TypeOf(models[0])))
	assert.False(t, registry.KVIndexMap["TestModel1"].HasIndex())
	assert.Equal(t, registry.KVIndexMap["TestModel2"], kvindex.NewKVIndex(reflect.TypeOf(models[1])))
	assert.False(t, registry.KVIndexMap["TestModel2"].HasIndex())

	tm3r, err := registry.KVIndexMap["TestModel3"].BuildIndexKey("whatever", "William", []byte("terra"))
	assert.Nil(t, err)
	tm3a, err := kvindex.NewKVIndex(reflect.TypeOf(models[2])).BuildIndexKey("whatever", "William", []byte("terra"))
	assert.Nil(t, err)
	// testmodel3 has index and
	assert.Equal(t, tm3r, tm3a)

}
