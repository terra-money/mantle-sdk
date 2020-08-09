package queryhandler

import (
	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
)

type HeightResolver struct {
	db           db.DB
	entityName   string
	indexName    string
	seekKey      []byte
}

// seek resolver
func NewHeightResolver(
	db db.DB,
	kvIndex *kvindex.KVIndex,
	entityName,
	indexName string,
	indexOption interface{},
) (QueryHandler, error) {
	if indexName != "Height" {
		return nil, nil
	}

	heightInUint64, _ := utils.GetUint64FromWhatever(indexOption)
	seekKey := utils.ConcatBytes([]byte(entityName), utils.LeToBe(heightInUint64))

	return &HeightResolver{
		db:           db,
		entityName:   entityName,
		indexName:    indexName,
		seekKey:      seekKey,
	}, nil
}

func (resolver HeightResolver) Resolve() (QueryHandlerIterator, error) {
	return NewHeightResolverIterator(resolver.seekKey), nil
}

// SeekResolverIterator never really iterates.
// Implemented this way because of interface acceptance.
// All methods (Valid, Next, Key, Close) will work to resolve documentKey
// only one time.
type HeightResolverIterator struct{
	key []byte
	isResolved bool
}

func NewHeightResolverIterator(key []byte) QueryHandlerIterator {
	return &HeightResolverIterator{
		key: key,
		isResolved: false,
	}
}

func (resolver *HeightResolverIterator) Valid() bool {
	return !resolver.isResolved
}
func (resolver *HeightResolverIterator) Next() {
	resolver.isResolved = true
}
func (resolver *HeightResolverIterator) Key() []byte {
	return resolver.key
}
func (resolver *HeightResolverIterator) Close() {
	// noop
}