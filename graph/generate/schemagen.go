package generate

import (
	"fmt"
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

	name := t.Name()
	fields := graphql.Fields{}

	for i := 0; i < t.NumField(); i++ {
		child := t.Field(i)
		iterate(name, fields, child.Name, child.Type)
	}

	entitiyRoot := graphql.Field{
		Name: t.Name(),
		Type: graphql.NewObject(graphql.ObjectConfig{
			Name:   t.Name(),
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
				// querier, ok := p.Context.Value(utils.QuerierKey).(*querier.Querier)
				// if !ok {
				// 	panic(fmt.Sprintf("Querier is either cleared or not set, in ResolverFunc for %s", name))
				// }
			}

			depsResolver, ok := p.Context.Value(utils.DepsResolverKey).(depsresolver.DepsResolver)
			if !ok {
				panic(fmt.Sprintf("DepsResolver is either cleared or not set, in ResolveFunc for %s", name))
			}

			return reflect.ValueOf(depsResolver.Resolve(t)), nil
		},
	}

	return name, &entitiyRoot
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

				switch gqlType {
				case graphql.String:
					return field.String(), nil
				case graphql.Int:
					return field.Int(), nil
				case graphql.Boolean:
					return field.Bool(), nil
				case graphql.Float:
					return field.Float(), nil
				}

				return nil, fmt.Errorf("Invalid type at %s", currentName)
			},
		}

		parentFields[name] = current
	}
}
