package querier

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/badger"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
)

type TestStruct struct {
	Foo struct {
		Bar uint64 `mantle:"index"`
	}
}

func TestNewQuerier(t *testing.T) {
	db := createTestDB()
	kvi := createTestKVIndex((*TestStruct)(nil))

	// set some test data
	for i := 0; i < 100; i++ {
		// setting index by
		cursor, err := kvi.BuildIndexKey("Bar", uint64(i), utils.LeToBe(uint64(i)))
		assert.Equal(t, err, nil)

		// since we're just indexing, data can be null
		db.Set(cursor, nil)
	}

	querier := NewQuerier(db, kvindex.NewKVIndexMap(kvi))

	// range query
	// range query returns actual document ids
	func() {
		query := "@range(1, 2)"
		queryhandler, err := querier.Build("TestStruct", "Bar", query)
		assert.Equal(t, err, nil)

		ret, err := queryhandler.Resolve()
		assert.Equal(t, reflect.TypeOf(ret).Kind(), reflect.Slice)
		assert.Equal(t, len(ret.([][]byte)), 2)

		for i, docKey := range ret.([][]byte) {
			idx := i + 1
			assert.Equal(t, docKey, append([]byte("TestStruct"), utils.LeToBe(uint64(idx))...))
		}
	}()

	// seek query
	func() {
		query := 50
		queryhandler, err := querier.Build("TestStruct", "Bar", query)
		assert.Equal(t, err, nil)

		// ret is the document key
		ret, err := queryhandler.Resolve()
		assert.Equal(t, reflect.TypeOf(ret).Kind(), reflect.Slice, "buffer slice")
		assert.Equal(t, ret, append([]byte("TestStruct"), utils.LeToBe(uint64(50))...))

		// index 101 does not exist, should return error
		query = 101
		queryhandler, err = querier.Build("TestStruct", "Bar", query)
		assert.Equal(t, err, nil)

		ret, err = queryhandler.Resolve()
		assert.Equal(t, ret, nil)
		assert.NotNil(t, err)
	}()
}

func createTestDB() db.DB {
	return badger.NewBadgerDB("") // memdb
}

func createTestKVIndex(data interface{}) *kvindex.KVIndex {
	return kvindex.NewKVIndex(reflect.TypeOf(data))
}
