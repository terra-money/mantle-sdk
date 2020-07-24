package committer

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/terra-project/mantle/db"
	"github.com/terra-project/mantle/db/kvindex"
	"github.com/terra-project/mantle/utils"
	"github.com/vmihailenco/msgpack/v5"
)

type CommitterInstance struct {
	db         db.DB
	kvIndexMap kvindex.KVIndexMap
}

func NewCommitter(
	db db.DB,
	kvIndexMap kvindex.KVIndexMap,
) Committer {
	return &CommitterInstance{
		db:         db,
		kvIndexMap: kvIndexMap,
	}
}

func (committer *CommitterInstance) Commit(height uint64, entities ...interface{}) error {
	// create a write batch
	var transaction = 0
	writeBatch := committer.db.Batch()
	defer func() {
		// only flush when all transactions are done
		if transaction == 2 {
			writeBatch.Flush()
		}
		writeBatch.Close()
	}()

	// committer job
	for _, entity := range entities {
		t := utils.GetType(reflect.TypeOf(entity))
		kvIndex := committer.kvIndexMap[t.Name()]
		kvIndexEntries := kvIndex.GetEntries()

		// convert height => BE
		HeightInBE := utils.LeToBe(height)

		// create document key
		documentKey := bytes.NewBuffer(nil)
		documentKey.Write([]byte(t.Name()))
		documentKey.Write(HeightInBE)

		// commit document to db
		documentValue, err := msgpack.Marshal(entity)
		if err != nil {
			return fmt.Errorf(
				"Failed to serialize entity, entityName=%s",
				t.Name(),
			)
		}

		writeBatch.Set(documentKey.Bytes(), documentValue)
		transaction = 1

		// creating indexes
		for _, entry := range kvIndexEntries {
			valuePath := entry.GetEntry().Path
			value := getLeafValue(reflect.ValueOf(entity), valuePath)
			indexKey, err := entry.ResolveKeyType(value.Interface())

			// ResolveKeyType may fail due to unmatching type
			if err != nil {
				return fmt.Errorf(
					"Index creation failed due to unmatching index key type, entityName=%s, indexName=%s, expectedIndexKeyType=%s, actual=%s",
					entry.GetEntityName(),
					entry.GetEntry().Name,
					entry.GetEntry().Type.String(),
					value.Kind().String(),
				)
			}

			// put together an index key
			indexDocumentKey := entry.BuildIndexKey(indexKey, HeightInBE)
			writeBatch.Set(indexDocumentKey, nil)
		}
		transaction = 2
	}

	// all good
	return nil
}

func getLeafValue(entity reflect.Value, valuePath []string) reflect.Value {
	if len(valuePath) > 0 {
		return getLeafValue(entity.FieldByName(valuePath[0]), valuePath[1:])
	}
	return entity
}
