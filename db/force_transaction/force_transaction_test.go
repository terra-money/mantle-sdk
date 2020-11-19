package force_transaction

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/db/leveldb"
	"testing"
)

func TestForceTransaction(t *testing.T) {
	testdb := WithGlobalTransactionManager(leveldb.NewLevelDB("test"))

	// set transaction boundary
	testdb.SetGlobalTransactionBoundary()

	// roughly simulate tendermint/cosmos state save
	// simulation is for the worst case where
	// there is no CacheKV or any batching
	writeAction(testdb)

	// deliberate panic
	defer func() {
		// recover
		if r := recover(); r != nil {
			var i []byte
			var e error

			i, e = testdb.Get([]byte("foo"))
			if e != nil {
				fmt.Println(e)
			}

			assert.Zero(t, i)

			testdb.Close()
		}
	}()

	// deliberate panic
	panic("oops")
}

func writeAction(testdb db.DBwithGlobalTransaction) {
	testdb.Set([]byte("foo"), []byte{1})
	testdb.Set([]byte("bar"), []byte{1})

	// roughly simulate mantle committer behaviour
	batch := testdb.Batch()
	batch.Set([]byte("batchFoo"), []byte{1})
	batch.Set([]byte("batchBar"), []byte{1})

	batch.Delete([]byte("batchBar"))
	testdb.Delete([]byte("bar"))
}