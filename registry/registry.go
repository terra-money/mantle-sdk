package registry

import (
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
	"reflect"
)

type Registry struct {
	Indexers       []types.Indexer
	IndexerOutputs [][]types.ModelType
	Models         []types.ModelType
	KVIndexMap     kvindex.KVIndexMap
}

func NewRegistry(indexRegisterers []types.IndexerRegisterer) Registry {
	registry := Registry{
		Indexers:       []types.Indexer{},
		IndexerOutputs: [][]types.ModelType{},
		Models:         []types.ModelType{},
	}

	// add BaseState to kvindexes
	kvindexes := []*kvindex.KVIndex{
		kvindex.NewKVIndex(reflect.TypeOf(types.BaseState{})),
	}

	r := func(indexer types.Indexer, models ...types.ModelType) {
		var actualModels []types.ModelType
		for _, model := range models {
			actualModels = append(actualModels, utils.GetType(model))
		}

		registry.Indexers = append(registry.Indexers, indexer)
		registry.IndexerOutputs = append(registry.IndexerOutputs, actualModels)

		for _, modelType := range actualModels {
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
