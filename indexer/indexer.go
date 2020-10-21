package indexer

import (
	"sync"

	"github.com/terra-project/mantle-sdk/graph"
	"github.com/terra-project/mantle-sdk/graph/generate"
	"github.com/terra-project/mantle-sdk/types"
)

type IndexerBaseInstance struct {
	indexers       []types.Indexer
	indexerOutputs [][]types.Model
	committer      types.GraphQLCommitter
	querier        types.GraphQLQuerier
}

func NewIndexerBaseInstance(
	indexers []types.Indexer,
	indexerOutputs [][]types.Model,
	querier types.GraphQLQuerier,
	committer types.GraphQLCommitter,
) *IndexerBaseInstance {
	return &IndexerBaseInstance{
		indexers:       indexers,
		indexerOutputs: indexerOutputs,
		committer:      committer,
		querier:        querier,
	}
}

func (instance *IndexerBaseInstance) RunIndexerRound() {
	// create wait group for ALL indexers
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(instance.indexers))

	for i, indexer := range instance.indexers {

		// indexer outputs are necessary for detecting self-reference
		indexerOutput := instance.indexerOutputs[i]

		runIndexer(
			&waitGroup,
			instance.committer,
			instance.querier,
			indexer,
			indexerOutput,
		)
	}

	waitGroup.Wait()
}

func runIndexer(
	wg *sync.WaitGroup,
	committer types.GraphQLCommitter,
	querier types.GraphQLQuerier,
	indexer types.Indexer,
	indexerOutput []types.Model,
) {
	var isolatedQuerier = createIsolatedQuerier(querier, indexerOutput)
	var isolatedCommitter = createIsolatedCommitter(committer)

	go func() {
		defer wg.Done()
		if indexerErr := indexer(isolatedQuerier, isolatedCommitter); indexerErr != nil {
			panic(indexerErr)
		}
	}()
}

func createIsolatedQuerier(
	querier types.GraphQLQuerier,
	indexerSelfOutput []types.Model,
) types.IndexerQuerier {
	return func(query interface{}, variables types.GraphQLParams) error {
		qs := generate.GenerateQuery(query, variables)
		result := querier(qs, variables, indexerSelfOutput)
		resultInternal := result.(*types.GraphQLInternalResult)

		return graph.UnmarshalInternalQueryResult(resultInternal, query)
	}
}

func createIsolatedCommitter(committer types.GraphQLCommitter) types.IndexerCommitter {
	return func(entity interface{}) error {
		return committer(entity)
	}
}
