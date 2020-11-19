package leveldb

import (
	"bytes"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbIterator "github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
	tmdb "github.com/tendermint/tm-db"
	"github.com/terra-project/mantle-sdk/db"
	"github.com/terra-project/mantle-sdk/db/leveldb/tm_adapter"
)

type LevelDB struct {
	db *leveldb.DB
	path string
	cosmosdb *tm_adapter.GoLevelDB
}

var maxPKRange = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

func NewLevelDB(path string) db.DB {
	dbInstance := &LevelDB{
		path: path,
	}

	dbInstance.db = dbInstance.open(path)

	return dbInstance
}

func (ldb *LevelDB) open(path string) *leveldb.DB {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		panic(err)
	}

	return db
}

func (ldb *LevelDB) Compact() error {
	// close & reopen db
	ldb.db.Close()
	ldb.db = ldb.open(ldb.path)
	ldb.cosmosdb.SetDB(ldb.db)

	return nil
}

func (ldb *LevelDB) GetCosmosAdapter() tmdb.DB {
	cosmosdb := tm_adapter.NewLevelDBCosmosAdapter(ldb.db)
	ldb.cosmosdb = cosmosdb
	return cosmosdb
}

func (ldb *LevelDB) GetDB() *leveldb.DB{
	return ldb.db
}

func (ldb *LevelDB) Get(key []byte) ([]byte, error) {
	return ldb.db.Get(key, nil)
}

func (ldb *LevelDB) Set(key, data []byte) error {
	return ldb.db.Put(key, data, nil)
}

func (ldb *LevelDB) Delete(key []byte) error {
	return ldb.db.Delete(key, nil)
}

func (ldb *LevelDB) GetSequence(key []byte, bandwidth uint64) (db.Sequence, error) {
	return NewLevelDBSequence(ldb, key, bandwidth)
}

func (ldb *LevelDB) Close() error {
	return ldb.db.Close()
}

type LevelIterator struct {
	it             leveldbIterator.Iterator
	start          []byte
	indexKeyLength int
	reverse        bool
}

func (ldb *LevelDB) Iterator(
	start []byte,
	reverse bool,
) db.Iterator {
	var it leveldbIterator.Iterator

	// if iterator goes backwards, so start needs to be the biggest of that index start range
	if reverse {
		it = ldb.db.NewIterator(&util.Range{Start: nil, Limit: start}, nil)
		it.Last()
	} else {
		it = ldb.db.NewIterator(&util.Range{Start: start, Limit: nil}, nil)
		it.First()
	}

	return &LevelIterator{
		it:             it,
		start:          start,
		indexKeyLength: len(start),
		reverse:        reverse,
	}
}

// no index-only iterator for ldb
func (ldb *LevelDB) IndexIterator(
	start []byte,
	reverse bool,
) db.Iterator {
	return ldb.Iterator(start, reverse)
}

func (it *LevelIterator) Close() {
	it.it.Release()
}

func (it *LevelIterator) Valid(prefix []byte) bool {
	if len(prefix) != 0 && !bytes.HasPrefix(it.it.Key(), prefix) {
		return false
	}
	return it.it.Valid()
}

func (it *LevelIterator) Next() {
	if it.reverse {
		it.it.Prev()
	} else {
		it.it.Next()
	}
}

func (it *LevelIterator) Key() []byte {
	return it.it.Key()
}

func (it *LevelIterator) Value() []byte {
	return it.it.Value()
}

func (it *LevelIterator) DocumentKey() []byte {
	return it.Key()[it.indexKeyLength:]
}

type Batch struct {
	batch *leveldb.Batch
	db *leveldb.DB
}

func (ldb *LevelDB) Batch() db.Batch {
	return &Batch{
		batch: new(leveldb.Batch),
		db: ldb.db,
	}
}

func (batch *Batch) Set(key, data []byte) error {
	batch.batch.Put(key, data)
	return nil
}

func (batch *Batch) Delete(key []byte) error {
	batch.batch.Delete(key)
	return nil
}

func (batch *Batch) Flush() error {
	return batch.db.Write(batch.batch, nil)
}

func (batch *Batch) Close() {
	// noop
	// batch.batch.Dump()
}
