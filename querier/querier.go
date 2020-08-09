package querier

import (
	"fmt"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/querier/queryhandler"
)

type QuerierInstance struct {
	db         db.DB
	kvindexMap *kvindex.KVIndexMap
}

// NewQuerier creates new query builder depending on the input
func NewQuerier(db db.DB, kvindexMap *kvindex.KVIndexMap) Querier {
	return &QuerierInstance{
		db:         db,
		kvindexMap: kvindexMap,
	}
}

// query pattern matcher
// note that precedence matters
var handlersList = []queryhandler.QueryHandlerBuilder{
	queryhandler.NewRangeResolver, // [1,2]
	queryhandler.NewHeightResolver, // Height: 2222
	queryhandler.NewSeekResolver, // someIndex: 2222
	// queryhandler.NewAggregationResolver,
}

// direct db getter
func (qi *QuerierInstance) Get(absoluteDocumentKey []byte) ([]byte, error) {
	return qi.db.Get(absoluteDocumentKey)
}

// Build returns a QueryHandler depending on entityName, indexName, query.
// if no appropriate QueryHandler is found, Build() returns an error.
func (qi *QuerierInstance) Build(entityName, indexName string, query interface{}) (queryhandler.QueryHandler, error) {
	kvIndex, ok := (*qi.kvindexMap)[entityName]

	if !ok {
		return nil, fmt.Errorf(
			"entity not found during queryhandler build. entityName=%s",
			entityName,
		)
	}

	for _, handlerBuilder := range handlersList {
		handler, err := handlerBuilder(qi.db, kvIndex, entityName, indexName, query)
		if err != nil {
			return nil, err
		}

		if handler == nil {
			continue
		} else {
			return handler, nil
		}
	}

	return nil, fmt.Errorf(
		"No matching handler found, entityName=%s, indexName=%s, query=%s",
		entityName,
		indexName,
		query,
	)
}
