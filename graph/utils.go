package graph

import (
	"github.com/terra-project/mantle/serdes"
	"github.com/terra-project/mantle/types"
	"reflect"
)

func UnmarshalInternalQueryResult(result *types.GraphQLInternalResult, target interface{}) error {
	targetValue := reflect.Indirect(reflect.ValueOf(target))
	for key, packBytes := range result.Data {
		targetField := targetValue.FieldByName(key)
		targetCache := reflect.New(targetField.Type())

		if unpackErr := serdes.Deserialize(targetField.Type(), packBytes, targetCache.Interface()); unpackErr != nil {
			return unpackErr
		}

		targetField.Set(targetCache.Elem())
	}

	return nil
}
