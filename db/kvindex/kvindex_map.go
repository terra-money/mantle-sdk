package kvindex

type KVIndexMap map[string]*KVIndex

func NewKVIndexMap(kvIndexes ...*KVIndex) KVIndexMap {
	kvindexMap := KVIndexMap{}
	for _, kvi := range kvIndexes {
		kvindexMap[kvi.ModelName()] = kvi
	}

	return kvindexMap
}
