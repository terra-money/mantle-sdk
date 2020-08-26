package queryhandler

import (
	"fmt"
	"strings"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/utils"
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
	_,
	entityName,
	indexName string,
	indexOption interface{},
) (QueryHandler, error) {
	if !strings.HasPrefix(indexName, "Height") {
		return nil, nil
	}

	if heightRange, isHeightRange := indexOption.([]uint64); isHeightRange {
		return &HeightResolver{
			db:          db,
			entityName:  entityName,
			indexName:   indexName,
			prefixStart: utils.LeToBe(heightRange[0]),
			prefixEnd:   utils.LeToBe(heightRange[1]),
		}, nil
	}

	if heightSingle, isHeightSingle := indexOption.(uint64); isHeightSingle {
		return &HeightResolver{
			db:          db,
			entityName:  entityName,
			indexName:   indexName,
			prefixStart: utils.LeToBe(heightSingle),
			prefixEnd:   nil,
		}, nil
	}

	return nil, fmt.Errorf("invalid height parameters, entityName=%s, indexOption=%v", entityName, indexOption)
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
	}

	return NewHeightResolverIterator(
		entityNameInBytes,
		prefixGroup,
		prefixStart,
		prefixEnd,
		resolver.db.IndexIterator(
			prefixStart,
			true,
		),
	), nil
}

type HeightResolverIterator struct {
	entityName  []byte
	prefixGroup []byte
	prefixStart []byte
	prefixEnd   []byte
	it          db.Iterator
	isResolved  bool
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
	// isPrefixValid := resolver.it.Valid(resolver.prefixGroup)
	return resolver.it.Valid(resolver.prefixGroup)
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
