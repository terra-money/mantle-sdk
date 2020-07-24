package depsresolver

import "reflect"

type DepsResolver interface {
	SetPredefinedState(interface{})
	Emit(interface{})
	GetState() map[string]interface{}
	Resolve(reflect.Type) interface{}
	Dispose()
}
