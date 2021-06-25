package querier

import (
	"github.com/terra-money/mantle-sdk/querier/queryhandler"
)

type Querier interface {
	Get([]byte) ([]byte, error)
	Build(entityName, indexName string, query interface{}) (queryhandler.QueryHandler, error)
}
