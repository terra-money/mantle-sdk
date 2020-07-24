package indexer

import (
	"context"
	"fmt"
	"sync"

	"github.com/terra-project/mantle/graph"
	"github.com/terra-project/mantle/graph/generate"
	"github.com/terra-project/mantle/types"
)

type IndexerBaseInstance struct {
	indexers  []types.Indexer
	committer types.GraphQLCommitter
	querier   types.GraphQLQuerier
}

func NewIndexerBaseInstance(
	indexers []types.Indexer,
	querier types.GraphQLQuerier,
	committer types.GraphQLCommitter,
) *IndexerBaseInstance {
	return &IndexerBaseInstance{
		indexers:  indexers,
		committer: committer,
		querier:   querier,
	}
}

func (instance *IndexerBaseInstance) RunIndexerRound() {
	// create wait group for ALL indexers
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(instance.indexers))

	// baseCtx, cancel := context.WithTimeout(context.TODO(), 10*1000*time.Millisecond)
	baseCtx := context.TODO()

	for idx, indexer := range instance.indexers {
		runIndexer(
			&waitGroup,
			context.WithValue(baseCtx, "round", idx),
			instance.committer,
			instance.querier,
			indexer,
		)
	}

	waitGroup.Wait()
}

func runIndexer(
	wg *sync.WaitGroup,
	ctx context.Context,
	committer types.GraphQLCommitter,
	querier types.GraphQLQuerier,
	indexer types.Indexer,
) {
	var isolatedQuerier types.IndexerQuerier = createIsolatedQuerier(querier)
	var isolatedCommitter types.IndexerCommitter = createIsolatedCommitter(committer)

	go func() {
		defer wg.Done()
		indexer(isolatedQuerier, isolatedCommitter)
	}()
}

func createIsolatedQuerier(querier types.GraphQLQuerier) types.IndexerQuerier {
	return func(query interface{}, variables types.GraphQLParams) error {
		qs := generate.GenerateQuery(query)
		result := querier(qs, variables)

		if result.HasErrors() {
			for _, err := range result.Errors {
				fmt.Println(err)
			}

			return fmt.Errorf("Isolated Querier error")
		}

		return graph.UnmarshalGraphQLResult(result, &query)
	}
}

func createIsolatedCommitter(committer types.GraphQLCommitter) types.IndexerCommitter {
	return func(entity interface{}) {
		committer(entity)
	}
}
