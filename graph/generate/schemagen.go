package generate

import (
	"fmt"
	"github.com/terra-project/mantle-sdk/serdes"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql/language/ast"

	"github.com/graphql-go/graphql"
	"github.com/terra-project/mantle-sdk/depsresolver"
	"github.com/terra-project/mantle-sdk/querier"
	"github.com/terra-project/mantle-sdk/utils"
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

func GenerateGraphResolver(modelType reflect.Type) (*graphql.Field, error) {
	t := utils.GetType(modelType)
	entityName := t.Name()

	responseType := buildResponseType(t, "", entityName)
	if responseType == nil {
		return nil, nil
	}

	return &graphql.Field{
		Name: entityName,
		Type: responseType,

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
			if args != nil && len(args) != 0 {
				q, ok := p.Context.Value(utils.QuerierKey).(querier.Querier)
				if !ok {
					panic(fmt.Errorf("querier is either cleared or not set, in ResolverFunc for %s", entityName))
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
					queryResolverIterator.Close()

					doc, err := q.Get(key)
					if err != nil {
						return nil, err
					}
					docValue := reflect.New(t)

					if err := serdes.Deserialize(t, doc, docValue.Interface()); err != nil {
						return nil, fmt.Errorf("could not unmarshal msgpack")
					}

					return docValue.Interface(), err
				}
			}

			depsResolver, depsResolverNotCleared := p.Context.Value(utils.DepsResolverKey).(depsresolver.DepsResolver)
			if !depsResolverNotCleared {
				panic(fmt.Sprintf("DepsResolver is either cleared or not set, in ResolveFunc for %s", entityName))
			}

			// if resolve immediately flag is set, don't await for dependency.
			// load up straight from disk and respond.
			if p.Context.Value(utils.ImmediateResolveFlagKey).(bool) {
				return depsResolver.ResolveLatest(t), nil
			}

			// if current round of graphql resolve is a self-referencing (i.e. awaiting on itself to be resolved)
			// we can't process it here.
			// Must error.
			var dependencies = p.Context.Value(utils.DependenciesKey).(utils.DependenciesKeyType)
			var isAwaitedDependency = p.Args["Height"] == nil
			var isSelfReferencing = dependencies[entityName] == true

			if isAwaitedDependency && isSelfReferencing {
				return nil, fmt.Errorf("Self reference is disallowed. entityName=%s", entityName)
			}

			return depsResolver.Resolve(t), nil
		},
	}, nil
}

func buildResponseType(t reflect.Type, tName string, parentName string) graphql.Output {
	kind := t.Kind()

	if strings.Contains(tName, "XXX_") {
		return nil
	}

	// check if this field is scalar
	if scalar, isScalar := IsCosmosScalar(t); isScalar {
		return scalar
	}

	switch kind {
	// in case of struct,
	case reflect.Struct:
		fields := graphql.Fields{}
		structName := fmt.Sprintf("%s%s", parentName, tName)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			var fieldType graphql.Output

			// see if this field should be implemented as scalar type
			fieldType = buildResponseType(field.Type, field.Name, structName)

			// skip nil fields
			if fieldType == nil {
				continue
			}

			fields[field.Name] = &graphql.Field{
				Name: field.Name,
				Type: fieldType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, isSourceValue := p.Source.(reflect.Value)
					if !isSourceValue {
						source = reflect.ValueOf(p.Source)
					}

					return reflect.Indirect(source).FieldByName(field.Name).Interface(), nil
				},
			}
		}

		return graphql.NewObject(graphql.ObjectConfig{
			Name:   structName,
			Fields: fields,
		})

	case reflect.Interface:
		return graphql.NewScalar(graphql.ScalarConfig{
			Name:         fmt.Sprintf("%s%s", parentName, tName),
			Serialize:    func(value interface{}) interface{} { return value },
			ParseValue:   func(value interface{}) interface{} { return value },
			ParseLiteral: func(valueAST ast.Value) interface{} { return nil },
		})

	// in case of ptr, take Elem() of the type and go deeper
	case reflect.Ptr:
		t := t.Elem()
		return buildResponseType(t, tName, parentName)

		// in case of slice,
	case reflect.Slice, reflect.Array:
		t := t.Elem()
		responseType := buildResponseType(t, tName, parentName)
		if responseType == nil {
			return nil
		}
		return graphql.NewList(responseType)

	default:
		return utils.GetGraphQLType(kind)
	}
}
