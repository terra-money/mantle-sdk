package kvindex

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/terra-project/mantle/utils"
)

const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

type KVIndex struct {
	entityType      reflect.Type
	entityName      string
	indexMapEntries map[string]KVIndexEntry
}

type KVIndexKeyTypeResolver func(val interface{}) ([]byte, error)

func NewKVIndex(
	entityType reflect.Type,
) *KVIndex {
	entityType = utils.GetType(entityType)

	entityName := entityType.Name()
	entries := createIndexMap(entityType)
	indexMapEntries := make(map[string]KVIndexEntry)

	for _, entry := range entries {
		if _, ok := indexMapEntries[entry.Name]; ok {
			panic(fmt.Errorf("Duplicate index is disallowed, %s:%s", entityName, entry.Name))
		}

		indexMapEntries[entry.Name] = KVIndexEntry{
			entityName: entityName,
			entry:      entry,
			resolver:   createSecondaryIndexGetter(entry.Type),
		}
	}

	return &KVIndex{
		entityType:      entityType,
		entityName:      entityType.Name(),
		indexMapEntries: indexMapEntries,
	}
}

func (kvi *KVIndex) HasIndex() bool {
	return len(kvi.indexMapEntries) != 0
}

func (kvi *KVIndex) GetEntries() map[string]KVIndexEntry {
	return kvi.indexMapEntries
}

func (kvi *KVIndex) GetIndexEntry(indexName string) *KVIndexEntry {
	kvIndexEntry, ok := kvi.indexMapEntries[indexName]
	if !ok {
		return nil
	}

	return &kvIndexEntry
}

func (kvi *KVIndex) GetPrefix(indexName string) []byte {
	if _, ok := kvi.indexMapEntries[indexName]; !ok {
		panic(fmt.Errorf("Index %s:%s does NOT exist.", kvi.entityName, indexName))
	}

	entityName := kvi.entityName

	newPrefix := []byte{}
	newPrefix = append(newPrefix, []byte(entityName)...)
	newPrefix = append(newPrefix, []byte(indexName)...)

	return newPrefix
}

func (kvi *KVIndex) GetCursor(indexName string, cursor interface{}) ([]byte, error) {
	if cursor == nil {
		return nil, fmt.Errorf("GetCursor expected cursor to be non-null type, %s:%s", kvi.entityName, indexName)
	}

	if reflect.TypeOf(cursor).Kind() == reflect.Ptr {
		return nil, fmt.Errorf("GetCursor exptected cursor to be non-ptr type, %s:%s", kvi.entityName, indexName)
	}

	indexMap := kvi.indexMapEntries[indexName]
	cursorType := reflect.TypeOf(cursor).Kind()

	if cursorType != indexMap.entry.Type {
		return nil, fmt.Errorf(
			"GetCursor received cursor different from index definition. Expected %s, received %s, %s:%s",
			indexMap.entry.Type.String(),
			cursorType.String(),
			kvi.entityName,
			indexName,
		)
	}

	prefix := kvi.GetPrefix(indexName)
	canonicalCursor, err := indexMap.resolver(cursor)

	if err != nil {
		return nil, err
	}

	if len(canonicalCursor) == 0 {
		return nil, fmt.Errorf("Calculaing cursor failed, %s:%s(%s)", kvi.entityName, indexName, cursor)
	}

	return append(prefix, canonicalCursor...), nil
}

func (kvi *KVIndex) BuildIndexKey(indexName string, cursor interface{}, documentCursor []byte) ([]byte, error) {
	indexCursor, err := kvi.GetCursor(indexName, cursor)
	if err != nil {
		return nil, err
	}

	return append(indexCursor, documentCursor...), nil
}

