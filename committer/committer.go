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
		entityName := t.Name()

		kvIndex, ok := committer.kvIndexMap[entityName]
		if !ok {
			return fmt.Errorf("Unknown Entity committed, entityName=%s", t.Name())
		}
		kvIndexEntries := kvIndex.Entries()

		// convert some properties to byte beforehand
		heightInBe := utils.LeToBe(height)
		var commit CommitFunc = func(key, value []byte) error {
			fmt.Println(string(utils.ConcatBytes([]byte(entityName), key)))
			return writeBatch.Set(utils.ConcatBytes([]byte(entityName), key), value)
		}
		var getSequence = NewSequenceGenerator(entityName, kvIndex, committer.db)

		// in case of slice entity, the leading index path must be adjusted (* has to go)
		// -- flag for that
		commitTargets, flattenErr := flattenCommitTargets(getSequence, v)
		if flattenErr != nil {
			return flattenErr
		}

		// commit document(s) to db
		for _, commitTarget := range commitTargets {
			// commit primary documents
			if commitErr := commitTarget.commit(commit); commitErr != nil {
				return commitErr
			}

			// commit height index. skip when pk model
			if !kvIndex.IsPrimaryKeyedModel() {
				if commitErr := commitTarget.commitIndex(commit, "Height", heightInBe); commitErr != nil {
					return commitErr
				}
			}

			// commit all indexes
			for _, kviEntry := range kvIndexEntries {
				// skip Height
				if kviEntry.Name() == "Height" {
					continue
				}

				indexName := kviEntry.Name()

				//// if entity is either slice or map, the leading '*' must go away
				//if commitTarget.entityType == EntityTypeSlice {
				//	valuePath = valuePath[1:]
				//}

				leafValues, leafValuesErr := kviEntry.ResolveIndexKeySingle(commitTarget.data)
				if leafValuesErr != nil {
					return leafValuesErr
				}

				for _, leafValue := range leafValues {
					// commit index
					indexValue, indexValueErr := utils.ConvertToLexicographicBytes(leafValue)

					if indexValueErr != nil {
						return fmt.Errorf(
							"index creation failed due to unmatching index key type, entityName=%s, indexName=%s, expectedIndexKeyType=%s, actual=%s",
							entityName,
							kviEntry.Name(),
							kviEntry.Type().Name(),
							reflect.TypeOf(indexValue).String(),
						)
					}

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
	pk   []byte
	data interface{}
}

// TODO: don't use reflect somehow
func flattenCommitTargets(
	getSequence sequenceGetter,
	v reflect.Value,
) (commitTargets []commitTarget, err error) {
	// create sequencer
	nextSeq, errSeq := getSequence(v.Interface())
	if errSeq != nil {
		return nil, errSeq
	}

	// flatten entities
	var entitiesInterface []interface{}
	switch v.Type().Kind() {
	case reflect.Slice:
		len := v.Len()
		entitiesInterface = make([]interface{}, len)
		for i := 0; i < len; i++ {
			entitiesInterface[i] = v.Index(i).Interface()
		}
	default:
		entitiesInterface = make([]interface{}, 1)
		entitiesInterface[0] = v.Interface()
	}

	// make commit targets
	length := len(entitiesInterface)
	commitTargets = make([]commitTarget, length)

	for i := 0; i < length; i++ {
		pk, pkErr := nextSeq()
		if pkErr != nil {
			return nil, pkErr
		}

		commitTargets[i] = commitTarget{
			pk:   pk,
			data: entitiesInterface[i],
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

//
type sequenceGetter func(entity interface{}) (sequenceNext, error)
type sequenceNext func() ([]byte, error)

func NewSequenceGenerator(modelName string, kvindex *kvindex.KVIndex, persistence db.DB) sequenceGetter {
	// if primary tag is set
	if kvindex.IsPrimaryKeyedModel() {
		return func(entity interface{}) (sequenceNext, error) {
			hashKeys, hashKeysErr := kvindex.ResolvePrimaryKey(entity)
			if hashKeysErr != nil {
				return nil, hashKeysErr
			}
			hashIdx := 0

			return func() ([]byte, error) {
				key, err := utils.ConvertToLexicographicBytes(hashKeys[hashIdx])
				if err != nil {
					return nil, err
				}

				hashIdx++

				return key, nil
			}, nil
		}
	}

	//
	return func(entity interface{}) (sequenceNext, error) {
		v := reflect.ValueOf(entity)
		t := v.Type()

		var len uint64 = 1

		//
		if t.Kind() == reflect.Slice {
			len = uint64(v.Len())
		}

		sequencer, sequencerErr := persistence.GetSequence([]byte(t.Name()), len)
		if sequencerErr != nil {
			return nil, sequencerErr
		}

		return func() ([]byte, error) {
			seq, err := sequencer.Next()
			if err != nil {
				return nil, err
			}
			return utils.LeToBe(seq), nil
		}, nil
	}
}
