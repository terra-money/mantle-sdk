package graph

import (
	"github.com/terra-project/mantle-sdk/serdes"
	"github.com/terra-project/mantle-sdk/types"
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

type Thunk func() (interface{}, error)
type ThunkResult struct {
	data interface{}
	err error
}
func CreateThunk(thunk Thunk) (Thunk, error) {
	ch := make(chan *ThunkResult, 1)

	go func() {
		defer close(ch)
		res, err := thunk()
		if err != nil {
			ch <- &ThunkResult{data: nil, err: err}
		} else {
			ch <- &ThunkResult{data: res, err: nil}
		}
	}()

	return func() (interface{}, error) {
		r := <-ch
		return r.data, r.err
	}, nil
}