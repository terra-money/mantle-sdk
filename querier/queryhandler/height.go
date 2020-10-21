package queryhandler

import (
	"bytes"
	"fmt"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/db/kvindex"
	"github.com/terra-project/mantle-sdk/utils"
)

type HeightResolver struct {
	db          db.DB
	entityName  string
	indexName   string
	prefixStart []byte
	prefixEnd   []byte
}

// seek resolver
func NewHeightResolver(
	db db.DB,
	_ *kvindex.KVIndex,
	entityName,
	indexName string,
	indexOption interface{},
) (QueryHandler, error) {
	if !(indexName == "Height" || indexName == "Height_range") {
		return nil, nil
	}

	switch indexOption.(type) {
	case []interface{}:
		heightRange, _ := indexOption.([]interface{})
		return &HeightResolver{
			db:          db,
			entityName:  entityName,
			indexName:   "Height",
			prefixStart: utils.LeToBe(uint64(heightRange[0].(int))),
			prefixEnd:   utils.LeToBe(uint64(heightRange[1].(int))),
		}, nil
	case int:
		height, _ := indexOption.(int)
		return &HeightResolver{
			db:          db,
			entityName:  entityName,
			indexName:   "Height",
			prefixStart: utils.LeToBe(uint64(height)),
			prefixEnd:   nil,
		}, nil
	default:
		return nil, fmt.Errorf("invalid height parameters, entityName=%s, indexOption=%v", entityName, indexOption)
	}
}

func (resolver HeightResolver) Resolve() (QueryHandlerIterator, error) {
	entityNameInBytes := []byte(resolver.entityName)
	indexNameInBytes := []byte(resolver.indexName)

	prefixGroup := utils.BuildIndexGroupPrefix(
		entityNameInBytes,
		indexNameInBytes,
	)

	prefixStart := utils.BuildIndexIteratorPrefix(
		entityNameInBytes,
		indexNameInBytes,
		resolver.prefixStart,
	)

	var prefixEnd []byte = nil
	if resolver.prefixEnd != nil {
		prefixEnd = utils.BuildIndexIteratorPrefix(
			entityNameInBytes,
			indexNameInBytes,
			resolver.prefixEnd,
		)
	} else {
		prefixEnd = prefixStart
	}

	return NewHeightResolverIterator(
		entityNameInBytes,
		prefixGroup,
		prefixStart,
		prefixEnd,
		resolver.db.IndexIterator(
			prefixStart,
			false,
		),
	), nil
}

type HeightResolverIterator struct {
	entityName  []byte
	prefixGroup []byte
	prefixStart []byte
	prefixEnd   []byte
	it          db.Iterator
}

func NewHeightResolverIterator(entityName, prefixGroup, prefixStart, prefixEnd []byte, it db.Iterator) QueryHandlerIterator {
	return &HeightResolverIterator{
		entityName:  entityName,
		prefixGroup: prefixGroup,
		prefixStart: prefixStart,
		prefixEnd:   prefixEnd,
		it:          it,
	}
}

func (resolver *HeightResolverIterator) Valid() bool {
	if len(resolver.prefixEnd) > 0 {
		psCompare := bytes.Compare(resolver.it.Key(), resolver.prefixStart)
		peCompare := bytes.Compare(resolver.it.Key(), utils.GetReverseSeekKeyFromIndexGroupPrefix(resolver.prefixEnd))
		return psCompare == 1 && peCompare <= 0
	} else {
		return resolver.it.Valid(resolver.prefixStart)
	}

}
func (resolver *HeightResolverIterator) Next() {
	resolver.it.Next()
}
func (resolver *HeightResolverIterator) Key() []byte {
	return utils.BuildDocumentKey(
		resolver.entityName,
		resolver.it.DocumentKey(),
	)
}
func (resolver *HeightResolverIterator) Close() {
	resolver.it.Close()
}