func createSecondaryIndexGetter(k reflect.Kind) KVIndexKeyTypeResolver {
	switch k {
	// if string, use string as is for index (lexicographic order)
	case reflect.String:
		return func(val interface{}) ([]byte, error) {
			valString, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("index type conversion failed. Expected String variants, got %s", reflect.TypeOf(val).Kind().String())
			}
			return []byte(valString), nil
		}
		// if number types, take big endian. flush them to uint64
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(val interface{}) ([]byte, error) {
			key := make([]byte, 8)
			n, err := utils.GetUint64FromWhatever(val)
			if err != nil {
				return nil, fmt.Errorf("Index type conversion failed. Expected Uint variants, got %s (suberror = %s)", reflect.TypeOf(val).Kind().String(), err)
			}
			binary.BigEndian.PutUint64(key, n)
			return key, nil
		}
		// for signed integers, take care of msb
		// TODO fix me
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(val interface{}) ([]byte, error) {
			key := make([]byte, 8)
			n, err := strconv.ParseInt(fmt.Sprintf("%s", val), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("Index type conversion failed. Expected Int variants, got %s", reflect.TypeOf(val).Kind().String())
			}
			binary.BigEndian.PutUint64(key, uint64(n))
			return key, nil
		}
		// if float type, take a stringified version of it. --> 8width, 8decimal points
		// TODO: fix me
	case reflect.Float32, reflect.Float64:
		return func(val interface{}) ([]byte, error) {
			return []byte(fmt.Sprintf("%8.8f", val)), nil
		}
	default:
		panic(fmt.Errorf("This type of data is disallowed for indexing: %s", k.String()))
	}
}

///
type KVIndexEntry struct {
	entityName string
	entry      IndexMapEntry
	resolver   KVIndexKeyTypeResolver
}

func (kvimap KVIndexEntry) GetEntityName() string {
	return kvimap.entityName
}

func (kvimap KVIndexEntry) GetEntry() IndexMapEntry {
	return kvimap.entry
}

func (kvimap KVIndexEntry) ResolveKeyType(key interface{}) ([]byte, error) {
	// graphql query arguments DO retain types, but in case of @range or others,
	// query params are read as strings because of regex matcher.
	// need to check further if type of `key` can be safely converted.
	//
	// TODO: investigate how to do this properly
	return kvimap.resolver(key)
}

// simpler BuildIndexKey if you know which kviMap to use
func (kviMap KVIndexEntry) BuildIndexKey(indexKeyInBytes []byte, documentKeyInBytes []byte) []byte {
	indexKey := bytes.NewBuffer(nil)
	indexKey.Write([]byte(kviMap.entityName))
	indexKey.Write([]byte(kviMap.entry.Name))
	indexKey.Write(indexKeyInBytes)
	indexKey.Write(documentKeyInBytes)

	return indexKey.Bytes()
}

////////////////////////////////////
const KEY_INDEX = "index"

type IndexMapEntry struct {
	Type reflect.Kind
	Name string
	Path []string
}

func createIndexMap(t reflect.Type) []IndexMapEntry {
	t = utils.GetType(t)
	indexMapSlice := []IndexMapEntry{}

	createIndexMapIter(t, &indexMapSlice, []string{}, false, "")

	return indexMapSlice
}

func createIndexMapIter(t reflect.Type, indexMapSlice *[]IndexMapEntry, path []string, isIndexed bool, indexName string) {
	k := t.Kind()

	switch k {
	case reflect.Array, reflect.Slice:
		createIndexMapIter(t.Elem(), indexMapSlice, append(path, "*"), isIndexed, "")
	case reflect.Ptr:
		createIndexMapIter(t.Elem(), indexMapSlice, path, isIndexed, "")
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ft := t.Field(i)

			tags, _ := ft.Tag.Lookup(utils.MantleKeyTag)
			tagsSplit := strings.Split(tags, ",")
			indexTag, indexed := sliceContainsString(tagsSplit, KEY_INDEX)

			// disallow if
			// - another struct appears as a direct child of this struct, 뭉
			// - index tag is set.
			if indexed && ft.Type.Kind() == reflect.Struct {
				panic(fmt.Errorf("Structs are disallowed for kvindex: %s", path))
			}

			indexName := ft.Name
			fmt.Sscanf(indexTag, "index=%s", &indexName)

			createIndexMapIter(ft.Type, indexMapSlice, append(path, ft.Name), indexed, indexName)
		}

	default:
		// noop if not indexed field
		if !isIndexed {
			return
		}

		sanitize := regexp.MustCompile(utils.GraphQLAllowedCharactersRegex)
		pass := sanitize.MatchString(indexName)

		if !pass {
			panic(fmt.Errorf("Index name contains disallowed characters: %s", indexName))
		}

		entry := IndexMapEntry{
			Type: k,
			Name: indexName,
			Path: path,
		}

		*indexMapSlice = append(*indexMapSlice, entry)
	}

}

func sliceContainsString(haystack []string, needle string) (string, bool) {
	for i := 0; i < len(haystack); i++ {
		if strings.Contains(haystack[i], needle) {
			return haystack[i], true
		}
	}

	return "", false
}
