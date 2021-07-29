package queryhandler

import (
	"github.com/stretchr/testify/assert"
	"github.com/terra-money/mantle-sdk/committer"
	"github.com/terra-money/mantle-sdk/db/kvindex"
	"github.com/terra-money/mantle-sdk/db/leveldb"
	"github.com/terra-money/mantle-sdk/utils"
	"reflect"
	"testing"
)

func TestRangeResolver(t *testing.T) {
	type TestEntity struct {
		Foo string `mantle:"index"`
		Bar uint64 `mantle:"index"`
	}

	db := leveldb.NewLevelDB("../test_fixtures")
	kvIndex, _ := kvindex.NewKVIndex(reflect.TypeOf((*TestEntity)(nil)))
	entityName := reflect.TypeOf((*TestEntity)(nil)).Elem().Name()

	// range resolver is not created when
	// - indexName does NOT end with "_range"
	// - or indexOption is not a slice
	// - or indexOption slice length is NOT 2
	func() {
		var rr QueryHandler
		var err error
		rr, err = NewRangeResolver(db, kvIndex, entityName, "Foo", nil)
		assert.Nil(t, rr)
		assert.Nil(t, err)

		rr, err = NewRangeResolver(db, kvIndex, entityName, "Foo_range", "hello")
		assert.Nil(t, err)

		rr, err = NewRangeResolver(db, kvIndex, entityName, "Foo_range", []string{"1", "2", "3"})
		assert.Nil(t, err)
	}()

	// RangeResolver.Resolve fails if types don't match with the underlying type
	func() {
		// Foo is of string but indexOption has non-string value 2
		rr, err := NewRangeResolver(
			db,
			kvIndex,
			entityName,
			"Foo_range",
			[]interface{}{"1", 2},
		)
		assert.Nil(t, err)

		_, err = rr.Resolve()
		assert.NotNil(t, err)
	}()

	// RangeResolver success
	func() {
		committer := committer.NewCommitter(db, kvindex.NewKVIndexMap(kvIndex))
		// generate some testdata
		for i := 0; i < 26; i++ {
			testEntity := TestEntity{
				Foo: string(rune('a' + i)),
				Bar: uint64(i),
			}

			// should not fail
			if err := committer.Commit(uint64(i), testEntity); err != nil {
				panic(err)
			}
		}

		// testing Foo (string)
		func() {
			rr, err := NewRangeResolver(
				db,
				kvIndex,
				entityName,
				"Foo_range",
				[]interface{}{"c", "f"},
			)
			assert.NotNil(t, rr)
			assert.Nil(t, err)

			it, err := rr.Resolve()
			assert.Nil(t, err)

			i := 2 // from "c"
			for it.Valid() {
				assert.Equal(
					t,
					utils.ConcatBytes([]byte(entityName), utils.LeToBe(uint64(i))),
					it.Key(),
				)
				it.Next()
				i++
			}

			it.Close()
		}()

		// testing Bar (uint64)
		func() {
			rr, err := NewRangeResolver(
				db,
				kvIndex,
				entityName,
				"Bar_range",
				[]interface{}{5, 12},
			)
			assert.NotNil(t, rr)
			assert.Nil(t, err)

			it, err := rr.Resolve()
			assert.Nil(t, err)

			i := 5 // from "c"
			for it.Valid() {
				assert.Equal(
					t,
					utils.ConcatBytes([]byte(entityName), utils.LeToBe(uint64(i))),
					it.Key(),
				)
				it.Next()
				i++
			}

			it.Close()
		}()
	}()
}
