package badger

import (
	"bytes"
	"errors"
	"fmt"
	tmdb "github.com/tendermint/tm-db"

	badger "github.com/dgraph-io/badger/v2"
)

type BadgerCosmosAdapter struct {
	db *badger.DB
}

var (
	// errBatchClosed is returned when a closed or written batch is used.
	errBatchClosed = errors.New("batch has been written or closed")

	// errKeyEmpty is returned when attempting to use an empty or nil key.
	errKeyEmpty = errors.New("key cannot be empty")

	// errValueNil is returned when attempting to set a nil value.
	errValueNil = errors.New("value cannot be nil")
)

// NewBadgerDB creates a Badger key-value store backed to the
// directory dir supplied. If dir does not exist, it will be created.
func NewBadgerCosmosAdapter(db *badger.DB) (*BadgerCosmosAdapter) {
	return &BadgerCosmosAdapter{
		db: db,
	}
}

func (b *BadgerCosmosAdapter) Get(key []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errKeyEmpty
	}
	var val []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err == badger.ErrKeyNotFound {
			return nil
		} else if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		if err == nil && val == nil {
			val = []byte{}
		}
		return err
	})
	return val, err
}

func (b *BadgerCosmosAdapter) Has(key []byte) (bool, error) {
	if len(key) == 0 {
		return false, errKeyEmpty
	}
	var found bool
	err := b.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		found = (err != badger.ErrKeyNotFound)
		return nil
	})
	return found, err
}

func (b *BadgerCosmosAdapter) Set(key, value []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	if value == nil {
		return errValueNil
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func withSync(db *badger.DB, err error) error {
	if err != nil {
		return err
	}
	return db.Sync()
}

func (b *BadgerCosmosAdapter) SetSync(key, value []byte) error {
	return withSync(b.db, b.Set(key, value))
}

func (b *BadgerCosmosAdapter) Delete(key []byte) error {
	if len(key) == 0 {
		return errKeyEmpty
	}
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerCosmosAdapter) DeleteSync(key []byte) error {
	return withSync(b.db, b.Delete(key))
}

func (b *BadgerCosmosAdapter) Close() error {
	return b.db.Close()
}

func (b *BadgerCosmosAdapter) Print() error {
	return nil
}

func (b *BadgerCosmosAdapter) iteratorOpts(start, end []byte, opts badger.IteratorOptions) (*badgerDBIterator, error) {
	if (start != nil && len(start) == 0) || (end != nil && len(end) == 0) {
		return nil, errKeyEmpty
	}
	txn := b.db.NewTransaction(false)
	iter := txn.NewIterator(opts)
	iter.Rewind()
	iter.Seek(start)
	if opts.Reverse && iter.Valid() && bytes.Equal(iter.Item().Key(), start) {
		// If we're going in reverse, our starting point was "end",
		// which is exclusive.
		iter.Next()
	}
	return &badgerDBIterator{
		reverse: opts.Reverse,
		start:   start,
		end:     end,

		txn:  txn,
		iter: iter,
	}, nil
}

func (b *BadgerCosmosAdapter) Iterator(start, end []byte) (tmdb.Iterator, error) {
	opts := badger.DefaultIteratorOptions
	return b.iteratorOpts(start, end, opts)
}

func (b *BadgerCosmosAdapter) ReverseIterator(start, end []byte) (tmdb.Iterator, error) {
	opts := badger.DefaultIteratorOptions
	opts.Reverse = true
	return b.iteratorOpts(end, start, opts)
}

func (b *BadgerCosmosAdapter) Stats() map[string]string {
	return nil
}

func (b *BadgerCosmosAdapter) NewBatch() tmdb.Batch {
	wb := &badgerDBBatch{
		db:         b.db,
		wb:         b.db.NewWriteBatch(),
		firstFlush: make(chan struct{}, 1),
	}
	wb.firstFlush <- struct{}{}
	return wb
}

type badgerDBBatch struct {
	db *badger.DB
	wb *badger.WriteBatch

	// Calling db.Flush twice panics, so we must keep track of whether we've
	// flushed already on our own. If Write can receive from the firstFlush
	// channel, then it's the first and only Flush call we should do.
	//
	// Upstream bug report:
	// https://github.com/dgraph-io/badger/issues/1394
	firstFlush chan struct{}
}

func (b *badgerDBBatch) Set(key, value []byte) {
	if len(key) == 0 {
		//return errKeyEmpty
	}
	if value == nil {
		//return errValueNil
	}
	b.wb.Set(key, value)
}

func (b *badgerDBBatch) Delete(key []byte) {
	if len(key) == 0 {
		// noop
	}
	b.wb.Delete(key)
}

func (b *badgerDBBatch) Write() error {
	select {
	case <-b.firstFlush:
		return b.wb.Flush()
	default:
		return fmt.Errorf("batch already flushed")
	}
}

func (b *badgerDBBatch) WriteSync() error {
	return withSync(b.db, b.Write())
}

func (b *badgerDBBatch) Close() {
	select {
	case <-b.firstFlush: // a Flush after Cancel panics too
	default:
	}
	b.wb.Cancel()
}

type badgerDBIterator struct {
	reverse    bool
	start, end []byte

	txn  *badger.Txn
	iter *badger.Iterator

	lastErr error
}

func (i *badgerDBIterator) Close() {
	i.iter.Close()
	i.txn.Discard()
}

func (i *badgerDBIterator) Domain() (start, end []byte) { return i.start, i.end }
func (i *badgerDBIterator) Error() error                { return i.lastErr }

func (i *badgerDBIterator) Next() {
	if !i.Valid() {
		panic("iterator is invalid")
	}
	i.iter.Next()
}

func (i *badgerDBIterator) Valid() bool {
	if !i.iter.Valid() {
		return false
	}
	if len(i.end) > 0 {
		key := i.iter.Item().Key()
		if c := bytes.Compare(key, i.end); (!i.reverse && c >= 0) || (i.reverse && c < 0) {
			// We're at the end key, or past the end.
			return false
		}
	}
	return true
}

func (i *badgerDBIterator) Key() []byte {
	if !i.Valid() {
		panic("iterator is invalid")
	}
	// Note that we don't use KeyCopy, so this is only valid until the next
	// call to Next.
	return i.iter.Item().KeyCopy(nil)
}

func (i *badgerDBIterator) Value() []byte {
	if !i.Valid() {
		panic("iterator is invalid")
	}
	val, err := i.iter.Item().ValueCopy(nil)
	if err != nil {
		i.lastErr = err
	}
	return val
}