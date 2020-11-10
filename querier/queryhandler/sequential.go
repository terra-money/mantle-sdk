package queryhandler

import (
	"github.com/terra-project/mantle-sdk/constants"
	"github.com/terra-project/mantle-sdk/utils"
	"math"

	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/db/kvindex"
)

type SequentialResolver struct {
	db             db.DB
	isPrimaryModel bool
	entityName     string
	reverse        bool
}

// NewSequentialResolver creates sequential resolver.
// SequentialResolver is used when NO index is given.
// SequentialResolver is only available to list queries
func NewSequentialResolver(
	db db.DB,
	kvIndex *kvindex.KVIndex,
	entityName,
	_ string,
	order interface{},
) (QueryHandler, error) {
	// iterate over document keys
	// note that because of this nature, list intersection calculation neeeds attention
	if !kvIndex.IsSliceModel() {
		return nil, nil
	}

	reverse := false
	if order.(int) == constants.DESC {
		reverse = true
	}

	return SequentialResolver{
		db:             db,
		entityName:     entityName,
		isPrimaryModel: kvIndex.IsPrimaryKeyedModel(),
		reverse:        reverse,
	}, nil
}

func (resolver SequentialResolver) Resolve() (QueryHandlerIterator, error) {
	entityNameInBytes := []byte(resolver.entityName)
	var seekKey []byte

	// if this is of primary keyed model,
	// we need to seek to the last key in the primary key index.
	// easiest way to achieve is this iterate from
	// a non-existent key that will always be larger than the
	// entire primary key group.
	// in our case, we can +1 to the last bit of document key group
	if resolver.reverse {
		if resolver.isPrimaryModel {
			seekKey = utils.BuildDocumentGroupPrefix(entityNameInBytes)
			seekKey[len(seekKey)-1] = seekKey[len(seekKey)-1] + 1
		} else {
			max, maxErr := utils.ConvertToLexicographicBytes(uint64(math.MaxUint64))
			if maxErr != nil {
				return nil, maxErr
			}
			seekKey = utils.BuildDocumentKey(entityNameInBytes, max)
		}
	} else {
		seekKey = utils.BuildDocumentGroupPrefix(entityNameInBytes)
	}

	db := resolver.db

	// sequential is always reverse.
	// TODO: fix me if there is a request
	it := db.IndexIterator(
		seekKey,
		resolver.reverse,
	)

	return NewSequentialResolverIterator(entityNameInBytes, seekKey, it), nil
}

type SequentialResolverIterator struct {
	entityName []byte
	seekKey    []byte
	it         db.Iterator
	prefix     []byte
}

func NewSequentialResolverIterator(
	entityName []byte,
	seekKey []byte,
	it db.Iterator,
) *SequentialResolverIterator {
	return &SequentialResolverIterator{
		entityName: entityName,
		seekKey:    seekKey,
		it:         it,
		prefix:     utils.BuildDocumentGroupPrefix(entityName),
	}
}

// Sequential validity is guaranteed as long as prefix is solid
func (resolver *SequentialResolverIterator) Valid() bool {
	return resolver.it.Valid(resolver.prefix)
}

func (resolver *SequentialResolverIterator) Next() {
	resolver.it.Next()
}

func (resolver *SequentialResolverIterator) Key() []byte {
	return resolver.it.Key()
}

func (resolver *SequentialResolverIterator) Close() {
	resolver.it.Close()
}
