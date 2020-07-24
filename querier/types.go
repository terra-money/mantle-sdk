package querier

import (
	"github.com/terra-project/mantle/querier/queryhandler"
)

type Querier interface {
	Get([]byte) ([]byte, error)
	Build(entityName, indexName string, query interface{}) (queryhandler.QueryHandler, error)
}
