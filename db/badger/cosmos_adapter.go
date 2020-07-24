package badger

import (
	"bytes"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v2"
	tmdb "github.com/tendermint/tm-db"
)

type BadgerCosmosAdapter struct {
	db *BadgerDB
}

func NewBadgerCosmosAdapter(
	db *BadgerDB,
) *BadgerCosmosAdapter {
	return &BadgerCosmosAdapter{
		db,
	}
}

func (adapter *BadgerCosmosAdapter) Get(key []byte) ([]byte, error) {
	ret, err := adapter.db.Get(key)
	if err != badger.ErrKeyNotFound {
		return nil, nil
	}
	return ret, nil
}

func (adapter *BadgerCosmosAdapter) Has(key []byte) (bool, error) {
	entry, err := adapter.Get(key)

	if err != nil {
		return false, err
	}

	if entry != nil {
		return false, nil
	} else {
		return false, nil
	}
}

func (adapter *BadgerCosmosAdapter) Set(key []byte, value []byte) error {
	return adapter.db.Set(key, value)
}

func (adapter *BadgerCosmosAdapter) SetSync(key []byte, value []byte) error {
	return adapter.Set(key, value)
}

func (adapter *BadgerCosmosAdapter) Delete(key []byte) error {
	err := adapter.db.Delete(key)
	if err != nil {
		fmt.Printf("db/badger/cosmos_adapter.Delete: entry doesn't exist, key=%s", key)
	}
	return err
}

func (adapter *BadgerCosmosAdapter) DeleteSync(key []byte) error {
	return adapter.Delete(key)
}

func (adapter *BadgerCosmosAdapter) Close() error {
	return adapter.Close()
}

func (adapter *BadgerCosmosAdapter) Print() error {
	return nil
}
func (adapter *BadgerCosmosAdapter) Stats() map[string]string {
	return nil
}

// iterator interface that follows dbm.Iterator
type BadgerCosmosAdapterIterator struct {
	iterator  *badger.Iterator
	txn       *badger.Txn
	reverse   bool
	startKey  []byte
	endKey    []byte
	lastError error
}

func (adapter *BadgerCosmosAdapter) Iterator(start, end []byte) (tmdb.Iterator, error) {
	bdb := adapter.db.GetDB()

	// create readonly transaction
	txn := bdb.NewTransaction(false)
	itOpts := badger.DefaultIteratorOptions
	itOpts.PrefetchSize = 10 // TODO: optimize me
	it := txn.NewIterator(itOpts)

	return &BadgerCosmosAdapterIterator{
		iterator: it,
		txn:      txn,
		reverse:  false,
		startKey: start,
		endKey:   end,
	}, nil
}

func (adapter *BadgerCosmosAdapter) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	bdb := adapter.db.GetDB()

	// create readonly transaction
	txn := bdb.NewTransaction(false)
	itOpts := badger.DefaultIteratorOptions
	itOpts.PrefetchSize = 10 // TODO: optimize me
	it := txn.NewIterator(itOpts)

	return &BadgerCosmosAdapterIterator{
		iterator: it,
		txn:      txn,
		reverse:  true,
		startKey: start,
		endKey:   end,
	}, nil
}

func (it *BadgerCosmosAdapterIterator) Error() error {
	return it.lastError
}

func (it *BadgerCosmosAdapterIterator) Close() {
	it.iterator.Close()
}

func (it *BadgerCosmosAdapterIterator) Domain() (start, end []byte) {
	return it.startKey, it.endKey
}

func (it *BadgerCosmosAdapterIterator) Valid() bool {
	if !it.iterator.Valid() {
		return false
	}
	// if end key is set, check the current key and see if endKey is higher than the current key
	if len(it.endKey) > 0 {
		key := it.iterator.Item().Key()
		if c := bytes.Compare(key, it.endKey); (!it.reverse && c >= 0) || (it.reverse && c < 0) {
			// We're at the end key, or past the end.
			return false
		}
	}
	return true
}

func (it *BadgerCosmosAdapterIterator) Next() {
	it.iterator.Next()
}

func (it *BadgerCosmosAdapterIterator) Key() []byte {
	if !it.Valid() {
		panic("db/badger/cosmos_adapter.BadgerCosmosAdapterIterator.Value: Iteration is invalid")
	}

	return it.iterator.Item().KeyCopy(nil)
}

func (it *BadgerCosmosAdapterIterator) Value() []byte {
	if !it.Valid() {
		panic("db/badger/cosmos_adapter.BadgerCosmosAdapterIterator.Value: Iteration is invalid")
	}

	val, err := it.iterator.Item().ValueCopy(nil)
	if err != nil {
		it.lastError = err
	}

	return val
}

// batch interface that follows dbm.Iterator
type BadgerCosmosAdapterBatch struct {
	db      *badger.DB
	adapter *BadgerCosmosAdapter
	wb      *badger.WriteBatch

	// Calling db.Flush twice panics, so we must keep track of whether we've
	// flushed already on our own. If Write can receive from the firstFlush
	// channel, then it's the first and only Flush call we should do.
	//
	// Upstream bug report:
	// https://github.com/dgraph-io/badger/issues/1394
	firstFlush chan struct{}
}

func (adapter *BadgerCosmosAdapter) NewBatch() tmdb.Batch {
	db := adapter.db.GetDB()
	batch := &BadgerCosmosAdapterBatch{
		db:         db,
		adapter:    adapter,
		wb:         db.NewWriteBatch(),
		firstFlush: make(chan struct{}, 1),
	}
	batch.firstFlush <- struct{}{}
	return batch
}

func (b *BadgerCosmosAdapterBatch) Set(key, value []byte) {
	if len(key) == 0 {
		log.Fatal(fmt.Errorf("db/badger/cosmos_adapter.BadgerCosmosAdapterBatch.Set: Key is empty"))
	}
	if value == nil {
		log.Fatal(fmt.Errorf("db/badger/cosmos_adapter.BadgerCosmosAdapterBatch.Set: Value is nil"))
	}

	b.wb.Set(key, value)
}

func (b *BadgerCosmosAdapterBatch) Delete(key []byte) {
	if len(key) == 0 {
		fmt.Errorf("db/badger/cosmos_adapter.BadgerCosmosAdapterBatch.Set: Key is empty")
	}
	b.wb.Delete(key)
}

func (b *BadgerCosmosAdapterBatch) Write() error {
	select {
	case <-b.firstFlush:
		return b.wb.Flush()
	default:
		return fmt.Errorf("db/badger/cosmos_adapter.BadgerCosmosAdapterBatch.Set: Batch is already closed")
	}
}

func (b *BadgerCosmosAdapterBatch) WriteSync() error {
	return b.Write()
}

func (b *BadgerCosmosAdapterBatch) Close() {
	select {
	case <-b.firstFlush: // a Flush after Cancel panics too
	default:
	}
	b.wb.Cancel()
}
