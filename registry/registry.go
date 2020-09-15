package registry

import (
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/types"
	"reflect"
)

type Registry struct {
	Indexers       []types.Indexer
	IndexerOutputs [][]types.Model
	Models         []types.Model
	KVIndexMap     kvindex.KVIndexMap
}

func NewRegistry(indexRegisterers []types.IndexerRegisterer) Registry {
	registry := Registry{
		Indexers: []types.Indexer{},
		Models:   []types.Model{},
	}

	// add BaseState to kvindexes
	baseStateKVIndex, baseStateKVIndexErr := kvindex.NewKVIndex(reflect.TypeOf(types.BaseState{}))
	if baseStateKVIndexErr != nil {
		panic(baseStateKVIndexErr)
	}
	kvindexes := []*kvindex.KVIndex{
		baseStateKVIndex,
	}

	r := func(indexer types.Indexer, models ...types.Model) {
		registry.Indexers = append(registry.Indexers, indexer)
		registry.IndexerOutputs = append(registry.IndexerOutputs, models)

		for _, model := range models {
			registry.Models = append(registry.Models, model)
			kvIndex, kvIndexErr := kvindex.NewKVIndex(model)
			if kvIndexErr != nil {
				panic(kvIndexErr)
			}
			kvindexes = append(kvindexes, kvIndex)
		}
	}

	for _, registerer := range indexRegisterers {
		registerer(r)
	}

	registry.KVIndexMap = kvindex.NewKVIndexMap(kvindexes...)

	return registry
}
