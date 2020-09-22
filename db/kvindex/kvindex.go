package kvindex

import (
	"fmt"
	"github.com/terra-project/mantle/types"
	"github.com/terra-project/mantle/utils"
	"reflect"
	"regexp"
)

type KVIndex struct {
	primaryKeyEntry IndexEntry
	isSliceModel    bool
	modelType       reflect.Type
	modelName       string
	indexEntryMap   map[string]IndexEntry
}

func NewKVIndex(model types.Model) (*KVIndex, error) {
	modelType := utils.GetType(model)
	modelName := modelType.Name()
	indexEntryMap, indexEntryMapErr := createIndexMap(modelType)
	if indexEntryMapErr != nil {
		return nil, indexEntryMapErr
	}

	// check if model is of slice/array type
	var isSliceModel = false
	if modelType.Kind() == reflect.Slice {
		isSliceModel = true
	}

	// make this KVIndex hashentity if struct has key tag defined
	var primaryKeyEntry IndexEntry
	for _, entry := range indexEntryMap {
		if entry.isPrimaryKey && primaryKeyEntry.indexName != "" {
			return nil, fmt.Errorf("duplicate hash key is disallowed, modelName=%s", modelName)
		}

		if entry.isPrimaryKey {
			primaryKeyEntry = entry
		}
	}

	return &KVIndex{
		primaryKeyEntry: primaryKeyEntry,
		isSliceModel:    isSliceModel,
		modelType:       modelType,
		modelName:       modelName,
		indexEntryMap:   indexEntryMap,
	}, nil
}

func (kvi *KVIndex) IsPrimaryKeyedModel() bool {
	return kvi.primaryKeyEntry.indexName != ""
}

func (kvi *KVIndex) IsSliceModel() bool {
	return kvi.isSliceModel
}

func (kvi *KVIndex) Entries() map[string]IndexEntry {
	return kvi.indexEntryMap
}

func (kvi *KVIndex) Entry(indexName string) (IndexEntry, bool) {
	entry, ok := kvi.indexEntryMap[indexName]
	return entry, ok
}

func (kvi *KVIndex) ModelName() string {
	return kvi.modelName
}

func (kvi *KVIndex) ResolvePrimaryKey(entity interface{}) ([]interface{}, error) {
	values, err := kvi.primaryKeyEntry.ResolveIndexKey(entity)
	if err != nil {
		return nil, err
	}

	return values, nil
}

type IndexEntry struct {
	isPrimaryKey bool
	indexType    reflect.Type
	indexName    string
	indexPath    []string
}

func (indexEntry IndexEntry) Name() string {
	return indexEntry.indexName
}

func (indexEntry IndexEntry) Type() reflect.Type {
	return indexEntry.indexType
}

func (indexEntry IndexEntry) ResolveIndexKey(entity interface{}) ([]interface{}, error) {
	entityValue := reflect.ValueOf(entity)
	values, err := getLeafValues(entityValue, indexEntry.indexPath)
	if err != nil {
		return nil, err
	}

	actualValues := make([]interface{}, len(values))
	for i, value := range values {
		actualValues[i] = value.Interface()
	}

	return actualValues, nil
}

func (indexEntry IndexEntry) ResolveIndexKeySingle(entity interface{}) ([]interface{}, error) {
	entityValue := reflect.ValueOf(entity)
	indexPath := indexEntry.indexPath
	if indexPath[0] == "*" {
		indexPath = indexPath[1:]
	}
	values, err := getLeafValues(entityValue, indexPath)
	if err != nil {
		return nil, err
	}

	actualValues := make([]interface{}, len(values))
	for i, value := range values {
		actualValues[i] = value.Interface()
	}

	return actualValues, nil
}

// private subroutines
func createIndexMap(t reflect.Type) (map[string]IndexEntry, error) {
	t = utils.GetType(t)
	indexEntryMap := make(map[string]IndexEntry)

	if err := createIndexMapIter(
		t,
		indexEntryMap,
		[]string{},
		false,
		"",
		false,
	); err != nil {
		return nil, err
	}

	return indexEntryMap, nil
}

func createIndexMapIter(
	t reflect.Type,
	indexEntryMap map[string]IndexEntry,
	path []string,
	isIndexed bool,
	indexName string,
	isPrimaryKey bool,
) error {
	k := t.Kind()

	switch k {
	case reflect.Array, reflect.Slice, reflect.Map:
		return createIndexMapIter(t.Elem(), indexEntryMap, append(path, "*"), isIndexed, "", isPrimaryKey)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)
			isIndexed := false
			indexName := ""
			isPrimaryKey := false

			modelTag, modelTagExists := utils.NewTagMap(ft.Tag.Lookup(utils.MantleModelTag))
			if modelTagExists {
				// disallow if subtype is struct and ANY tag is set
				if ft.Type.Kind() == reflect.Struct {
					return fmt.Errorf("structs are disallowed for kvindex: %s", path)
				}

				// if index option is set, index this field
				if _, indexed := modelTag.Option("index"); indexed {
					isIndexed = true
					indexName = ft.Name
				}

				if _, isPrimary := modelTag.Option("primary"); isPrimary {
					isPrimaryKey = true
				}
			}

			if err := createIndexMapIter(ft.Type, indexEntryMap, append(path, ft.Name), isIndexed, indexName, isPrimaryKey); err != nil {
				return err
			}
		}

		return nil

	default:
		// noop if not indexed field
		if !isIndexed {
			return nil
		}

		sanitize := regexp.MustCompile(utils.GraphQLAllowedCharactersRegex)
		pass := sanitize.MatchString(indexName)

		if !pass {
			return fmt.Errorf("index name contains disallowed characters: %s", indexName)
		}

		if _, ok := indexEntryMap[indexName]; ok {
			return fmt.Errorf("duplicate index is disallowed, indexName=%s", indexName)
		}

		indexEntryMap[indexName] = IndexEntry{
			isPrimaryKey: isPrimaryKey,
			indexType:    t,
			indexName:    indexName,
			indexPath:    path,
		}

		return nil
	}
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
			//case reflect.Map:
			//	for _, key := range entity.MapKeys() {
			//		if err := _getLeafValues(entity.MapIndex(key), valuePath[1:], values); err != nil {
			//			return err
			//		}
			//	}
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
