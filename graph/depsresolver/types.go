package depsresolver

import "reflect"

type DepsResolver interface {
	SetPredefinedState(interface{})
	Emit(interface{}) error

	GetState() map[string]interface{}
	Resolve(reflect.Type) interface{}
	Dispose()
}
