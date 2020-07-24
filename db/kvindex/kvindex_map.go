package kvindex

type KVIndexMap map[string]*KVIndex

func NewKVIndexMap(kvindexes ...*KVIndex) KVIndexMap {
	kvindexMap := KVIndexMap{}
	for _, kvi := range kvindexes {
		kvindexMap[kvi.entityType.Name()] = kvi
	}

	return kvindexMap
}
