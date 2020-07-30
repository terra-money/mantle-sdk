package depsresolver

import (
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

func (resolver *DepsResolverInstance) Emit(entity interface{}) {
	resolver.rmux.Lock()
	event := getEvent(entity)
	resolver.published[event] = entity
	resolver.rmux.Unlock()

	for _, subscription := range resolver.channels[event] {
		subscription <- entity
	}
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

	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}

	if t.Kind() == reflect.Struct {
		return t
	}

	panic("Invalid type entity provided")
}
