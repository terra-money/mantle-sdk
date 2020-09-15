package committer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle/utils"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/terra-project/mantle/db/badger"
	"github.com/terra-project/mantle/db/kvindex"
)

func TestCommitter(t *testing.T) {
	// sim simple struct
	func() {
		type TestStruct struct {
			Foo string
			Bar string `mantle:"index"`
		}

		testdb := badger.NewBadgerDB("")
		kvIndexMap := kvindex.NewKVIndexMap(
			kvindex.NewKVIndex(reflect.TypeOf((*TestStruct)(nil))),
		)
		committer := NewCommitter(testdb, kvIndexMap)

		// commit generates the following keys (key handles are converted to big endians):
		// TestStruct#1
		// TestStruct@Bar:ToBeIndexed#1
		// TestStruct@h:100#1
		// (pk is always 0 in this case)
		entity := TestStruct{
			Foo: "foo",
			Bar: "ToBeIndexed",
		}
		err := committer.Commit(uint64(100), entity)

		assert.Nil(t, err)

		// primary document exists
		val, valErr := testdb.Get(utils.BuildDocumentKey(
			[]byte("TestStruct"),
			utils.LeToBe(0),
		))
		assert.Nil(t, valErr)

		queryEntity := TestStruct{}
		unpackErr := msgpack.Unmarshal(val, &queryEntity)
		if unpackErr != nil {
			t.Fail()
		}
		assert.Equal(t, entity, queryEntity)

		// height index exists
		val, valErr = testdb.Get(utils.BuildIndexedDocumentKey(
			[]byte("TestStruct"),
			[]byte("Height"),
			utils.LeToBe(100),
			utils.LeToBe(0),
		))

		// testing the existence of key is enough,
		// because key itself has the seq pointer. we can always point back to
		// the original document.
		assert.Nil(t, valErr)

		// arbitrary index exists
		val, valErr = testdb.Get(utils.BuildIndexedDocumentKey(
			[]byte("TestStruct"),
			[]byte("Bar"),
			[]byte("ToBeIndexed"),
			utils.LeToBe(0),
		))

		assert.Nil(t, valErr)
	}()

	// sim slice struct
	func() {
		type TestSliceStruct []struct {
			Foo string
			Bar string `mantle:"index"`
		}

		testdb := badger.NewBadgerDB("")
		kvIndexMap := kvindex.NewKVIndexMap(
			kvindex.NewKVIndex(reflect.TypeOf((*TestSliceStruct)(nil))),
		)
		committer := NewCommitter(testdb, kvIndexMap)

		// commit generates the following keys (key handles are converted to big endians):
		// TestStruct#1
		// TestStruct@Bar:ToBeIndexed#1
		// TestStruct@h:100#1
		// (pk is always 0 in this case)
		entity := TestSliceStruct{
			{
				Foo: "foo",
				Bar: "Bar1",
			},
			// sim overlap as well
			{
				Foo: "foo",
				Bar: "Bar1",
			},
			{
				Foo: "foo",
				Bar: "Bar2",
			},
		}
		err := committer.Commit(uint64(100), entity)

		assert.Nil(t, err)

		// primary documents exist
		for i := 0; i < len(entity); i++ {
			_, valErr := testdb.Get(utils.BuildDocumentKey(
				[]byte("TestSliceStruct"),
				utils.LeToBe(uint64(i)),
			))
			assert.Nil(t, valErr)
		}

		// indexed documents exist
		// bar1
		prefix := utils.BuildIndexIteratorPrefix(
			[]byte("TestSliceStruct"),
			[]byte("Bar"),
			[]byte("Bar1"),
		)

		it := testdb.Iterator(prefix, false)

		keys := make([][]byte, 0)
		for it.Valid(prefix) {
			docKey := it.DocumentKey()
			keys = append(keys, docKey)
			it.Next()
		}

		it.Close()

		// bar2
		prefix = utils.BuildIndexIteratorPrefix(
			[]byte("TestSliceStruct"),
			[]byte("Bar"),
			[]byte("Bar2"),
		)

		for it.Valid(prefix) {
			docKey := it.DocumentKey()

			keys = append(keys, docKey)
			it.Next()
		}
		it.Close()

		for i, key := range keys {
			assert.Equal(t, utils.LeToBe(uint64(i)), key)
		}

	}()

	// sim map struct
	func() {
		type Entity struct {
			Foo string
			Bar string `mantle:"index"`
		}
		type TestMapStruct map[string]Entity

		testdb := badger.NewBadgerDB("")
		kvIndexMap := kvindex.NewKVIndexMap(
			kvindex.NewKVIndex(reflect.TypeOf((TestMapStruct)(nil))),
		)
		committer := NewCommitter(testdb, kvIndexMap)
		docKeys := []string{"doc1", "doc2", "doc3"}

		// commit generates the following keys (key handles are converted to big endians):
		// TestStruct#1
		// TestStruct@Bar:ToBeIndexed#1
		// TestStruct@h:100#1
		// (pk is always 0 in this case)
		entity := TestMapStruct{
			docKeys[0]: {
				Foo: "foo",
				Bar: "Bar1",
			},
			// sim overlap as well
			docKeys[1]: {
				Foo: "foo",
				Bar: "Bar1",
			},
			docKeys[2]: {
				Foo: "foo",
				Bar: "Bar2",
			},
		}
		err := committer.Commit(uint64(100), entity)

		assert.Nil(t, err)

		// primary documents exist
		for i := 0; i < len(entity); i++ {
			_, valErr := testdb.Get(utils.BuildDocumentKey(
				[]byte("TestMapStruct"),
				[]byte(docKeys[i]),
			))
			assert.Nil(t, valErr)
		}

		// indexed documents exist
		// bar1
		prefix := utils.BuildIndexIteratorPrefix(
			[]byte("TestSliceStruct"),
			[]byte("Bar"),
			[]byte("Bar1"),
		)

		it := testdb.Iterator(prefix, false)

		keys := make([][]byte, 0)
		for it.Valid(prefix) {
			docKey := it.DocumentKey()
			keys = append(keys, docKey)
			it.Next()
		}

		it.Close()

		// bar2
		prefix = utils.BuildIndexIteratorPrefix(
			[]byte("TestSliceStruct"),
			[]byte("Bar"),
			[]byte("Bar2"),
		)

		for it.Valid(prefix) {
			docKey := it.DocumentKey()
			keys = append(keys, docKey)
			it.Next()
		}
		it.Close()

		for i, key := range keys {
			assert.Equal(t, utils.LeToBe(uint64(i)), key)
		}

	}()
}
