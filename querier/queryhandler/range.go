package queryhandler

import (
	"bytes"
	"fmt"
	"github.com/terra-money/mantle-sdk/utils"
	"strings"

	"github.com/terra-money/mantle-sdk/db"
	"github.com/terra-money/mantle-sdk/db/kvindex"
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
		reverse:      true,
	}, nil
}

func (resolver RangeResolver) Resolve() (QueryHandlerIterator, error) {
	kviEntry := resolver.kvindexEntry
	rangeStart, err := utils.ConvertToIndexValueToCorrectType(kviEntry.Type(), resolver.startKey)
	if err != nil {
		return nil, fmt.Errorf(
			"range parameter `start` cannot be converted to underlying index type, entityName=%s, indexName=%s, start=%s. %s",
			resolver.entityName,
			resolver.indexName,
			rangeStart,
			err,
		)
	}

	rangeEnd, err := utils.ConvertToIndexValueToCorrectType(kviEntry.Type(), resolver.endKey)
	if err != nil {
		return nil, fmt.Errorf(
			"range parameter `end` cannot be converted to underlying index type, entityName=%s, indexName=%s, start=%s",
			resolver.entityName,
			resolver.indexName,
			rangeEnd,
		)
	}

	db := resolver.db
	indexName := kviEntry.Name()
	entityName := resolver.entityName

	var seekKey, endKey []byte
	if resolver.reverse {
		seekKey = rangeEnd
		endKey = rangeStart
	} else {
		seekKey = rangeStart
		endKey = rangeEnd
	}

	it := db.IndexIterator(
		utils.BuildIndexIteratorPrefix(
			[]byte(entityName),
			[]byte(indexName),
			seekKey,
		),
		resolver.reverse,
	)

	return NewRangeResolverIterator(
		entityName,
		indexName,
		endKey,
		it,
		resolver.reverse,
	), nil
}

type RangeResolverIterator struct {
	entityName string
	indexName  string
	endKey     []byte
	it         db.Iterator
	prefix     []byte
	reverse    bool
}

func NewRangeResolverIterator(
	entityName string,
	indexName string,
	endKey []byte,
	it db.Iterator,
	reverse bool,
) *RangeResolverIterator {
	return &RangeResolverIterator{
		entityName: entityName,
		indexName:  indexName,
		endKey:     endKey,
		it:         it,
		prefix: utils.BuildIndexGroupPrefix(
			[]byte(entityName),
			[]byte(indexName),
		),
		reverse: reverse,
	}
}

func (iterator *RangeResolverIterator) Valid() bool {
	prefixValid := iterator.it.Valid(iterator.prefix)
	var withinRangeValid = false

	currentIndexKey := iterator.it.Key()[len(iterator.prefix):]
	currentIndexKey = currentIndexKey[:len(iterator.endKey)]

	comparison := bytes.Compare(
		currentIndexKey,
		iterator.endKey,
	)

	if iterator.reverse {
		// in case of reverse, currentKey should be the same or bigger than current
		withinRangeValid = comparison >= 0
	} else {
		withinRangeValid = comparison <= 0
	}

	return prefixValid && withinRangeValid
}

func (iterator *RangeResolverIterator) Next() {
	iterator.it.Next()
}

func (iterator *RangeResolverIterator) Key() []byte {
	return utils.BuildDocumentKey(
		[]byte(iterator.entityName),
		iterator.it.DocumentKey(),
	)
}

func (iterator *RangeResolverIterator) Close() {
	iterator.it.Close()
}
