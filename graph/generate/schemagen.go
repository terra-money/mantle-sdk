package generate

import (
	"fmt"
	"github.com/terra-project/mantle/querier"
	"github.com/vmihailenco/msgpack/v5"
	"reflect"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle/graph/depsresolver"
	"github.com/terra-project/mantle/utils"
)

var goTypeToGraphqlType = map[string]graphql.Type{
	"string":     graphql.String,
	"rune":       graphql.Int,
	"int":        graphql.Int,
	"byte":       graphql.Int,
	"int8":       graphql.Int,
	"int16":      graphql.Int,
	"int32":      graphql.Int,
	"int64":      graphql.Int,
	"uint8":      graphql.Int,
	"uint16":     graphql.Int,
	"uint32":     graphql.Int,
	"uint64":     graphql.Int,
	"bool":       graphql.Boolean,
	"float32":    graphql.Float,
	"float64":    graphql.Float,
	"complex64":  graphql.Float,
	"complex128": graphql.Float,
}

func GenerateGraphResolver(modelType reflect.Type) (string, *graphql.Field) {
	t := utils.GetType(modelType)

	entityName := t.Name()
	fields := graphql.Fields{}

	for i := 0; i < t.NumField(); i++ {
		child := t.Field(i)
		iterate(entityName, fields, child.Name, child.Type)
	}

	entitiyRoot := graphql.Field{
		Name: entityName,
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   entityName,
			Fields: fields,
		}),

		// Resolve function here defines a root entity resolver.
		// it should resolve data dependency through subscriber
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// if an object is requested with args,
			// then the querier is not expecting this entity
			// to be resolved in this round.
			// grab from querier
			args := p.Args

			// in this schema resolver, only single object is to be handled.
			// for list type resolvers, see CreateListSchemaBuilder
			if args != nil {
				q, ok := p.Context.Value(utils.QuerierKey).(querier.Querier)
				if !ok {
					panic(fmt.Sprintf("Querier is either cleared or not set, in ResolverFunc for %s", entityName))
				}

				for indexKey, indexParam := range args {
					queryResolver, err := q.Build(entityName, indexKey, indexParam)
					if err != nil {
						return nil, err
					}

					//
					queryResolverIterator, err := queryResolver.Resolve()
					if err != nil {
						return nil, err
					}

					key := []byte{}
					for queryResolverIterator.Valid() {
						key = queryResolverIterator.Key()
						queryResolverIterator.Next()
					}

					defer queryResolverIterator.Close()

					doc, err := q.Get(key)
					docInterface := reflect.New(t)

					if err := msgpack.Unmarshal(doc, docInterface.Interface()); err != nil {
						return nil, err
					}

					return docInterface, err
				}
			}

			// if current round of graphql resolve is a self-referencing (i.e. awaiting on itself to be resolved)
			// we can't process it here.
			// Must error.
			var dependencies, ok = p.Context.Value(utils.DependenciesKey).(utils.DependenciesKeyType)
			var isAwaitedDependency = p.Args["Height"] == nil
			var isSelfReferencing = dependencies[entityName] == true

			if isAwaitedDependency && isSelfReferencing {
				return nil, fmt.Errorf("Self reference is disallowed. entityName=%s", entityName)
			}

			depsResolver, ok := p.Context.Value(utils.DepsResolverKey).(depsresolver.DepsResolver)
			if !ok {
				panic(fmt.Sprintf("DepsResolver is either cleared or not set, in ResolveFunc for %s", entityName))
			}

			return reflect.ValueOf(depsResolver.Resolve(t)), nil
		},
	}

	return entityName, &entitiyRoot
}

func iterate(
	parentName string,
	parentFields graphql.Fields,
	name string,
	t reflect.Type,
) {
	currentName := fmt.Sprintf("%s_%s", parentName, name)

	// only array, struct and primitives allowed
	switch t.Kind() {

	case reflect.Ptr:
		iterate(parentName, parentFields, name, utils.GetType(t))

	case reflect.Array, reflect.Slice:
		parentFields[name] = &graphql.Field{
			Name: name,
			Type: graphql.String,
		}
		// fields := graphql.Fields{}

		// for i := 0; i < t.NumField(); i++ {
		// 	child := t.Field(i)
		// 	iterate(name, fields, child.Name, child.Type)
		// }

		// obj := graphql.NewObject(graphql.ObjectConfig{
		// 	Name:   t.Name(),
		// 	Fields: fields,
		// })

		// list := graphql.NewList(obj)

		// parentFields[t.Name()] = &graphql.Field{
		// 	Name: t.Name(),
		// 	Type: list,
		// }

	case reflect.Interface:
		parentFields[name] = &graphql.Field{
			Name: name,
			Type: graphql.String,
		}
		// can't make it, noop

	case reflect.Struct:
		fields := graphql.Fields{}

		for i := 0; i < t.NumField(); i++ {
			child := t.Field(i)
			iterate(currentName, fields, child.Name, child.Type)
		}

		// don't do anything for an empty struct
		if len(fields) == 0 {
			return
		}

		obj := graphql.NewObject(graphql.ObjectConfig{
			Name:   currentName,
			Fields: fields,
		})

		parentFields[name] = &graphql.Field{
			Name: name,
			Type: obj,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// otherwise convert to struct here
				// TODO: benchmark against states all written in map[string]interface{}
				// return reflect.Indirect(p.Source).FieldByName(name), nil
				source, ok := p.Source.(reflect.Value)
				if !ok {
					return nil, fmt.Errorf("Invalid source was given, at %s", currentName)
				}

				field := reflect.Indirect(source).FieldByName(name)
				if !field.IsValid() {
					return nil, fmt.Errorf("Field accessor failed, at %s", currentName)
				}

				return field, nil
			},
		}

	// do not do anything for primitives (let the parent handler handle it)
	default:
		gqlType := goTypeToGraphqlType[t.Kind().String()]
		current := &graphql.Field{
			Name: name,
			Type: gqlType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// otherwise convert to struct here
				// TODO: benchmark against states all written in map[string]interface{}
				source, ok := p.Source.(reflect.Value)
				if !ok {
					return nil, fmt.Errorf("Invalid struct was passed down as source, at %s", currentName)
				}

				field := reflect.Indirect(source).FieldByName(name)
				if !field.IsValid() {
					return nil, fmt.Errorf("Field accessor failed, at %s", currentName)
				}

				return field.Interface(), nil
			},
		}

		parentFields[name] = current
	}
}
