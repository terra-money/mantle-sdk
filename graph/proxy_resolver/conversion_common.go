package proxy_resolver

import (
	"fmt"
	"github.com/graphql-go/graphql"
)

type (
	Definitions map[string]TypeDescriptor
)

var (
	errNoRoot              = fmt.Errorf("proxy resolver context is not root")
	errUnknownArgumentType = func(name string) error { return fmt.Errorf("unknown argument type: %s", name) }
	errUnknownCustomScalar = func(name string) error { return fmt.Errorf("unknown scalar type: %s", name) }
	errInvalidSource       = func(path []interface{}) error {
		return fmt.Errorf("invalid proxy resolver context or source, path: %v", path)
	}
	errNoName = func(whatever interface{}) error { return fmt.Errorf("type without name is given: %v", whatever) }
)

var objectCache = make(map[string]graphql.Input) // use Input for both Input and Output

func getTypeFromCache(name string) graphql.Input {
	return objectCache[name]
}
func setTypeInCache(name string, input graphql.Input) {
	if getTypeFromCache(name) != nil {
		panic(fmt.Errorf("duplicate input cache set: %s", name))
	}
	objectCache[name] = input
}
