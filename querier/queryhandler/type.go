package queryhandler

import (
	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
)

type QueryHandlerBuilder func(
	db db.DB,
	kvIndex *kvindex.KVIndex,
	entityName,
	indexName string,
	indexOption interface{},
) (QueryHandler, error)

type QueryHandler interface {
	Resolve() (QueryHandlerIterator, error)
}

type QueryHandlerIterator interface {
	Valid() bool
	Next()
	Key() []byte
	Close()
}
