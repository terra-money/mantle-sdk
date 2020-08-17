package committer

import (
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
	var transaction = false
	writeBatch := committer.db.Batch()

	// committer job
	for _, rawEntity := range entities {
		v := reflect.ValueOf(rawEntity)
		t := utils.GetType(v.Type())
		kvIndex, ok := committer.kvIndexMap[t.Name()]
		if !ok {
			return fmt.Errorf("Unknown Entity committed, entityName=%s", t.Name())
		}
		entityName := t.Name()
		kvIndexEntries := kvIndex.GetEntries()

		// creating indexes
		// the entity might be of slice type - convert it to slice anyway for easier handling
		entitiesInSlice := make([]interface{}, 0)

		// in case of slice entity, the leading index path must be adjusted (* has to go)
		// -- flag for that
		isSliceEntity := false
		if t.Kind() == reflect.Slice {
			isSliceEntity = true
			for i := 0; i < v.Len(); i++ {
				entitiesInSlice = append(entitiesInSlice, v.Index(i).Interface())
			}
		} else {
			entitiesInSlice = append(entitiesInSlice, rawEntity)
		}

		seq, errSeq := committer.db.GetSequence([]byte(t.Name()), uint64(len(entitiesInSlice)))
		if errSeq != nil {
			return fmt.Errorf("sequence generation failed")
		}

		// convert some properties to byte beforehand
		entityNameInBytes := []byte(entityName)
		heightInBe := utils.LeToBe(height)

		// commit document(s) to db
		for _, entity := range entitiesInSlice {
			pk, errSeq := seq.Next()
			pkInBE := utils.LeToBe(pk)

			if errSeq != nil {
				return errSeq
			}

			// underlying document key, entityName
			documentKey := utils.BuildDocumentKey(
				entityNameInBytes,
				pkInBE,
			)
			documentValue, err := msgpack.Marshal(entity)
			if err != nil {
				return fmt.Errorf(
					"Failed to serialize entity, entityName=%s",
					t.Name(),
				)
			}

			// write document to db
			writeBatch.Set(documentKey, documentValue)

			// generate height indexes
			heightIndexKey := utils.BuildIndexedDocumentKey(
				entityNameInBytes,
				utils.DocumentHeightIndex,
				heightInBe,
				pkInBE,
			)

			writeBatch.Set(heightIndexKey, nil)

			// generate the rest of indexes
			for _, kviEntry := range kvIndexEntries {
				indexName := kviEntry.GetEntry().Name
				indexNameInBytes := []byte(indexName)
				valuePath := kviEntry.GetEntry().Path

				// if isSliceEntity is true, then this entity is a slice entity
				// the leading '*' must go away
				if isSliceEntity {
					valuePath = valuePath[1:]
				}

				leafValues, leafValuesErr := getLeafValues(reflect.ValueOf(entity), valuePath)

				if leafValuesErr != nil {
					return leafValuesErr
				}

				for _, leafValue := range leafValues {
					index, indexErr := kviEntry.ResolveKeyType(leafValue.Interface())

					if indexErr != nil {
						return fmt.Errorf(
							"index creation failed due to unmatching index key type, entityName=%s, indexName=%s, expectedIndexKeyType=%s, actual=%s",
							kviEntry.GetEntityName(),
							kviEntry.GetEntry().Name,
							kviEntry.GetEntry().Type.String(),
							leafValue.Kind().String(),
						)
					}

					indexKey := utils.BuildIndexedDocumentKey(
						entityNameInBytes,
						indexNameInBytes,
						index,
						pkInBE,
					)

					writeBatch.Set(indexKey, nil)
				}
			}
		}

		seqReleaseErr := seq.Release()
		if seqReleaseErr != nil {
			return fmt.Errorf("sequence release failed")
		}
	}

	transaction = true

	defer func() {
		// only flush when all transactions are done
		if transaction == true {
			writeBatch.Flush()
		}
		writeBatch.Close()
	}()

	// all good
	return nil
}

func getLeafValues(entity reflect.Value, valuePath []string) ([]reflect.Value, error) {
	values := make([]reflect.Value, 0)
	if err := _getLeafValues(entity, valuePath, &values); err != nil {
		return nil, err
	}

	return values, nil
}

func _getLeafValues(entity reflect.Value, valuePath []string, values *[]reflect.Value) error {
	if len(valuePath) > 0 {
		currentPath := valuePath[0]
		if currentPath == "*" {
			switch entity.Type().Kind() {
			case reflect.Slice, reflect.Array:
				len := entity.Len()
				for i := 0; i < len; i++ {
					if err := _getLeafValues(entity.Index(i), valuePath[1:], values); err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("entity is not slice yet path * is given")
			}
		} else {
			return _getLeafValues(entity.FieldByName(currentPath), valuePath[1:], values)
		}
	}

	*values = append(*values, entity)

	return nil
}
