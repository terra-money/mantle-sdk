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

const (
	EntityTypeSingular = 0
	EntityTypeSlice    = 1
	EntityTypeMap      = 2
)

type (
	GetSequenceFunc func(length uint64) (db.Sequence, error)
	CommitFunc      func(key, data []byte) error
)

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

		// convert some properties to byte beforehand
		heightInBe := utils.LeToBe(height)
		var commit CommitFunc = func(key, value []byte) error {
			return writeBatch.Set(utils.ConcatBytes([]byte(entityName), key), value)
		}
		var getSequence GetSequenceFunc = func(length uint64) (db.Sequence, error) {
			return committer.db.GetSequence([]byte(t.Name()), length)
		}

		// in case of slice entity, the leading index path must be adjusted (* has to go)
		// -- flag for that
		commitTargets, flattenErr := flattenCommitTargets(getSequence, v, t)
		if flattenErr != nil {
			return flattenErr
		}

		// commit document(s) to db
		for _, commitTarget := range commitTargets {
			// commit primary documents
			if commitErr := commitTarget.commit(commit); commitErr != nil {
				return commitErr
			}

			// commit height index
			if commitErr := commitTarget.commitIndex(commit, "Height", heightInBe); commitErr != nil {
				return commitErr
			}

			// commit all indexes
			for _, kviEntry := range kvIndexEntries {
				// skip Height
				if kviEntry.GetEntry().Name == "Height" {
					continue
				}

				indexName := kviEntry.GetEntry().Name
				valuePath := kviEntry.GetEntry().Path

				// if entity is either slice or map, the leading '*' must go away
				if commitTarget.entityType == EntityTypeSlice || commitTarget.entityType == EntityTypeMap {
					valuePath = valuePath[1:]
				}

				leafValues, leafValuesErr := getLeafValues(reflect.ValueOf(commitTarget.data), valuePath)
				if leafValuesErr != nil {
					return leafValuesErr
				}

				for _, leafValue := range leafValues {
					indexValue, indexErr := kviEntry.ResolveKeyType(leafValue.Interface())

					if indexErr != nil {
						return fmt.Errorf(
							"index creation failed due to unmatching index key type, entityName=%s, indexName=%s, expectedIndexKeyType=%s, actual=%s",
							kviEntry.GetEntityName(),
							kviEntry.GetEntry().Name,
							kviEntry.GetEntry().Type.String(),
							leafValue.Kind().String(),
						)
					}

					// commit index
					if commitErr := commitTarget.commitIndex(commit, indexName, indexValue); commitErr != nil {
						return commitErr
					}
				}
			}
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

type commitTarget struct {
	entityType int
	pk         []byte
	data       interface{}
}

// TODO: don't use reflect somehow
func flattenCommitTargets(getSequence GetSequenceFunc, v reflect.Value, t reflect.Type) (commitTargets []commitTarget, err error) {
	switch t.Kind() {
	case reflect.Slice:
		length := v.Len()
		seq, errSeq := getSequence(uint64(length))
		if errSeq != nil {
			return nil, errSeq
		}
		defer seq.Release()

		commitTargets = make([]commitTarget, length)

		for i := 0; i < v.Len(); i++ {
			pk, pkErr := seq.Next()
			if pkErr != nil {
				return nil, pkErr
			}

			commitTargets[i] = commitTarget{
				entityType: EntityTypeSlice,
				pk:         utils.LeToBe(pk),
				data:       v.Index(i).Interface(),
			}
		}
	case reflect.Map:
		length := v.Len()
		commitTargets = make([]commitTarget, length)

		for i, mapKey := range v.MapKeys() {
			commitTargets[i] = commitTarget{
				entityType: EntityTypeMap,
				pk:         []byte(mapKey.String()),
				data:       v.MapIndex(mapKey).Interface(),
			}
		}
	default:
		seq, errSeq := getSequence(uint64(1))
		if errSeq != nil {
			return nil, errSeq
		}
		defer seq.Release()

		pk, pkErr := seq.Next()
		if pkErr != nil {
			return nil, pkErr
		}

		commitTargets = []commitTarget{
			{
				entityType: EntityTypeSingular,
				pk:         utils.LeToBe(pk),
				data:       v.Interface(),
			},
		}
	}

	return commitTargets, nil
}

func (ct commitTarget) commit(commit CommitFunc) error {
	documentKey := utils.BuildDocumentKey(nil, ct.pk)
	documentValue, err := msgpack.Marshal(ct.data)
	if err != nil {
		return fmt.Errorf("failed to serialize entity")
	}

	// write document to db
	return commit(documentKey, documentValue)
}

func (ct commitTarget) commitIndex(commit CommitFunc, indexName string, indexValue []byte) error {
	return commit(utils.BuildIndexedDocumentKey(
		nil,
		[]byte(indexName),
		indexValue,
		ct.pk,
	), nil)
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
			case reflect.Map:
				for _, key := range entity.MapKeys() {
					if err := _getLeafValues(entity.MapIndex(key), valuePath[1:], values); err != nil {
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
