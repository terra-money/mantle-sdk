package queryhandler

import (
	"bytes"
	"fmt"
	"github.com/terra-project/mantle/utils"
	"strings"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
)

type RangeResolver struct {
	db           db.DB
	kvindexEntry *kvindex.IndexEntry
	entityName   string
	indexName    string
	startKey     interface{}
	endKey       interface{}
	reverse      bool
}

func NewRangeResolver(
	db db.DB,
	kvIndex *kvindex.KVIndex,
	entityName,
	indexName string,
	indexOption interface{},
) (QueryHandler, error) {
	indexOptionSlice, isIndexOptionSlice := indexOption.([]interface{})
	isIndexNameRange := strings.HasSuffix(indexName, "_range")

	if !(isIndexOptionSlice && isIndexNameRange && len(indexOptionSlice) == 2) {
		return nil, nil
	}

	kvIndexEntry, kvIndexEntryExists := kvIndex.Entry(indexName[:len(indexName)-6])
	if !kvIndexEntryExists {
		return nil, fmt.Errorf("acquiring kvIndexEntry failed, entityName=%s, indexName=%s", entityName, indexName)
	}

	return RangeResolver{
		db:           db,
		kvindexEntry: &kvIndexEntry,
		entityName:   entityName,
		indexName:    indexName,
		startKey:     indexOptionSlice[0],
		endKey:       indexOptionSlice[1],
	}, nil
}

func (resolver RangeResolver) Resolve() (QueryHandlerIterator, error) {
	kviEntry := resolver.kvindexEntry
	startKey, err := utils.ConvertToLexicographicBytes(resolver.startKey)
	if err != nil {
		return nil, fmt.Errorf(
			"range parameter `start` cannot be converted to underlying index type, entityName=%s, indexName=%s, start=%s. %s",
			resolver.entityName,
			resolver.indexName,
			startKey,
			err,
		)
	}

	endKey, err := utils.ConvertToLexicographicBytes(resolver.endKey)
	if err != nil {
		return nil, fmt.Errorf(
			"range parameter `end` cannot be converted to underlying index type, entityName=%s, indexName=%s, start=%s",
			resolver.entityName,
			resolver.indexName,
			endKey,
		)
	}

	db := resolver.db
	indexName := kviEntry.Name()
	entityName := resolver.entityName
	it := db.IndexIterator(
		utils.BuildIndexIteratorPrefix(
			[]byte(entityName),
			[]byte(indexName),
			startKey,
		),
		resolver.reverse,
	)

	return NewRangeResolverIterator(entityName, indexName, endKey, it), nil
}

type RangeResolverIterator struct {
	entityName string
	indexName  string
	endKey     []byte
	it         db.Iterator
	prefix     []byte
}

func NewRangeResolverIterator(
	entityName string,
	indexName string,
	endKey []byte,
	it db.Iterator,
) *RangeResolverIterator {
	return &RangeResolverIterator{
		entityName: entityName,
		indexName:  indexName,
		endKey:     endKey,
		it:         it,
		prefix:     utils.ConcatBytes([]byte(entityName), []byte(indexName)),
	}
}

func (resolver *RangeResolverIterator) Valid() bool {
	isPrefixValid := resolver.it.Valid(resolver.prefix)

	if !isPrefixValid {
		return false
	}

	// iteration is valid until
	// the compare{slice(len(item.Key)), endKey} is equal or lower
	key := resolver.it.Key()
	comp := key[:len(key)-8][len(resolver.prefix):]
	isKeyValid := bytes.Compare(
		comp,
		resolver.endKey,
	) <= 0

	return isKeyValid
}

func (resolver *RangeResolverIterator) Next() {
	resolver.it.Next()
}

func (resolver *RangeResolverIterator) Key() []byte {
	return append([]byte(resolver.entityName), resolver.it.DocumentKey()...)
}

func (resolver *RangeResolverIterator) Close() {
	resolver.it.Close()
}
