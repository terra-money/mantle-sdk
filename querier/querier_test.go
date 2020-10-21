package querier

import (
	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle-sdk/querier/queryhandler"
	"reflect"
	"testing"

	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/db/badger"
	"github.com/terra-project/mantle-sdk/db/kvindex"
)

type TestStruct struct {
	Foo struct {
		Bar uint64 `mantle:"index"`
	}
}

// TestNewQuerier we only sim upto whether queryhandlers are created.
// detailed operation of queryhandlers are tested per handler (in queryhandler directory)
// EVERYTHING in handlers must be tested
var handlers = handlersList

func TestNewQuerier(t *testing.T) {
	db := createTestDB()
	kvi := createTestKVIndex((*TestStruct)(nil))
	kvIndexMap := kvindex.NewKVIndexMap(kvi)

	querier := NewQuerier(db, kvIndexMap)

	// fails if entityName is not found in kvIndexMap
	func() {
		handler, err := querier.Build("NotFound", "index", "query")
		assert.Nil(t, handler)
		assert.NotNil(t, err)
	}()

	// fails if no matching handler is found
	// this is very hard to sim because the rule for "no matching query handler" is a little arbitrary.
	// seek resolver might pick this up
	func() {
		handler, err := querier.Build("TestStruct", "NonExistentIndex", "whateverQuery")
		assert.Nil(t, handler)
		assert.NotNil(t, err)
	}()

	// Range resolver (detailed sim in queryhandler/range)
	// - has _range suffix
	// - indexOption is a slice and the length is exactly 2
	// - otherwise error
	func() {
		handler, err := querier.Build("TestStruct", "Bar_range", []interface{}{1, 100})
		assert.Nil(t, err)
		assert.NotNil(t, handler)
		assert.Implements(t, (*queryhandler.QueryHandler)(nil), handler)
	}()

	// Height resolver
	// - indexName is Height
	// - indexOption __can__ be converted to uint64
	// - otherwise error
	func() {
		handler, err := querier.Build("TestStruct", "Height", 64)
		assert.Nil(t, err)
		assert.NotNil(t, handler)
		assert.Implements(t, (*queryhandler.QueryHandler)(nil), handler)
	}()

	// Seek resolver
	// - indexName is given
	// - kvIndexEntry.ResolveKeyType succeeds
	// - otherwise error
	func() {
		handler, err := querier.Build("TestStruct", "Bar", 64)
		assert.Nil(t, err)
		assert.NotNil(t, handler)
		assert.Implements(t, (*queryhandler.QueryHandler)(nil), handler)
	}()
}

func createTestDB() db.DB {
	return badger.NewBadgerDB("") // memdb
}

func createTestKVIndex(data interface{}) *kvindex.KVIndex {
	return kvindex.NewKVIndex(reflect.TypeOf(data))
}
