package types

import "reflect"

type Model interface{}
type ModelType reflect.Type

type Register func(
	Indexer,
	...ModelType,
)
