package depsresolver

import (
	"fmt"
	"reflect"
	"sync"
)

type DepsResolverInstance struct {
	mux       sync.RWMutex
	rmux      sync.RWMutex
	channels  map[reflect.Type]([]chan interface{})
	published map[reflect.Type]interface{}
}

func NewDepsResolver() DepsResolver {
	return &DepsResolverInstance{
		mux:       sync.RWMutex{},
		rmux:      sync.RWMutex{},
		channels:  make(map[reflect.Type][]chan interface{}),
		published: make(map[reflect.Type]interface{}),
	}
}

func (resolver *DepsResolverInstance) SetPredefinedState(entity interface{}) {
	event := getEvent(entity)

	resolver.rmux.Lock()
	resolver.published[event] = entity
	resolver.rmux.Unlock()
}

func (resolver *DepsResolverInstance) Emit(entity interface{}) error {
	resolver.rmux.Lock()
	event := getEvent(entity)

	if _, alreadyEmitted := resolver.published[event]; alreadyEmitted {
		return fmt.Errorf("cannot commit same entity more than once.")
	}

	resolver.published[event] = entity
	resolver.rmux.Unlock()

	for _, subscription := range resolver.channels[event] {
		subscription <- entity
	}

	return nil
}

func (resolver *DepsResolverInstance) GetState() map[string]interface{} {
	state := map[string]interface{}{}
	for key, entity := range resolver.published {
		state[key.Name()] = entity
	}

	return state
}

// Resolver is really just a subscriber
func (resolver *DepsResolverInstance) Resolve(event reflect.Type) interface{} {
	// check if this event has been delivered already
	// in such case, get data directly from resolver.published
	resolver.rmux.RLock()
	if entity, ok := resolver.published[event]; ok {
		resolver.rmux.RUnlock()
		return entity
	}
	resolver.rmux.RUnlock()

	// otherwise start polling on the event channel
	subchannel := make(chan interface{})

	resolver.mux.Lock()
	resolver.channels[event] = append(resolver.channels[event], subchannel)
	resolver.mux.Unlock()

	select {
	case e := <-subchannel:
		return e
	}
}

func (resolver *DepsResolverInstance) Dispose() {
	for _, entity := range resolver.channels {
		for _, channel := range entity {
			close(channel)
		}
	}

	// dispose the previous published data
	resolver.published = make(map[reflect.Type]interface{})
}

func getEvent(entity interface{}) reflect.Type {
	t := reflect.TypeOf(entity)

	switch t.Kind() {
	case reflect.Ptr:
		return t.Elem()
	case reflect.Struct:
		return t
	case reflect.Slice:
		return t
	default:
		panic("Invalid type entity provided")
	}
}
