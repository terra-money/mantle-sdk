package queryhandler

import (
	"fmt"
	"log"

	"bytes"
	"reflect"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
)

type SeekResolver struct {
	db           db.DB
	kvindexEntry *kvindex.KVIndexEntry
	entityName   string
	indexName    string
	seekKey      []byte
}

// seek resolver
func NewSeekResolver(
	db db.DB,
	kvindexEntry *kvindex.KVIndexEntry,
	entityName,
	indexName string,
	indexOption interface{},
) QueryHandler {
	seekKey, err := kvindexEntry.ResolveKeyType(indexOption)

	// if ResolveKeyType fails, that means
	if err != nil {
		log.Fatalf(
			"Hash index is given but the given option can't be used for index %s, entity=%s, indexOptionType=%s",
			indexName,
			entityName,
			reflect.TypeOf(indexOption).Kind().String(),
		)
		return nil
	}

	return &SeekResolver{
		db:           db,
		kvindexEntry: kvindexEntry,
		entityName:   entityName,
		indexName:    indexName,
		seekKey:      seekKey,
	}
}

func (resolver SeekResolver) Resolve() (QueryHandlerIterator, error) {
	var seekKeyPrefix = utils.ConcatBytes([]byte(resolver.entityName), []byte(resolver.indexName))
	var seekKeyActual = utils.ConcatBytes([]byte(resolver.entityName), []byte(resolver.indexName), resolver.seekKey)
	it := resolver.db.IndexIterator(
		seekKeyActual,
		false,
	)

	documentKey := new(bytes.Buffer)
	_, err := documentKey.Write([]byte(resolver.entityName))
	if err != nil {
		return nil, err
	}

	if it.Valid(seekKeyPrefix) {
		_, err := documentKey.Write([]byte(it.DocumentKey()))
		if err != nil {
			return nil, err
		}
		it.Close() // close immediately
	} else {
		return nil, fmt.Errorf(
			"Index does not exist, entityName=%s, indexName=%s, indexKey=%v",
			resolver.entityName,
			resolver.indexName,
			resolver.seekKey,
		)
	}

	return NewSeekResolverIterator(documentKey.Bytes()), nil
}

// SeekResolverIterator never really iterates.
// Implemented this way because of interface acceptance.
// All methods (Valid, Next, Key, Close) will work to resolve documentKey
// only one time.
type SeekResolverIterator struct {
	documentKey []byte
	isResolved  bool
}

func NewSeekResolverIterator(documentKey []byte) QueryHandlerIterator {
	return &SeekResolverIterator{
		documentKey: documentKey,
		isResolved:  false,
	}
}

func (resolver *SeekResolverIterator) Valid() bool {
	return !resolver.isResolved
}
func (resolver *SeekResolverIterator) Next() {
	resolver.isResolved = true
}
func (resolver *SeekResolverIterator) Key() []byte {
	return resolver.documentKey
}
func (resolver *SeekResolverIterator) Close() {
	// noop
}
