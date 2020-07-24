package queryhandler

import (
	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
)

type QueryHandlerBuilder func(
	db db.DB,
	kvindexEntry *kvindex.KVIndexEntry,
	entityName,
	indexName string,
	indexOption interface{},
) QueryHandler

type QueryHandler interface {
	Resolve() (interface{}, error)
}
