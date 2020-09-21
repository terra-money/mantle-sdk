package graph

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/types"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"
)

// TODO: Make a better version of this or scrap
func UnmarshalGraphQLResult(result *graphql.Result, target interface{}) error {
	// leave BaseState alone
	res, err := json.Marshal(result.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(res, target)

	if err != nil {
		return err
	}

	return nil
}

func UnmarshalInternalQueryResult(result *types.GraphQLInternalResult, target interface{}) error {
	targetValue := reflect.Indirect(reflect.ValueOf(target))
	for key, packBytes := range result.Data {
		targetField := targetValue.FieldByName(key)
		targetCache := reflect.New(targetField.Type())
		if unpackErr := msgpack.Unmarshal(packBytes, targetCache.Interface()); unpackErr != nil {
			return unpackErr
		}

		targetField.Set(targetCache.Elem())
	}

	return nil
}
