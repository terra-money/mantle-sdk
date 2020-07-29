package committer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/badger"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
	"github.com/vmihailenco/msgpack/v5"
)

func TestCommitter(t *testing.T) {
	var height = uint64(20288336)
	committer, db := createTestCommitter()

	err := committer.Commit(
		height,
		TestStruct{
			Foo: "foo",
			Bar: 1337,
		},
	)

	assert.Nil(t, err)

	// search by documentKey
	key := append([]byte("TestStruct"), utils.LeToBe(height)...)
	testStructBytes, err := db.Get(key)
	assert.Nil(t, err)
	testStructDocument := TestStruct{}
	msgpack.Unmarshal(testStructBytes, &testStructDocument)
	assert.Equal(t, testStructDocument.Foo, "foo")
	assert.Equal(t, testStructDocument.Bar, uint64(1337))

	// search by index (foo)
	// note that we are not using querier here
	// we're just checking if index exists and the document key matches
	// with the original document key
	func() {
		indexKey := []byte("TestStructfoofoo")
		it := db.Iterator(indexKey, indexKey, false)
		if it.Valid() {
			assert.Equal(t, utils.LeToBe(height), it.DocumentKey())
		} else {
			assert.FailNow(t, "index was never found")
		}
		it.Close()
	}()

	// search by index (bar)
	// same as before
	func() {
		indexKey := append([]byte("TestStructBar"), utils.LeToBe(uint64(1337))...)
		it := db.Iterator(indexKey, indexKey, false)
		if it.Valid() {
			fmt.Println(it.DocumentKey())
			assert.Equal(t, utils.LeToBe(height), it.DocumentKey())
		} else {
			assert.FailNow(t, "index was never found")
		}
		it.Close()
	}()
}

// TODO: do a failing test

type TestStruct struct {
	Foo string `mantle:"index=foo"`
	Bar uint64 `mantle:"index"`
}

func createTestCommitter() (Committer, db.DB) {
	db := badger.NewBadgerDB("")
	kvIndexes := kvindex.NewKVIndexMap(
		kvindex.NewKVIndex(reflect.TypeOf((*TestStruct)(nil))),
	)

	return NewCommitter(db, kvIndexes), db
}
