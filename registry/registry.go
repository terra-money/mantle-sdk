package registry

import (
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/types"
)

type Registry struct {
	Indexers   []types.Indexer
	Models     []types.ModelType
	KVIndexMap kvindex.KVIndexMap
}

func NewRegistry(indexRegisterers []types.IndexerRegisterer) Registry {
	registry := Registry{
		Indexers: []types.Indexer{},
		Models:   []types.ModelType{},
	}

	kvindexes := []*kvindex.KVIndex{}

	r := func(indexer types.Indexer, models ...types.ModelType) {
		registry.Indexers = append(registry.Indexers, indexer)
		for _, modelType := range models {
			registry.Models = append(registry.Models, modelType)
			kvindexes = append(kvindexes, kvindex.NewKVIndex(modelType))
		}
	}

	for _, registerer := range indexRegisterers {
		registerer(r)
	}

	registry.KVIndexMap = kvindex.NewKVIndexMap(kvindexes...)

	return registry
}
