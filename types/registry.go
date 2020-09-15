package types

import (
	"reflect"
)

type Model reflect.Type
type Register func(
	Indexer,
	...Model,
)
