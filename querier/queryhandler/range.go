package queryhandler

import (
	"bytes"
	"fmt"
	"github.com/terra-project/mantle/utils"
	"regexp"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
)

var RangeResolverParser = regexp.MustCompile(`@range\(\s*(?P<start>[\w\W]+)\s*,\s*(?P<end>[\w\W]+)\s*\)`)

type RangeResolver struct {
	db           db.DB
	kvindexEntry *kvindex.KVIndexEntry
	entityName   string
	indexName    string
	startKey     string
	endKey       string
	reverse      bool
}

func NewRangeResolver(
	db db.DB,
	kvindexEntry *kvindex.KVIndexEntry,
	entityName,
	indexName string,
	indexOption interface{},
) QueryHandler {
	indexOptionString, ok := indexOption.(string)
	if !ok {
		return nil
	}

	match := RangeResolverParser.FindStringSubmatch(indexOptionString)

	// not valid for this resolver, return nil
	if len(match) == 0 {
		return nil
	}

	var start string
	var end string
	for i := range RangeResolverParser.SubexpNames() {
		if i == 1 {
			start = match[1]
		} else if i == 2 {
			end = match[2]
		}
	}

	var reverse = false
	if bytes.Compare([]byte(start), []byte(end)) > 1 {
		reverse = true
	}

	return RangeResolver{
		db:           db,
		kvindexEntry: kvindexEntry,
		entityName:   entityName,
		indexName:    indexName,
		startKey:     start,
		endKey:       end,
		reverse:      reverse,
	}
}

func (resolver RangeResolver) Resolve() (QueryHandlerIterator, error) {
	kviEntry := resolver.kvindexEntry
	startKey, err := kviEntry.ResolveKeyType(resolver.startKey)
	if err != nil {
		return nil, fmt.Errorf(
			"Range parameter `start` cannot be converted to underlying index type, entityName=%s, indexName=%s, start=%s. %s",
			resolver.entityName,
			resolver.indexName,
			startKey,
			err,
		)
	}

	endKey, err := kviEntry.ResolveKeyType(resolver.endKey)
	if err != nil {
		return nil, fmt.Errorf(
			"Range parameter `end` cannot be converted to underlying index type, entityName=%s, indexName=%s, start=%s",
			resolver.entityName,
			resolver.indexName,
			endKey,
		)
	}

	db := resolver.db
	indexName := kviEntry.GetEntry().Name
	entityName := kviEntry.GetEntityName()
	it := db.IndexIterator(
		append([]byte(entityName), append([]byte(indexName), startKey...)...),
		resolver.reverse,
	)

	return NewRangeResolverIterator(entityName, indexName, endKey, it), nil
}

type RangeResolverIterator struct {
	entityName string
	indexName	string
	endKey []byte
	it db.Iterator
}

func NewRangeResolverIterator(
	entityName string,
	indexName string,
	endKey []byte,
	it db.Iterator,
) *RangeResolverIterator {
	return &RangeResolverIterator{
		entityName: entityName,
		indexName: indexName,
		endKey: endKey,
		it: it,
	}
}

func (resolver *RangeResolverIterator) Valid() bool {
	prefix := utils.ConcatBytes([]byte(resolver.entityName), []byte(resolver.indexName))

	isPrefixValid := resolver.it.Valid(prefix)

	// iteration is valid until
	// the compare{slice(len(item.Key)), endKey} is equal or lower
	isKeyValid := bytes.Compare(
		resolver.it.Key()[:len(resolver.endKey)],
		resolver.endKey,
	) <= 0

	return isPrefixValid && isKeyValid
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
