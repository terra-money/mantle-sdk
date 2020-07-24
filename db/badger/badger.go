package badger

import (
	"bytes"

	bd "github.com/dgraph-io/badger/v2"
	tmdb "github.com/tendermint/tm-db"
	"github.com/terra-project/mantle/db"
)

type BadgerDB struct {
	db *bd.DB
}

func NewBadgerDB(path string) *BadgerDB {
	var inMemory bool
	if path == "" {
		inMemory = true
	}
	// TODO: tweak me
	options := bd.
		DefaultOptions(path).
		WithInMemory(inMemory)

	db, err := bd.Open(options)
	if err != nil {
		panic(err)
	}

	return &BadgerDB{
		db: db,
	}
}

func (bdb *BadgerDB) GetCosmosAdapter() tmdb.DB {
	return NewBadgerCosmosAdapter(bdb)
}

func (bdb *BadgerDB) GetDB() *bd.DB {
	return bdb.db
}

func (bdb *BadgerDB) Get(key []byte) ([]byte, error) {
	txn := bdb.db.NewTransaction(false)
	defer txn.Discard()

	ret, err := txn.Get(key)
	if err != nil {
		return nil, err
	}

	val, err := ret.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (bdb *BadgerDB) Set(key, data []byte) error {
	txn := bdb.db.NewTransaction(true)
	defer txn.Commit()

	if err := txn.Set(key, data); err != nil {
		return err
	}

	return nil
}

func (bdb *BadgerDB) Delete(key []byte) error {
	txn := bdb.db.NewTransaction(true)
	defer txn.Commit()

	if err := txn.Delete(key); err != nil {
		return err
	}

	return nil
}

type BadgerIterator struct {
	it             *bd.Iterator
	txn            *bd.Txn
	start          []byte
	end            []byte
	indexKeyLength int
	reverse        bool
}

func (bdb *BadgerDB) Iterator(start, end []byte, reverse bool) db.Iterator {
	txn := bdb.db.NewTransaction(false)
	itOpts := bd.DefaultIteratorOptions
	itOpts.PrefetchValues = true
	itOpts.Reverse = reverse
	it := txn.NewIterator(itOpts)

	it.Seek(start)

	return &BadgerIterator{
		txn:            txn,
		it:             it,
		start:          start,
		end:            end,
		indexKeyLength: len(start),
		reverse:        reverse,
	}
}

func (bdb *BadgerDB) IndexIterator(start, end []byte, reverse bool) db.Iterator {
	txn := bdb.db.NewTransaction(false)
	itOpts := bd.DefaultIteratorOptions
	itOpts.PrefetchValues = false
	itOpts.Reverse = reverse
	it := txn.NewIterator(itOpts)

	it.Seek(start)

	return &BadgerIterator{
		txn:            txn,
		it:             it,
		start:          start,
		end:            end,
		indexKeyLength: len(start),
		reverse:        reverse,
	}
}

func (it *BadgerIterator) Close() {
	it.it.Close()
	it.txn.Discard()
}

func (it *BadgerIterator) Valid() bool {
	if !it.it.Valid() {
		return false
	}

	currentKey := it.Key()
	endLength := len(it.end)
	comp := bytes.Compare(it.end, currentKey[:endLength])

	return comp > -1
}

func (it *BadgerIterator) Next() {
	it.it.Next()
}

func (it *BadgerIterator) Key() []byte {
	return it.it.Item().Key()
}

func (it *BadgerIterator) Value() []byte {
	ret, err := it.it.Item().ValueCopy(nil)
	if err != nil {
		panic(err)
	}
	return ret
}

func (it *BadgerIterator) DocumentKey() []byte {
	k := it.Key()
	return k[it.indexKeyLength:]
}

type BadgerBatch struct {
	batch *bd.WriteBatch
}

func (bdb *BadgerDB) Batch() db.Batch {
	return &BadgerBatch{
		batch: bdb.db.NewWriteBatch(),
	}
}

func (batch *BadgerBatch) Set(key, data []byte) error {
	return batch.batch.Set(key, data)
}

func (batch *BadgerBatch) Delete(key []byte) error {
	return batch.batch.Delete(key)
}

func (batch *BadgerBatch) Flush() error {
	return batch.batch.Flush()
}

func (batch *BadgerBatch) Close() {
	batch.batch.Cancel()
}
