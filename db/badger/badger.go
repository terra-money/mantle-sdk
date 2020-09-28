package badger

import (
	"fmt"
	bd "github.com/dgraph-io/badger/v2"
	tmdb "github.com/tendermint/tm-db"
	compatbadger "github.com/terra-project/mantle-compatibility/badger"
	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/utils"
)

type BadgerDB struct {
	db *bd.DB
}

var maxPKRange = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

func NewBadgerDB(path string) db.DB {
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

func (bdb *BadgerDB) Compact() error {
	bdb.db.Flatten(8)
	if err := bdb.db.RunValueLogGC(0.5); err == bd.ErrNoRewrite {
		fmt.Println("nothing to compact!")
	}

	return nil
}

func (bdb *BadgerDB) GetCosmosAdapter() tmdb.DB {
	return compatbadger.NewBadgerCosmosAdapter(bdb.db)
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

func (bdb *BadgerDB) GetSequence(key []byte, bandwidth uint64) (db.Sequence, error) {
	return bdb.db.GetSequence(key, bandwidth)
}

func (bdb *BadgerDB) Close() error {
	return bdb.db.Close()
}

type BadgerIterator struct {
	it             *bd.Iterator
	txn            *bd.Txn
	start          []byte
	indexKeyLength int
	reverse        bool
}

func (bdb *BadgerDB) Iterator(
	start []byte,
	reverse bool,
) db.Iterator {
	txn := bdb.db.NewTransaction(false)
	itOpts := bd.DefaultIteratorOptions
	itOpts.PrefetchValues = true
	itOpts.Reverse = reverse
	it := txn.NewIterator(itOpts)

	// if iterator goes backwards, so start needs to be the biggest of that index start range

	if reverse {
		it.Seek(utils.ConcatBytes(start, maxPKRange))
	} else {
		it.Seek(start)
	}

	return &BadgerIterator{
		txn:            txn,
		it:             it,
		start:          start,
		indexKeyLength: len(start),
		reverse:        reverse,
	}
}

func (bdb *BadgerDB) IndexIterator(
	start []byte,
	reverse bool,
) db.Iterator {
	txn := bdb.db.NewTransaction(false)
	itOpts := bd.DefaultIteratorOptions
	itOpts.PrefetchValues = false
	itOpts.Reverse = reverse
	it := txn.NewIterator(itOpts)

	// if iterator goes backwards, so start needs to be the biggest of that index start range
	if reverse {
		it.Seek(utils.ConcatBytes(start, maxPKRange))
	} else {
		it.Seek(start)
	}

	return &BadgerIterator{
		txn:            txn,
		it:             it,
		start:          start,
		indexKeyLength: len(start),
		reverse:        reverse,
	}
}

func (it *BadgerIterator) Close() {
	it.it.Close()
	it.txn.Discard()
}

func (it *BadgerIterator) Valid(prefix []byte) bool {
	if len(prefix) != 0 && !it.it.ValidForPrefix(prefix) {
		return false
	}

	return it.it.Valid()
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
	return it.Key()[it.indexKeyLength:]
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
