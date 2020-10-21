package queryhandler

import (
	"fmt"
	"reflect"

	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/db/kvindex"
	"github.com/terra-project/mantle-sdk/utils"
)

type SeekResolver struct {
	db           db.DB
	kvindexEntry *kvindex.IndexEntry
	entityName   string
	indexName    string
	seekKey      []byte
}

// seek resolver
func NewSeekResolver(
	db db.DB,
	kvIndex *kvindex.KVIndex,
	entityName,
	indexName string,
	indexOption interface{},
) (QueryHandler, error) {
	// not seek op if indexName is not given
	if indexName == "" {
		return nil, nil
	}

	kvIndexEntry, kvIndexEntryExists := kvIndex.Entry(indexName)
	if !kvIndexEntryExists {
		return nil, fmt.Errorf("acquiring kvIndexEntry failed, queryHandler=seek, entityName=%s, indexName=%s", entityName, indexName)
	}

	seekKey, seekKeyErr := utils.ConvertToIndexValueToCorrectType(kvIndexEntry.Type(), indexOption)

	// if ResolveKeyType fails, that means
	if seekKeyErr != nil {
		return nil, fmt.Errorf("Hash index is given but the given option can't be used for index %s, entity=%s, indexOptionType=%s",
			indexName,
			entityName,
			reflect.TypeOf(indexOption).Kind().String())
	}

	return &SeekResolver{
		db:           db,
		kvindexEntry: &kvIndexEntry,
		entityName:   entityName,
		indexName:    indexName,
		seekKey:      seekKey,
	}, nil
}

func (resolver SeekResolver) Resolve() (QueryHandlerIterator, error) {
	entityNameInBytes := []byte(resolver.entityName)
	var seekKey = utils.BuildIndexIteratorPrefix(
		entityNameInBytes,
		[]byte(resolver.indexName),
		resolver.seekKey,
	)

	return NewSeekResolverIterator(
		entityNameInBytes,
		seekKey,
		resolver.db.IndexIterator(
			seekKey,
			true,
		),
	), nil
}

type SeekResolverIterator struct {
	entityName []byte
	prefix     []byte
	it         db.Iterator
	isResolved bool
}

func NewSeekResolverIterator(entityName, prefix []byte, it db.Iterator) QueryHandlerIterator {
	return &SeekResolverIterator{
		entityName: entityName,
		prefix:     prefix,
		it:         it,
	}
}

func (resolver *SeekResolverIterator) Valid() bool {
	return resolver.it.Valid(resolver.prefix)
}
func (resolver *SeekResolverIterator) Next() {
	resolver.it.Next()
}
func (resolver *SeekResolverIterator) Key() []byte {
	return utils.BuildDocumentKey(
		resolver.entityName,
		resolver.it.DocumentKey(),
	)
}
func (resolver *SeekResolverIterator) Close() {
	resolver.it.Close()
}
